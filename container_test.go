package simpledi_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/eerzho/simpledi"
)

func TestRegister(t *testing.T) {
	type A struct{}

	t.Run("successful", func(t *testing.T) {
		c := simpledi.NewContainer()
		err := c.Register(simpledi.Option{
			Key:  "a",
			Ctor: func() any { return &A{} },
			Dtor: func() error { return nil },
		})
		assertNoErr(t, err)
	})

	t.Run("empty key should return ErrEmptyKey", func(t *testing.T) {
		c := simpledi.NewContainer()
		err := c.Register(simpledi.Option{
			Ctor: func() any { return &A{} },
		})
		assertErrType(t, err, simpledi.ErrEmptyKey)
	})

	t.Run("nil constructor should return ErrNilCtor", func(t *testing.T) {
		c := simpledi.NewContainer()
		err := c.Register(simpledi.Option{
			Key: "a",
		})
		assertErrType(t, err, simpledi.ErrNilCtor)
	})
}

func TestMustRegister(t *testing.T) {
	type A struct{}

	t.Run("successful", func(t *testing.T) {
		c := simpledi.NewContainer()
		assertNoPanic(t, func() {
			c.MustRegister(simpledi.Option{
				Key: "a",
				Ctor: func() any {
					return &A{}
				},
			})
		})
	})

	t.Run("should panic", func(t *testing.T) {
		c := simpledi.NewContainer()
		assertPanic(t, func() {
			c.MustRegister(simpledi.Option{
				Ctor: func() any {
					return &A{}
				},
			})
		})
	})
}

