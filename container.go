package simpledi

import "fmt"

type Container struct {
	deps      map[string][]string
	ctors     map[string]func() any
	instances map[string]any
}

func NewContainer() *Container {
	return &Container{
		deps:      make(map[string][]string),
		ctors:     make(map[string]func() any),
		instances: make(map[string]any),
	}
}

func (c *Container) Register(key string, deps []string, constrcutor func() any) {
	c.deps[key] = deps
	c.ctors[key] = constrcutor
}

func (c *Container) Get(key string) any {
	return c.instances[key]
}

func (c *Container) Resolve() error {
	order, err := c.sort()
	if err != nil {
		return err
	}
	for _, key := range order {
		if _, exists := c.instances[key]; exists {
			continue
		}
		constructor := c.ctors[key]
		if constructor == nil {
			return fmt.Errorf("no constructor for [%s]", key)
		}
		c.instances[key] = constructor()
	}
	return nil
}

func (c *Container) sort() ([]string, error) {
	var queue []string
	graph := make(map[string][]string)
	inDegree := make(map[string]int)
	for key, deps := range c.deps {
		count := len(deps)
		if count == 0 {
			queue = append(queue, key)
		} else {
			inDegree[key] = len(deps)
		}
		for _, subKey := range deps {
			if _, exisits := c.ctors[subKey]; !exisits {
				return nil, fmt.Errorf("missig dependency [%s] required by [%s]", subKey, key)
			}
			graph[subKey] = append(graph[subKey], key)
		}
	}
	var sorted []string
	for len(queue) != 0 {
		key := queue[0]
		queue = queue[1:]
		sorted = append(sorted, key)
		for _, subKey := range graph[key] {
			inDegree[subKey]--
			if inDegree[subKey] == 0 {
				queue = append(queue, subKey)
				delete(inDegree, subKey)
			}
		}
	}
	if len(inDegree) != 0 {
		var cycles []string
		for node := range inDegree {
			cycles = append(cycles, node)
		}
		return nil, fmt.Errorf("cyclic dependency detected %v", cycles)
	}
	return sorted, nil
}
