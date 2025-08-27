// Package simpledi provides a simple dependency injection container for Go.
// Zero dependencies, no reflection, no code generation.
package simpledi

import (
	"errors"
	"sync"
)

type Def struct {
	Key  string
	Deps []string
	Ctor func() any
	Dtor func() error
}

type Container struct {
	mu           sync.Mutex
	resolved     bool
	objects      map[string]any
	dependencies map[string][]string
	constructors map[string]func() any
	destructors  map[string]func() error
}

func NewContainer() *Container {
	return &Container{
		objects:      make(map[string]any),
		dependencies: make(map[string][]string),
		constructors: make(map[string]func() any),
		destructors:  make(map[string]func() error),
	}
}

func (c *Container) Register(def Def) error {
	if def.Key == "" {
		return errEmptyKey()
	}
	if def.Ctor == nil {
		return errNilCtor(def.Key)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.dependencies[def.Key] = def.Deps
	c.constructors[def.Key] = def.Ctor
	if def.Dtor != nil {
		c.destructors[def.Key] = def.Dtor
	}
	c.resolved = false

	return nil
}

func (c *Container) MustRegister(def Def) {
	if err := c.Register(def); err != nil {
		panic(err)
	}
}

func (c *Container) Get(key string) (any, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	object, ok := c.objects[key]
	if !ok {
		return nil, errNotFound(key)
	}
	return object, nil
}

func (c *Container) MustGet(key string) any {
	object, err := c.Get(key)
	if err != nil {
		panic(err)
	}
	return object
}

func (c *Container) Resolve() error {
	c.mu.Lock()

	if c.resolved {
		c.mu.Unlock()
		return nil
	}

	sorted, err := c.sort()
	if err != nil {
		c.mu.Unlock()
		return err
	}

	for _, key := range sorted {
		constructor := c.constructors[key]

		c.mu.Unlock()
		object := constructor()
		c.mu.Lock()

		c.objects[key] = object
	}
	c.resolved = true

	c.mu.Unlock()

	return nil
}

func (c *Container) MustResolve() {
	if err := c.Resolve(); err != nil {
		panic(err)
	}
}

func (c *Container) Reset() error {
	c.mu.Lock()

	if !c.resolved {
		c.mu.Unlock()
		return nil
	}

	sorted, _ := c.sort()

	var errs []error
	for i := len(sorted) - 1; i >= 0; i-- {
		key := sorted[i]
		if destructor, ok := c.destructors[key]; ok {
			c.mu.Unlock()
			if err := destructor(); err != nil {
				errs = append(errs, err)
			}
			c.mu.Lock()
		}
	}

	c.objects = make(map[string]any)
	c.dependencies = make(map[string][]string)
	c.constructors = make(map[string]func() any)
	c.destructors = make(map[string]func() error)
	c.resolved = false

	c.mu.Unlock()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (c *Container) MustReset() {
	if err := c.Reset(); err != nil {
		panic(err)
	}
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
				return nil, errMissingDep(key, dep)
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
		return nil, errCyclicDeps(cycles)
	}

	return sorted, nil
}
