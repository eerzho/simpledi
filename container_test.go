package simpledi_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/eerzho/simpledi"
)

func TestResolve_BasicDependencies(t *testing.T) {
	c := simpledi.NewContainer()

	c.Register("a", nil, func() any { return "a-object" })
	c.Register("b", nil, func() any { return "b-object" })
	c.Register("d", []string{"a", "b"}, func() any {
		a := c.Get("a").(string)
		b := c.Get("b").(string)
		return fmt.Sprintf("d-object: '%s' + '%s'", a, b)
	})
	c.Register("f", []string{"a", "b", "d"}, func() any {
		a := c.Get("a").(string)
		b := c.Get("b").(string)
		d := c.Get("d").(string)
		return fmt.Sprintf("f-object: '%s' + '%s' + '%s'", a, b, d)
	})

	err := c.Resolve()
	if err != nil {
		t.Fatalf("got: %s, want: nil", err.Error())
	}
}

func TestResolve_DuplicateRegistration(t *testing.T) {
	c := simpledi.NewContainer()

	c.Register("a", nil, func() any { return "a-object" })
	c.Register("a", []string{"b"}, func() any {
		b := c.Get("b").(string)
		return fmt.Sprintf("a-bject: '%s'", b)
	})
	c.Register("b", nil, func() any { return "b-object" })
	c.Register("d", []string{"a", "b"}, func() any {
		a := c.Get("a").(string)
		b := c.Get("b").(string)
		return fmt.Sprintf("d-object: '%s' + '%s'", a, b)
	})
	c.Register("f", []string{"a", "b", "d"}, func() any {
		a := c.Get("a").(string)
		b := c.Get("b").(string)
		d := c.Get("d").(string)
		return fmt.Sprintf("f-object: '%s' + '%s' + '%s'", a, b, d)
	})

	err := c.Resolve()
	if err != nil {
		t.Fatalf("got: %s, want: nil", err.Error())
	}
}

func TestResolve_AnyRegistrationOrder(t *testing.T) {
	c := simpledi.NewContainer()

	c.Register("f", []string{"a", "b", "d"}, func() any {
		a := c.Get("a").(string)
		b := c.Get("b").(string)
		d := c.Get("d").(string)
		return fmt.Sprintf("f-object: '%s' + '%s' + '%s'", a, b, d)
	})
	c.Register("d", []string{"a", "b"}, func() any {
		a := c.Get("a").(string)
		b := c.Get("b").(string)
		return fmt.Sprintf("d-object: '%s' + '%s'", a, b)
	})
	c.Register("b", nil, func() any { return "b-object" })
	c.Register("a", nil, func() any { return "a-object" })

	err := c.Resolve()
	if err != nil {
		t.Fatalf("got: %s, want: nil", err.Error())
	}
}

func TestResolve_AnyOrderWithDuplicates(t *testing.T) {
	c := simpledi.NewContainer()

	c.Register("f", []string{"a", "b", "d"}, func() any {
		a := c.Get("a").(string)
		b := c.Get("b").(string)
		d := c.Get("d").(string)
		return fmt.Sprintf("f-object: '%s' + '%s' + '%s'", a, b, d)
	})
	c.Register("d", []string{"a", "b"}, func() any {
		a := c.Get("a").(string)
		b := c.Get("b").(string)
		return fmt.Sprintf("d-object: '%s' + '%s'", a, b)
	})
	c.Register("b", nil, func() any { return "b-object" })
	c.Register("a", []string{"b"}, func() any {
		b := c.Get("b").(string)
		return fmt.Sprintf("a-bject: '%s'", b)
	})
	c.Register("a", nil, func() any { return "a-object" })

	err := c.Resolve()
	if err != nil {
		t.Fatalf("got: %s, want: nil", err.Error())
	}
}

func TestResolve_BuilderIsNil(t *testing.T) {
	c := simpledi.NewContainer()

	c.Register("a", nil, nil)
	c.Register("b", []string{"a"}, func() any { return "b-object" })

	err := c.Resolve()
	if err == nil {
		t.Fatal("got: nil, want: error")
	}
	if !strings.Contains(err.Error(), "builder is nil") {
		t.Fatalf("got: %s, want: builder is nil", err.Error())
	}
}

func TestResolve_CyclicDependencies(t *testing.T) {
	c := simpledi.NewContainer()

	c.Register("a", []string{"b"}, func() any {
		b := c.Get("b").(string)
		return fmt.Sprintf("a-object: '%s'", b)
	})
	c.Register("b", []string{"a"}, func() any {
		a := c.Get("a").(string)
		return fmt.Sprintf("b-object: '%s'", a)
	})

	err := c.Resolve()
	if err == nil {
		t.Fatal("got: nil, want: error")
	}
	if !strings.Contains(err.Error(), "cyclic detected") {
		t.Fatalf("got: %s, want: cyclic detected", err.Error())
	}
}

func TestResolve_MissingDependency(t *testing.T) {
	c := simpledi.NewContainer()

	c.Register("a", nil, func() any { return "a-object" })
	c.Register("b", []string{"a", "d"}, func() any {
		a := c.Get("a").(string)
		c := c.Get("d").(string)
		return fmt.Sprintf("b-object: '%s' + '%s'", a, c)
	})

	err := c.Resolve()
	if err == nil {
		t.Fatal("got: nil, want: error")
	}
	if !strings.Contains(err.Error(), "not declared") {
		t.Fatalf("got: %s, want: not declared", err.Error())
	}
}

func TestGet_BeforeResolve(t *testing.T) {
	c := simpledi.NewContainer()

	a := c.Get("a")
	if a != nil {
		t.Errorf("got: %v, want: nil", a)
	}
}

func TestGet_AfterResolve(t *testing.T) {
	c := simpledi.NewContainer()

	c.Register("a", nil, func() any { return "a-object" })

	err := c.Resolve()
	if err != nil {
		t.Errorf("got: %s, want: nil", err.Error())
	}

	a := c.Get("a")
	if a == nil {
		t.Error("got: nil, want: a-object")
	}
}

func TestRegister_UsesLastBuilder(t *testing.T) {
	c := simpledi.NewContainer()

	c.Register("a", nil, func() any { return "a-object-1" })
	c.Register("a", nil, func() any { return "a-object-2" })
	c.Register("a", nil, func() any { return "a-object-3" })

	err := c.Resolve()
	if err != nil {
		t.Errorf("got: %s, want: nil", err.Error())
	}

	a := c.Get("a").(string)
	if a != "a-object-3" {
		t.Errorf("got: %s, want: a-object-3", a)
	}
}

func TestRegister_EmptyKey(t *testing.T) {
	c := simpledi.NewContainer()

	c.Register("", nil, func() any { return "-object" })
	c.Register("a", nil, func() any { return "a-object" })

	err := c.Resolve()
	if err != nil {
		t.Errorf("got: %s, want: nil", err.Error())
	}

	empty := c.Get("")
	if empty != "-object" {
		t.Errorf("got: %v, want: -object", empty)
	}

	a := c.Get("a")
	if a != "a-object" {
		t.Errorf("got: %v, want: a-object", a)
	}
}

func TestResolve_EmptyContainer(t *testing.T) {
	c := simpledi.NewContainer()

	err := c.Resolve()
	if err != nil {
		t.Errorf("got: %s, want: nil", err.Error())
	}
}
