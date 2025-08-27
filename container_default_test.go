package simpledi_test

import (
	"fmt"
	"testing"

	"github.com/eerzho/simpledi"
)

func TestDefaultRegister(t *testing.T) {
	type A struct{}

	t.Run("successful", func(t *testing.T) {
		err := simpledi.Register(simpledi.Def{
			Key:  "a",
			Ctor: func() any { return &A{} },
		})
		assertNoErr(t, err)

		simpledi.MustResolve()
		simpledi.MustReset()
	})

	t.Run("empty key should return ErrEmptyKey", func(t *testing.T) {
		err := simpledi.Register(simpledi.Def{
			Ctor: func() any { return &A{} },
		})
		assertErrType(t, err, simpledi.ErrEmptyKey)

		simpledi.MustResolve()
		simpledi.MustReset()
	})

	t.Run("nil constructor should return ErrNilCtor", func(t *testing.T) {
		err := simpledi.Register(simpledi.Def{
			Key: "a",
		})
		assertErrType(t, err, simpledi.ErrNilCtor)

		simpledi.MustResolve()
		simpledi.MustReset()
	})
}

func TestDefaultMustRegister(t *testing.T) {
	type A struct{}

	t.Run("successful", func(t *testing.T) {
		assertNoPanic(t, func() {
			simpledi.MustRegister(simpledi.Def{
				Key: "a",
				Ctor: func() any {
					return &A{}
				},
			})
		})

		simpledi.MustResolve()
		simpledi.MustReset()
	})

	t.Run("should panic", func(t *testing.T) {
		assertPanic(t, func() {
			simpledi.MustRegister(simpledi.Def{
				Ctor: func() any {
					return &A{}
				},
			})
		})

		simpledi.MustResolve()
		simpledi.MustReset()
	})
}

