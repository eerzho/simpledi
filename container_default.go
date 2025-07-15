package simpledi

import (
	"fmt"
	"sync"
)

var (
	c    *Container
	once sync.Once
)

// Register a dependency by key
//   - key:     unique name for the dependency
//   - needs:   list of dependency keys this object depends on
//   - builder: function that returns the object instance
func Register(key string, needs []string, builder func() any) {
	defaultC().Register(key, needs, builder)
}

// Get a dependency by key
//   - key: unique name of the dependency
func Get(key string) (any, bool) {
	return defaultC().Get(key)
}

// Get a dependency by key and casts it to the specified type
//   - key: unique name of the dependency
func GetAs[T any](key string) (T, bool) {
	var zero T

	object, ok := Get(key)
	if !ok {
		return zero, false
	}

	typed, ok := object.(T)
	if !ok {
		return zero, false
	}

	return typed, true
}

// Get a dependency by key or panics
//   - key: unique name of the dependency
func MustGet(key string) any {
	return defaultC().MustGet(key)
}

// Get a dependency by key and casts it to the specified type or panics
//   - key: unique name of the dependency
func MustGetAs[T any](key string) T {
	var zero T

	object := MustGet(key)
	typed, ok := object.(T)
	if !ok {
		panic(fmt.Sprintf("dependency [%s] cannot be cast to %T", key, zero))
	}

	return typed
}

// Resolve all dependencies
func Resolve() error {
	return defaultC().Resolve()
}

func defaultC() *Container {
	once.Do(func() {
		c = NewContainer()
	})
	return c
}
