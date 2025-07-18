// Package simpledi provides a simple dependency injection container for Go.
// Zero dependencies, no reflection, no code generation.
package simpledi

import (
	"fmt"
	"sync"
)

type ConstructorFunc func() any
type DestructorFunc func()

// Container is a DI container that stores created objects.
// Dependency resolution is based on topological sorting.
type Container struct {
	mu           sync.Mutex
	resolved     bool
	objects      map[string]any
	dependencies map[string][]string
	constructors map[string]ConstructorFunc
	destructors  map[string]DestructorFunc
}

type Option struct {
	Key          string
	Dependencies []string
	Constructor  ConstructorFunc
	Destructor   DestructorFunc
}

// Creates a new DI container.
func NewContainer() *Container {
	return &Container{
		objects:      make(map[string]any),
		dependencies: make(map[string][]string),
		constructors: make(map[string]ConstructorFunc),
		destructors:  make(map[string]DestructorFunc),
	}
}

// Register a dependency by key.
//   - key: unique name for the dependency
//   - deps: list of dependency keys this object depends on
//   - constructor: function that returns the object instance
// func (c *Container) Register(key string, deps []string, constructor func() any) {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()
// 	c.dependencies[key] = deps
// 	c.constructors[key] = constructor
// 	c.resolved = false
// }

func (c *Container) Register(option Option) error {
	if option.Key == "" {
		return fmt.Errorf("some error")
	}

	if option.Constructor == nil {
		return fmt.Errorf("some error")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.dependencies[option.Key] = option.Dependencies
	c.constructors[option.Key] = option.Constructor

	if option.Destructor != nil {
		c.destructors[option.Key] = option.Destructor
	}

	c.resolved = false

	return nil
}

func (c *Container) MustRegister(option Option) {
	if err := c.Register(option); err != nil {
		panic(err)
	}
}

// Get a dependency by key.
//   - key: unique name of the dependency
func (c *Container) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	object, ok := c.objects[key]
	return object, ok
}

// Get a dependency by key or panics.
//   - key: unique name of the dependency
func (c *Container) MustGet(key string) any {
	object, ok := c.Get(key)
	if !ok {
		panic(fmt.Sprintf("dependency [%s] not found", key))
	}
	return object
}

// Resolve all dependencies.
func (c *Container) Resolve() error {
	if c.resolved {
		return nil
	}

	sorted, err := c.sort()
	if err != nil {
		return err
	}

	c.construct(sorted)
	c.resolved = true
	return nil
}

// Resolve all dependencies or panic.
func (c *Container) MustResolve() {
	if err := c.Resolve(); err != nil {
		panic(err)
	}
}

func (c *Container) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.resolved {
		return
	}

	for _, destructor := range c.destructors {
		destructor()
	}

	c.objects = make(map[string]any)
	c.dependencies = make(map[string][]string)
	c.constructors = make(map[string]ConstructorFunc)
	c.destructors = make(map[string]DestructorFunc)

	c.resolved = false
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

func (c *Container) construct(keys []string) {
	for _, key := range keys {
		c.mu.Lock()
		constructor := c.constructors[key]
		c.mu.Unlock()

		c.objects[key] = constructor()
	}
}