func TestGet(t *testing.T) {
	type A struct{}

	t.Run("successful", func(t *testing.T) {
		c := simpledi.NewContainer()
		c.MustRegister(simpledi.Option{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		c.MustResolve()
		object, err := c.Get("a")
		assertNoErr(t, err)
		assertObjectType(t, object, (*A)(nil))
	})

	t.Run("wrong key should return ErrNotFound", func(t *testing.T) {
		c := simpledi.NewContainer()
		c.MustRegister(simpledi.Option{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		c.MustResolve()
		_, err := c.Get("wrong_key")
		assertErrType(t, err, simpledi.ErrNotFound)
	})
}

func TestMustGet(t *testing.T) {
	type A struct{}

	t.Run("successful", func(t *testing.T) {
		c := simpledi.NewContainer()
		c.MustRegister(simpledi.Option{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		c.MustResolve()
		assertNoPanic(t, func() {
			_ = c.MustGet("a")
		})
	})

	t.Run("should panic", func(t *testing.T) {
		c := simpledi.NewContainer()
		c.MustRegister(simpledi.Option{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		c.MustResolve()
		assertPanic(t, func() {
			_ = c.MustGet("wrong_key")
		})
	})
}

func TestResolve(t *testing.T) {
	type A struct{}
	type B struct{}

	t.Run("successful", func(t *testing.T) {
		c := simpledi.NewContainer()
		c.MustRegister(simpledi.Option{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		c.MustRegister(simpledi.Option{
			Key:  "b",
			Deps: []string{"a"},
			Ctor: func() any {
				return &B{}
			},
		})
		err := c.Resolve()
		assertNoErr(t, err)
	})

	t.Run("missing dep should return ErrMissingDep", func(t *testing.T) {
		c := simpledi.NewContainer()
		c.MustRegister(simpledi.Option{
			Key:  "a",
			Deps: []string{"b"},
			Ctor: func() any {
				return &A{}
			},
		})
		err := c.Resolve()
		assertErrType(t, err, simpledi.ErrMissingDep)
	})

	t.Run("cyclic dependency should return ErrCyclicDeps", func(t *testing.T) {
		c := simpledi.NewContainer()
		c.MustRegister(simpledi.Option{
			Key:  "a",
			Deps: []string{"b"},
			Ctor: func() any {
				return &A{}
			},
		})
		c.MustRegister(simpledi.Option{
			Key:  "b",
			Deps: []string{"a"},
			Ctor: func() any {
				return &B{}
			},
		})
		err := c.Resolve()
		assertErrType(t, err, simpledi.ErrCyclicDeps)
	})

	t.Run("resolved container should skip resolve", func(t *testing.T) {
		c := simpledi.NewContainer()
		c.MustRegister(simpledi.Option{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		c.MustResolve()
		err := c.Resolve()
		assertNoErr(t, err)
	})

	t.Run("try to resolve empty container", func(t *testing.T) {
		c := simpledi.NewContainer()
		err := c.Resolve()
		assertNoErr(t, err)
	})
}

func TestMustResolve(t *testing.T) {
	type A struct{}

	t.Run("successful", func(t *testing.T) {
		c := simpledi.NewContainer()
		c.MustRegister(simpledi.Option{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		assertNoPanic(t, func() {
			c.MustResolve()
		})
	})

	t.Run("should panic", func(t *testing.T) {
		c := simpledi.NewContainer()
		c.MustRegister(simpledi.Option{
			Key:  "a",
			Deps: []string{"b"},
			Ctor: func() any {
				return &A{}
			},
		})
		assertPanic(t, func() {
			c.MustResolve()
		})
	})
}

func TestReset(t *testing.T) {
	type A struct{}

	t.Run("successful", func(t *testing.T) {
		c := simpledi.NewContainer()
		c.MustRegister(simpledi.Option{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		c.MustResolve()
		err := c.Reset()
		assertNoErr(t, err)
	})

	t.Run("destructor should return error", func(t *testing.T) {
		c := simpledi.NewContainer()
		c.MustRegister(simpledi.Option{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
			Dtor: func() error {
				return fmt.Errorf("some error")
			},
		})
		c.MustResolve()
		err := c.Reset()
		assertErr(t, err)
	})

	t.Run("not resolved container should skip reset", func(t *testing.T) {
		c := simpledi.NewContainer()
		c.MustRegister(simpledi.Option{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		c.MustResolve()
		c.MustReset()
		err := c.Reset()
		assertNoErr(t, err)
	})
}

func TestMustReset(t *testing.T) {
	type A struct{}

	t.Run("successful", func(t *testing.T) {
		c := simpledi.NewContainer()
		c.MustRegister(simpledi.Option{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
		})
		c.MustResolve()
		assertNoPanic(t, func() {
			c.MustReset()
		})
	})

	t.Run("should panic", func(t *testing.T) {
		c := simpledi.NewContainer()
		c.MustRegister(simpledi.Option{
			Key: "a",
			Ctor: func() any {
				return &A{}
			},
			Dtor: func() error {
				return errors.New("some error")
			},
		})
		c.MustResolve()
		assertPanic(t, func() {
			c.MustReset()
		})
	})
}

func assertNoErr(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("got: %v, want: no error", err)
	}
}

func assertErr(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatalf("got: no error, want: error")
	}
}

func assertErrType(t *testing.T, err error, want simpledi.ErrorType) {
	t.Helper()

	assertErr(t, err)

	var diErr *simpledi.Error
	if !errors.As(err, &diErr) {
		t.Fatalf("got: %T, want: %T", err, diErr)
	}

	if diErr.Type != want {
		t.Fatalf("got: %v, want: %v", diErr.Type, want)
	}
}

func assertObjectType(t *testing.T, got, want any) {
	t.Helper()

	gotStr := fmt.Sprintf("%T", got)
	wantStr := fmt.Sprintf("%T", want)

	if gotStr != wantStr {
		t.Fatalf("got: %s, want: %s", gotStr, wantStr)
	}
}

func assertPanic(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("got: no panic, want: panic")
		}
	}()

	fn()
}

func assertNoPanic(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("got: panic, want: no panic")
		}
	}()

	fn()
}
