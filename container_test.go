package simpledi_test

import (
	"strings"
	"testing"

	"github.com/eerzho/simpledi"
)

func TestResolve(t *testing.T) {
	c := simpledi.NewContainer()
	c.MustRegister(simpledi.Option{
		Key:         "a",
		Constructor: func() any { return "a" },
	})
	c.MustRegister(simpledi.Option{
		Key:          "b",
		Dependencies: []string{"a"},
		Constructor:  func() any { return "b" },
	})
	c.MustRegister(simpledi.Option{
		Key:          "c",
		Dependencies: []string{"a", "b"},
		Constructor:  func() any { return "c" },
	})

	err := c.Resolve()
	if err != nil {
		t.Fatalf("got: err - %s, want: err - nil", err.Error())
	}
}

func TestResolveSkipWhenResolved(t *testing.T) {
	c := simpledi.NewContainer()
	count := 0
	c.MustRegister(simpledi.Option{
		Key: "a",
		Constructor: func() any {
			count++
			return "a"
		},
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
	c.MustRegister(simpledi.Option{
		Key:          "b",
		Dependencies: []string{"a"},
		Constructor:  func() any { return "b" },
	})
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
	c.MustRegister(simpledi.Option{
		Key:          "a",
		Dependencies: []string{"b"},
		Constructor:  func() any { return "a" },
	})
	c.MustRegister(simpledi.Option{
		Key:          "b",
		Dependencies: []string{"a"},
		Constructor:  func() any { return "b" },
	})
	err := c.Resolve()
	if err == nil {
		t.Fatalf("got: err - nil, want: err - %s", want)
	}
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("got: err - %s, want: err - %s", err.Error(), want)
	}
}
