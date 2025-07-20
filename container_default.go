package simpledi

import (
	"sync"
)

var (
	c    *Container
	once sync.Once
)

func Register(option Option) error {
	return defaultC().Register(option)
}

func MustRegister(option Option) {
	defaultC().MustRegister(option)
}

func Get(key string) (any, error) {
	return defaultC().Get(key)
}

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

func MustGet(key string) any {
	return defaultC().MustGet(key)
}

func MustGetAs[T any](key string) T {
	object, err := GetAs[T](key)
	if err != nil {
		panic(err)
	}
	return object
}

func Resolve() error {
	return defaultC().Resolve()
}

func MustResolve() {
	defaultC().MustResolve()
}

func Reset() error {
	return defaultC().Reset()
}

func MustReset() {
	defaultC().MustReset()
}

func defaultC() *Container {
	once.Do(func() {
		c = NewContainer()
	})
	return c
}