func TestDefaultGet(t *testing.T) {
	type A struct{}

	t.Run("successful", func(t *testing.T) {
		simpledi.MustRegister(simpledi.Def{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		simpledi.MustResolve()
		object, err := simpledi.Get("a")
		assertNoErr(t, err)
		assertObjectType(t, object, (*A)(nil))

		simpledi.MustReset()
	})

	t.Run("wrong key should return ErrNotFound", func(t *testing.T) {
		simpledi.MustRegister(simpledi.Def{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		simpledi.MustResolve()
		_, err := simpledi.Get("wrong_key")
		assertErrType(t, err, simpledi.ErrNotFound)

		simpledi.MustReset()
	})
}

func TestDefaultGetAs(t *testing.T) {
	type A struct{}
	type B struct{}

	t.Run("successful", func(t *testing.T) {
		simpledi.MustRegister(simpledi.Def{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		simpledi.MustResolve()
		_, err := simpledi.GetAs[*A]("a")
		assertNoErr(t, err)

		simpledi.MustReset()
	})

	t.Run("wrong key should return ErrNotFound", func(t *testing.T) {
		simpledi.MustRegister(simpledi.Def{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		simpledi.MustResolve()
		_, err := simpledi.GetAs[*A]("wrong_key")
		assertErrType(t, err, simpledi.ErrNotFound)

		simpledi.MustReset()
	})

	t.Run("wrong type should return ErrWrongType", func(t *testing.T) {
		simpledi.MustRegister(simpledi.Def{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		simpledi.MustResolve()
		_, err := simpledi.GetAs[*B]("a")
		assertErrType(t, err, simpledi.ErrWrongType)

		simpledi.MustReset()
	})
}

func TestDefaultMustGet(t *testing.T) {
	type A struct{}

	t.Run("successful", func(t *testing.T) {
		simpledi.MustRegister(simpledi.Def{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		simpledi.MustResolve()
		assertNoPanic(t, func() {
			_ = simpledi.MustGet("a")
		})

		simpledi.MustReset()
	})

	t.Run("should panic", func(t *testing.T) {
		simpledi.MustRegister(simpledi.Def{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		simpledi.MustResolve()
		assertPanic(t, func() {
			_ = simpledi.MustGet("wrong_key")
		})

		simpledi.MustReset()
	})
}

func TestDefaultMustGetAs(t *testing.T) {
	type A struct{}

	t.Run("successful", func(t *testing.T) {
		simpledi.MustRegister(simpledi.Def{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		simpledi.MustResolve()
		assertNoPanic(t, func() {
			_ = simpledi.MustGetAs[*A]("a")
		})

		simpledi.MustReset()
	})

	t.Run("should panic", func(t *testing.T) {
		simpledi.MustRegister(simpledi.Def{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		simpledi.MustResolve()
		assertPanic(t, func() {
			_ = simpledi.MustGetAs[*A]("wrong_key")
		})

		simpledi.MustReset()
	})
}

func TestDefaultResolve(t *testing.T) {
	type A struct{}
	type B struct{}

	t.Run("successful", func(t *testing.T) {
		simpledi.MustRegister(simpledi.Def{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		simpledi.MustRegister(simpledi.Def{
			Key:  "b",
			Deps: []string{"a"},
			Ctor: func() any {
				return &B{}
			},
		})
		err := simpledi.Resolve()
		assertNoErr(t, err)

		simpledi.MustReset()
	})

	t.Run("missing dep should return ErrMissingDep", func(t *testing.T) {
		simpledi.MustRegister(simpledi.Def{
			Key:  "a",
			Deps: []string{"b"},
			Ctor: func() any {
				return &A{}
			},
		})
		err := simpledi.Resolve()
		assertErrType(t, err, simpledi.ErrMissingDep)

		simpledi.MustReset()
	})

	t.Run("cyclic dependency should return ErrCyclicDeps", func(t *testing.T) {
		simpledi.MustRegister(simpledi.Def{
			Key:  "a",
			Deps: []string{"b"},
			Ctor: func() any {
				return &A{}
			},
		})
		simpledi.MustRegister(simpledi.Def{
			Key:  "b",
			Deps: []string{"a"},
			Ctor: func() any {
				return &B{}
			},
		})
		err := simpledi.Resolve()
		assertErrType(t, err, simpledi.ErrCyclicDeps)

		simpledi.MustReset()
	})

	t.Run("resolved container should skip resolve", func(t *testing.T) {
		simpledi.MustRegister(simpledi.Def{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		simpledi.MustResolve()
		err := simpledi.Resolve()
		assertNoErr(t, err)

		simpledi.MustReset()
	})

	t.Run("try to resolve empty container", func(t *testing.T) {
		err := simpledi.Resolve()
		assertNoErr(t, err)

		simpledi.MustReset()
	})
}

func TestDefaultMustResolve(t *testing.T) {
	type A struct{}

	t.Run("successful", func(t *testing.T) {
		simpledi.MustRegister(simpledi.Def{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		assertNoPanic(t, func() {
			simpledi.MustResolve()
		})

		simpledi.MustReset()
	})

	t.Run("should panic", func(t *testing.T) {
		simpledi.MustRegister(simpledi.Def{
			Key:  "a",
			Deps: []string{"b"},
			Ctor: func() any {
				return &A{}
			},
		})
		assertPanic(t, func() {
			simpledi.MustResolve()
		})

		simpledi.MustReset()
	})
}

func TestDefaultReset(t *testing.T) {
	type A struct{}

	t.Run("successful", func(t *testing.T) {
		simpledi.MustRegister(simpledi.Def{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		simpledi.MustResolve()
		err := simpledi.Reset()
		assertNoErr(t, err)
	})

	t.Run("destructor should return error", func(t *testing.T) {
		simpledi.MustRegister(simpledi.Def{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
			Dtor: func() error {
				return fmt.Errorf("some error")
			},
		})
		simpledi.MustResolve()
		err := simpledi.Reset()
		assertErr(t, err)
	})

	t.Run("not resolved container should skip reset", func(t *testing.T) {
		simpledi.MustRegister(simpledi.Def{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		simpledi.MustResolve()
		simpledi.MustReset()
		err := simpledi.Reset()
		assertNoErr(t, err)
	})
}

func TestDefaultMustReset(t *testing.T) {
	type A struct{}

	t.Run("successful", func(t *testing.T) {
		simpledi.MustRegister(simpledi.Def{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		simpledi.MustResolve()
		assertNoPanic(t, func() {
			simpledi.MustReset()
		})
	})

	t.Run("should panic", func(t *testing.T) {
		simpledi.MustRegister(simpledi.Def{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
			Dtor: func() error {
				return fmt.Errorf("some error")
			},
		})
		simpledi.MustResolve()
		assertPanic(t, func() {
			simpledi.MustReset()
		})
	})
}
