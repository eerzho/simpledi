package simpledi_test

import (
	"strings"
	"testing"

	"github.com/eerzho/simpledi"
)

func TestResolve(t *testing.T) {
	c := simpledi.NewContainer()
	c.Register("a", nil, func() any { return "a" })
	c.Register("b", []string{"a"}, func() any { return "b" })
	c.Register("c", []string{"a", "b"}, func() any { return "c" })
	err := c.Resolve()
	if err != nil {
		t.Fatalf("got: err - %s, want: err - nil", err.Error())
	}
}

func TestResolveSkipWhenResolved(t *testing.T) {
	c := simpledi.NewContainer()
	count := 0
	c.Register("a", nil, func() any {
		count++
		return "a"
	})
	err := c.Resolve()
	if err != nil {
		t.Fatalf("got: err - %s, want: err - nil", err.Error())
	}
	if count != 1 {
		t.Fatalf("got: count - %d, want: count - 1", count)
	}
	err = c.Resolve()
	if err != nil {
		t.Fatalf("got: err - %s, want: err - nil", err.Error())
	}
	if count != 1 {
		t.Fatalf("got: count - %d, want: count - 1", count)
	}
}

func TestResolveWhenEmptyContainer(t *testing.T) {
	c := simpledi.NewContainer()
	err := c.Resolve()
	if err != nil {
		t.Fatalf("got: err - %s, want: err - nil", err.Error())
	}
}

func TestResolveWhenKeyNotDeclared(t *testing.T) {
	want := "not declared"
	c := simpledi.NewContainer()
	c.Register("b", []string{"a"}, func() any { return "b" })
	err := c.Resolve()
	if err == nil {
		t.Fatalf("got: err - nil, want: err - %s", want)
	}
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("got: err - %s, want: err - %s", err.Error(), want)
	}
}

func TestResolveWhenCyclicDetected(t *testing.T) {
	want := "cyclic detected"
	c := simpledi.NewContainer()
	c.Register("a", []string{"b"}, func() any { return "a" })
	c.Register("b", []string{"a"}, func() any { return "b" })
	err := c.Resolve()
	if err == nil {
		t.Fatalf("got: err - nil, want: err - %s", want)
	}
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("got: err - %s, want: err - %s", err.Error(), want)
	}
}

func TestResolveWhenBuilderIsNil(t *testing.T) {
	want := "builder is nil"
	c := simpledi.NewContainer()
	c.Register("a", nil, nil)
	err := c.Resolve()
	if err == nil {
		t.Fatalf("got: err - nil, want: err - %s", want)
	}
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("got: err - %s, want: err - %s", err.Error(), want)
	}
}

func TestMustGet(t *testing.T) {
	c := simpledi.NewContainer()
	c.Register("a", nil, func() any { return "a" })
	c.MustResolve()

	a := c.MustGet("a")
	if a != "a" {
		t.Fatalf("got: %s, want: a", a)
	}
}

func TestMustGetPanic(t *testing.T) {
	c := simpledi.NewContainer()
	c.MustResolve()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("got: no panic, want: panic")
		}
	}()

	c.MustGet("nonexistent")
}

func TestGetAs(t *testing.T) {
	simpledi.Register("a", nil, func() any { return "a" })
	simpledi.MustResolve()

	a, ok := simpledi.GetAs[string]("a")
	if !ok {
		t.Fatalf("got: false, want: true")
	}
	if a != "a" {
		t.Fatalf("got: %s, want: a", a)
	}
}

func TestGetAsFail(t *testing.T) {
	simpledi.Register("a", nil, func() any { return "a" })
	simpledi.MustResolve()

	_, ok := simpledi.GetAs[int]("a")
	if ok {
		t.Fatalf("got: true, want: false")
	}
}

func TestMustGetAs(t *testing.T) {
	simpledi.Register("a", nil, func() any { return "a" })
	simpledi.MustResolve()

	a := simpledi.MustGetAs[string]("a")
	if a != "a" {
		t.Fatalf("got: %s, want: a", a)
	}
}

func TestMustGetAsPanic(t *testing.T) {
	simpledi.Register("a", nil, func() any { return "a" })
	simpledi.MustResolve()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("got: no panic, want: panic")
		}
	}()

	simpledi.MustGetAs[int]("a")
}
