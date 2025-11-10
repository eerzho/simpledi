package simpledi

import (
	"fmt"
	"sync"
)

var container = sync.OnceValue(func() *Container {
	return New()
})

// Set adds a definition to the container.
func Set(d Definition) {
	if err := container().Set(d); err != nil {
		panic(err)
	}
}

// Get returns an instance by ID.
func Get[T any](id string) T {
	const op = "simpledi.Get"

	var zero T
	instance, err := container().Get(id)
	if err != nil {
		panic(err)
	}
	if instance == nil {
		return zero
	}
	typedInstance, ok := instance.(T)
	if !ok {
		err := fmt.Errorf("%s: %w (ID: %s, Want: %T, Got: %T)", op, ErrTypeMismatch, id, zero, instance)
		panic(err)
	}

	return typedInstance
}

// Resolve creates instances for all registered definitions.
// Dependencies are resolved in topological order based on Deps.
func Resolve() {
	if err := container().Resolve(); err != nil {
		panic(err)
	}
}

// Close calls Close for all definitions that provide it, in reverse order.
// Returns a combined error if any Close calls fail.
// The container is then cleared and can be reused.
func Close() error {
	return container().Close()
}
