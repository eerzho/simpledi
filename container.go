// A simple dependency injection container for Go — zero dependencies, no reflection, no code generation.
package simpledi

import (
	"fmt"
	"sync"
)

// DI Container that stores created objects.
// Dependency resolution is based on topological sorting.
type Container struct {
	ts       *topoSort
	mu       sync.RWMutex
	objects  map[string]any
	builders map[string]func() any
}

// Creates and returns a new DI container.
func NewContainer() *Container {
	return &Container{
		ts:       newTopoSort(),
		objects:  make(map[string]any),
		builders: make(map[string]func() any),
	}
}

// Register a dependency by key
//   - key:     unique name for the dependency
//   - needs:   list of dependency keys this object depends on
//   - builder: function that returns the object instance
func (c *Container) Register(key string, needs []string, builder func() any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ts.append(key, needs)
	c.builders[key] = builder
}

// Get a dependency by key
//   - key: unique name of the dependency
func (c *Container) Get(key string) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.objects[key]
}

// Resolve all dependencies
func (c *Container) Resolve() error {
	sorted, err := c.ts.sort()
	if err != nil {
		return err
	}
	for _, key := range sorted {
		builder := c.builders[key]
		if builder == nil {
			return fmt.Errorf("[%s] builder is nil", key)
		}
		object := builder()
		c.mu.Lock()
		c.objects[key] = object
		c.mu.Unlock()
	}
	return nil
}
