package simpledi

import (
	"errors"
	"fmt"
	"sync"
)

var (
	errContainerResolved      = errors.New("container already resolved")
	errIDEmpty                = errors.New("ID is empty")
	errConstructorRequired    = errors.New("constructor required")
	errContainerNotResolved   = errors.New("container not resolved")
	errNotFound               = errors.New("not found")
	errCyclicDependency       = errors.New("cyclic dependency detected")
)

var ctr = sync.OnceValue(func() *container {
	return &container{
		definitions: make([]Definition, 0),
		instances:   make(map[string]any),
	}
})

type Definition struct {
	ID    string
	Deps  []string
	New   func() any
	Close func() error
}

func Set(d Definition) {
	if err := ctr().set(d); err != nil {
		panic(err)
	}
}

func Get[T any](id string) T {
	object, err := ctr().get(id)
	if err != nil {
		panic(err)
	}
	typedObject, ok := object.(T)
	if !ok {
		panic(fmt.Errorf("Get: %q type mismatch, got %T", id, object))
	}
	return typedObject
}

func Close() error {
	return ctr().close()
}

func Resolve() {
	if err := ctr().resolve(); err != nil {
		panic(err)
	}
}

type container struct {
	resolved    bool
	definitions []Definition
	instances   map[string]any
}

func (c *container) set(d Definition) error {
	if c.resolved {
		return fmt.Errorf("set: container already resolved")
	}
	if d.ID == "" {
		return fmt.Errorf("set: ID is empty")
	}
	if d.New == nil {
		return fmt.Errorf("set: %q has no constructor", d.ID)
	}
	c.definitions = append(c.definitions, d)
	return nil
}

func (c *container) get(id string) (any, error) {
	if !c.resolved {
		return nil, fmt.Errorf("get: container not resolved")
	}
	object, ok := c.instances[id]
	if !ok {
		return nil, fmt.Errorf("get: %q not found", id)
	}
	return object, nil
}

func (c *container) resolve() error {
	if c.resolved {
		return fmt.Errorf("resolve: already resolved")
	}
	if err := c.sort(); err != nil {
		return err
	}
	for _, definition := range c.definitions {
		instance := definition.New()
		c.instances[definition.ID] = instance
	}
	c.resolved = true
	return nil
}

func (c *container) close() error {
	if !c.resolved {
		return fmt.Errorf("close: not resolved")
	}
	errs := make([]error, 0)
	for i := len(c.definitions) - 1; i >= 0; i-- {
		definition := c.definitions[i]
		if definition.Close != nil {
			if err := definition.Close(); err != nil {
				errs = append(errs, fmt.Errorf("close: %q failed: %w", definition.ID, err))
			}
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (c *container) sort() error {
	definitionsCount := len(c.definitions)
	if definitionsCount == 0 {
		return nil
	}

	inDegree := make(map[string]int, definitionsCount)
	for _, definition := range c.definitions {
		inDegree[definition.ID] = len(definition.Deps)
	}

	queue := make([]Definition, 0, definitionsCount)
	graph := make(map[string][]Definition, definitionsCount)
	for _, definition := range c.definitions {
		if inDegree[definition.ID] == 0 {
			queue = append(queue, definition)
			continue
		}
		for _, dependency := range definition.Deps {
			if _, ok := inDegree[dependency]; !ok {
				return fmt.Errorf("sort: %q depends on %q which is not registered", definition.ID, dependency)
			}
			graph[dependency] = append(graph[dependency], definition)
		}
	}

	sortedDefinitions := make([]Definition, 0, definitionsCount)
	for len(queue) > 0 {
		definition := queue[0]
		queue = queue[1:]
		sortedDefinitions = append(sortedDefinitions, definition)
		for _, subDefinition := range graph[definition.ID] {
			inDegree[subDefinition.ID]--
			if inDegree[subDefinition.ID] == 0 {
				queue = append(queue, subDefinition)
			}
		}
	}

	sortedDefinitionsCount := len(sortedDefinitions)
	if definitionsCount != sortedDefinitionsCount {
		cycles := make([]string, 0, definitionsCount-sortedDefinitionsCount)
		for key, degree := range inDegree {
			if degree > 0 {
				cycles = append(cycles, key)
			}
		}
		return fmt.Errorf("sort: cyclic dependency detected: %v", cycles)
	}
	c.definitions = sortedDefinitions
	return nil
}
