package simpledi

import (
	"fmt"
	"sync"
)

type Container struct {
	ts       *topoSort
	mu       sync.RWMutex
	objects  map[string]any
	builders map[string]func() any
}

func NewContainer() *Container {
	return &Container{
		ts:       newTopoSort(),
		objects:  make(map[string]any),
		builders: make(map[string]func() any),
	}
}

func (c *Container) Register(key string, needs []string, builder func() any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ts.append(key, needs)
	c.builders[key] = builder
}

func (c *Container) Get(key string) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.objects[key]
}

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
