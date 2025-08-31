package simpledi

import (
	"sync"
)

var (
	c    *Container
	once sync.Once
)

// Register registers a dependency.
//   - def: this is a dependency description
func Register(def Def) error {
	return defaultC().Register(def)
}

// MustRegister registers a dependency or panics on error.
//   - def: this is a dependency description
func MustRegister(def Def) {
	defaultC().MustRegister(def)
}

// Get gets dependency.
//   - key: this is a unique name for the dependency
func Get(key string) (any, error) {
	return defaultC().Get(key)
}

// GetAs gets dependency as type T.
//   - key: this is a unique name for the dependency
func GetAs[T any](key string) (T, error) {
	var zero T

	object, err := Get(key)
	if err != nil {
		return zero, err
	}

	typed, ok := object.(T)
	if !ok {
		return zero, errWrongType(key, zero, object)
	}

	return typed, nil
}

// MustGet gets dependency or panics on error.
//   - key: this is a unique name for the dependency
func MustGet(key string) any {
	return defaultC().MustGet(key)
}

// MustGetAs gets dependency as type T or panics on error.
//   - key: this is a unique name for the dependency
func MustGetAs[T any](key string) T {
	object, err := GetAs[T](key)
	if err != nil {
		panic(err)
	}
	return object
}

// Resolve resolves all dependencies.
func Resolve() error {
	return defaultC().Resolve()
}

// MustResolve resolves all dependencies or panics on error.
func MustResolve() {
	defaultC().MustResolve()
}

// Reset resets the container.
func Reset() error {
	return defaultC().Reset()
}

// MustReset resets the container or panics on error.
func MustReset() {
	defaultC().MustReset()
}

func defaultC() *Container {
	once.Do(func() {
		c = NewContainer()
	})
	return c
}
