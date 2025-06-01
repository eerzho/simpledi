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
