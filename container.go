package simpledi

import "fmt"

type Container struct {
	dependencies map[string][]string
	constructors map[string]func() any
	instances    map[string]any
}

func NewContainer() *Container {
	return &Container{
		dependencies: make(map[string][]string),
		constructors: make(map[string]func() any),
		instances:    make(map[string]any),
	}
}

func (c *Container) Register(name string, deps []string, constrcutor func() any) {
	c.dependencies[name] = deps
	c.constructors[name] = constrcutor
}

func (c *Container) Get(name string) any {
	return c.instances[name]
}

func (c *Container) Resolve() error {
	order, err := c.sort()
	if err != nil {
		return err
	}
	for _, name := range order {
		if _, exists := c.instances[name]; exists {
			continue
		}
		constructor := c.constructors[name]
		if constructor == nil {
			return fmt.Errorf("no constructor for %s", name)
		}
		c.instances[name] = constructor()
	}
	return nil
}

func (c *Container) sort() ([]string, error) {
	var queue []string
	graph := make(map[string][]string)
	inDegree := make(map[string]int)
	for name, deps := range c.dependencies {
		count := len(deps)
		if count == 0 {
			queue = append(queue, name)
		} else {
			inDegree[name] = len(deps)
		}
		for _, nb := range deps {
			if _, exisits := c.constructors[nb]; !exisits {
				return nil, fmt.Errorf("missig dependency [%s] required by [%s]", nb, name)
			}
			graph[nb] = append(graph[nb], name)
		}
	}
	var sorted []string
	for len(queue) != 0 {
		name := queue[0]
		queue = queue[1:]
		sorted = append(sorted, name)
		for _, nb := range graph[name] {
			inDegree[nb]--
			if inDegree[nb] == 0 {
				queue = append(queue, nb)
				delete(inDegree, nb)
			}
		}
	}
	if len(inDegree) != 0 {
		var cycles []string
		for node := range inDegree {
			cycles = append(cycles, node)
		}
		return nil, fmt.Errorf("cyclic dependency detected: %v", cycles)
	}
	return sorted, nil
}
