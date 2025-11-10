package simpledi

import (
	"fmt"
	"sync"
)

var container = sync.OnceValue(func() *Container {
	return New()
})

type Definition struct {
	ID    string
	Deps  []string
	New   func() any
	Close func() error
}

func Set(d Definition) {
	if err := container().Set(d); err != nil {
		panic(err)
	}
}

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

func Resolve() {
	if err := container().Resolve(); err != nil {
		panic(err)
	}
}

func Close() error {
	return container().close()
}
