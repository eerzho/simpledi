// Package simpledi provides a simple dependency injection container for Go.
// Zero dependencies, no reflection, no code generation.
package simpledi

import (
	"errors"
	"fmt"
)

var (
	// ErrContainerResolved indicates that the operation was called
	// after the container has already been resolved.
	ErrContainerResolved = errors.New("Container resolved")

	// ErrIDRequired indicates that a definition has no ID.
	ErrIDRequired = errors.New("ID required")

	// ErrNewRequired indicates that a definition has no constructor function.
	ErrNewRequired = errors.New("New required")

	// ErrContainerNotResolved indicates that an operation requires
	// a resolved container, but the container is not yet resolved.
	ErrContainerNotResolved = errors.New("Container not resolved")

	// ErrIDNotFound indicates that no instance with the given ID exists.
	ErrIDNotFound = errors.New("ID not found")

	// ErrIDDuplicate indicates that a definition with the same ID already exists.
	ErrIDDuplicate = errors.New("ID duplicate")

	// ErrDependencyNotFound indicates that a dependency in Deps is not defined.
	ErrDependencyNotFound = errors.New("Dependency not found")

	// ErrDependencyCycle indicates that a circular dependency was detected.
	ErrDependencyCycle = errors.New("Dependency cycle detected")

	// ErrTypeMismatch indicates that a requested instance type does not match.
	ErrTypeMismatch = errors.New("Type mismatch")
)

// Definition describes a dependency definition.
type Definition struct {
	// ID is the unique identifier of the definition. Required.
	ID string
	// Deps is the list of dependency IDs. Optional.
	Deps []string
	// New is the function that returns a new instance. Required.
	New func() any
	// Close is the function called on container close. Optional.
	Close func() error
}

// Container is a simple dependency injection container.
//
// A container stores definitions, resolves their dependencies,
// creates instances, and manages cleanup.
type Container struct {
	resolved    bool
	definitions []Definition
	instances   map[string]any
}

// New returns a new Container.
func New() *Container {
	return &Container{
		definitions: make([]Definition, 0),
		instances:   make(map[string]any),
	}
}

// Set adds a definition to the container.
func (c *Container) Set(d Definition) error {
	const op = "simpledi.Set"

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

// Get returns an instance by ID.
func (c *Container) Get(id string) (any, error) {
	const op = "simpledi.Get"

	if id == "" {
		return nil, fmt.Errorf("%s: %w", op, ErrIDRequired)
	}
	instance, ok := c.instances[id]
	if !ok {
		return nil, fmt.Errorf("%s: %w (ID: %s)", op, ErrIDNotFound, id)
	}

	return instance, nil
}

// Resolve creates instances for all registered definitions.
// Dependencies are resolved in topological order based on Deps.
func (c *Container) Resolve() error {
	const op = "simpledi.Resolve"

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

// Close calls Close for all definitions that provide it, in reverse order.
// Returns a combined error if any Close calls fail.
// The container is then cleared and can be reused.
func (c *Container) Close() error {
	const op = "simpledi.Close"

	errs := make([]error, 0)
	if c.resolved {
		for i := len(c.definitions) - 1; i >= 0; i-- {
			definition := c.definitions[i]
			if definition.Close != nil {
				if err := definition.Close(); err != nil {
					errs = append(errs, fmt.Errorf("%s: %w (ID: %s)", op, err, definition.ID))
				}
			}
		}
	}

	c.definitions = make([]Definition, 0)
	c.instances = make(map[string]any)
	c.resolved = false

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (c *Container) sort() error {
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
	queueIdx := 0
	for queueIdx < len(queue) {
		definition := queue[queueIdx]
		sortedDefinitions = append(sortedDefinitions, definition)
		for _, subDefinition := range graph[definition.ID] {
			inDegree[subDefinition.ID]--
			if inDegree[subDefinition.ID] == 0 {
				queue = append(queue, subDefinition)
			}
		}
		queueIdx++
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
