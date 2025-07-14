// Package simpledi provides a simple dependency injection container for Go.
// Zero dependencies, no reflection, no code generation.
package simpledi

import (
	"fmt"
	"sync"
)

// Container is a DI container that stores created objects.
// Dependency resolution is based on topological sorting.
type Container struct {
	mu           sync.Mutex
	resolved     bool
	objects      map[string]any
	builders     map[string]func() any
	dependencies map[string][]string
}

// NewContainer creates and returns a new DI container.
func NewContainer() *Container {
	return &Container{
		objects:      make(map[string]any),
		builders:     make(map[string]func() any),
		dependencies: make(map[string][]string),
	}
}

// Register registers a dependency by key.
//   - key: unique name for the dependency
//   - deps: list of dependency keys this object depends on
//   - builder: function that returns the object instance
func (c *Container) Register(key string, deps []string, builder func() any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dependencies[key] = deps
	c.builders[key] = builder
	c.resolved = false
}

// Get retrieves a dependency by key.
func (c *Container) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	object, ok := c.objects[key]
	return object, ok
}

// MustGet retrieves a dependency by key or panics if not found.
func (c *Container) MustGet(key string) any {
	object, ok := c.Get(key)
	if !ok {
		panic(fmt.Sprintf("dependency [%s] not found", key))
	}
	return object
}

// Resolve resolves all dependencies.
func (c *Container) Resolve() error {
	if c.resolved {
		return nil
	}
	sorted, err := c.sort()
	if err != nil {
		return err
	}
	if err := c.build(sorted); err != nil {
		return err
	}
	c.resolved = true
	return nil
}

func (c *Container) sort() ([]string, error) {
	depsCount := len(c.dependencies)
	if depsCount == 0 {
		return nil, nil
	}
	inDegree := make(map[string]int, depsCount)
	for key := range c.dependencies {
		inDegree[key] = 0
	}
	for key, deps := range c.dependencies {
		for _, dep := range deps {
			if _, exists := c.dependencies[dep]; !exists {
				return nil, fmt.Errorf("[%s] not declared", dep)
			}
			inDegree[key]++
		}
	}
	queue := make([]string, 0, depsCount)
	for key, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, key)
		}
	}
	sorted := make([]string, 0, depsCount)
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		sorted = append(sorted, current)
		for key, deps := range c.dependencies {
			for _, dep := range deps {
				if dep == current {
					inDegree[key]--
					if inDegree[key] == 0 {
						queue = append(queue, key)
					}
					break
				}
			}
		}
	}
	sortedCount := len(sorted)
	if sortedCount != depsCount {
		cycles := make([]string, 0, depsCount-sortedCount)
		for key, degree := range inDegree {
			if degree > 0 {
				cycles = append(cycles, key)
			}
		}
		return nil, fmt.Errorf("cyclic detected: %v", cycles)
	}
	return sorted, nil
}

func (c *Container) build(keys []string) error {
	for _, key := range keys {
		c.mu.Lock()
		builder := c.builders[key]
		c.mu.Unlock()
		if builder == nil {
			return fmt.Errorf("[%s] builder is nil", key)
		}
		c.objects[key] = builder()
	}
	return nil
}
