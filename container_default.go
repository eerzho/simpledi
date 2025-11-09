package simpledi

import (
	"fmt"
	"sync"
)

var ctr = sync.OnceValue(func() *container {
	return &container{
		definitions: make([]Definition, 0),
		instances:   make(map[string]any),
	}
})

type Definition struct {
	ID    string
	Deps  []string
	New   func() any
	Close func() error
}

func Set(d Definition) {
	if err := ctr().set(d); err != nil {
		panic(err)
	}
}

func Get[T any](id string) T {
	const op = "simpledi.get"

	var zero T
	instance, err := ctr().get(id)
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
	if err := ctr().resolve(); err != nil {
		panic(err)
	}
}

func Close() error {
	return ctr().close()
}
