package simpledi

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrContainerResolved    = errors.New("Container resolved")
	ErrIDRequired           = errors.New("ID required")
	ErrNewRequired          = errors.New("New required")
	ErrContainerNotResolved = errors.New("Container not resolved")
	ErrIDNotFound           = errors.New("ID not found")
	ErrIDDuplicate          = errors.New("ID duplicate")
	ErrDependencyNotFound   = errors.New("Dependency not found")
	ErrDependencyCycle      = errors.New("Dependency cycle detected")
	ErrTypeMismatch         = errors.New("Type mismatch")
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
	const op = "simpledi.get"

	var zero T
	object, err := ctr().get(id)
	if err != nil {
		panic(err)
	}
	typedObject, ok := object.(T)
	if !ok {
		err := fmt.Errorf("%s: %w (ID: %s, Want: %T, Got: %T)", op, ErrTypeMismatch, id, zero, object)
		panic(err)
	}

	return typedObject
}

func Resolve() {
	if err := ctr().resolve(); err != nil {
		panic(err)
	}
}

func Close() error {
	return ctr().close()
}

type container struct {
	resolved    bool
	definitions []Definition
	instances   map[string]any
}

func (c *container) set(d Definition) error {
	const op = "simpledi.set"

	if c.resolved {
		return fmt.Errorf("%s: %w", op, ErrContainerResolved)
	}
	if d.ID == "" {
		return fmt.Errorf("%s: %w", op, ErrIDRequired)
	}
	if d.New == nil {
		return fmt.Errorf("%s: %w (ID: %s)", op, ErrNewRequired, d.ID)
	}
	c.definitions = append(c.definitions, d)

	return nil
}

func (c *container) get(id string) (any, error) {
	const op = "simpledi.get"

	if !c.resolved {
		return nil, fmt.Errorf("%s: %w", op, ErrContainerNotResolved)
	}
	object, ok := c.instances[id]
	if !ok {
		return nil, fmt.Errorf("%s: %w (ID: %s)", op, ErrIDNotFound, id)
	}

	return object, nil
}

func (c *container) resolve() error {
	const op = "simpledi.resolve"

	if c.resolved {
		return fmt.Errorf("%s: %w", op, ErrContainerResolved)
	}
	if err := c.sort(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	for _, definition := range c.definitions {
		instance := definition.New()
		c.instances[definition.ID] = instance
	}
	c.resolved = true

	return nil
}

func (c *container) sort() error {
	const op = "simpledi.sort"

	definitionsCount := len(c.definitions)
	if definitionsCount == 0 {
		return nil
	}

	inDegree := make(map[string]int, definitionsCount)
	for _, definition := range c.definitions {
		if _, ok := inDegree[definition.ID]; ok {
			return fmt.Errorf("%s: %w (ID: %s)", op, ErrIDDuplicate, definition.ID)
		}
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
				return fmt.Errorf("%s: %w (ID: %s, Dependency: %s)", op, ErrDependencyNotFound, definition.ID, dependency)
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
		return fmt.Errorf("%s: %w (Cycles: %v)", op, ErrDependencyCycle, cycles)
	}

	c.definitions = sortedDefinitions

	return nil
}

func (c *container) close() error {
	const op = "simpledi.close"

	if c.resolved {
		errs := make([]error, 0)
		for i := len(c.definitions) - 1; i >= 0; i-- {
			definition := c.definitions[i]
			if definition.Close != nil {
				if err := definition.Close(); err != nil {
					errs = append(errs, fmt.Errorf("%s: %w (ID: %s)", op, err, definition.ID))
				}
			}
		}
		if len(errs) > 0 {
			return errors.Join(errs...)
		}
	}

	c.definitions = make([]Definition, 0)
	c.instances = make(map[string]any)
	c.resolved = false

	return nil
}
