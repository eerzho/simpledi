package simpledi_test

import (
	"testing"

	"github.com/eerzho/simpledi"
)

func TestContainer_RegisterAndGet(t *testing.T) {
	c := simpledi.NewContainer()
	c.Register("A", nil, func() any { return "InstanceA" })
	err := c.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	instance := c.Get("A")
	if instance != "InstanceA" {
		t.Errorf("expected 'InstanceA', got %v", instance)
	}
}

func TestContainer_ResolveWithDependencies(t *testing.T) {
	c := simpledi.NewContainer()
	c.Register("A", nil, func() any { return "A" })
	c.Register("B", []string{"A"}, func() any { return "B" })

	err := c.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.Get("A") != "A" || c.Get("B") != "B" {
		t.Errorf("instances not resolved correctly: A=%v, B=%v", c.Get("A"), c.Get("B"))
	}
}

func TestContainer_MissingDependency(t *testing.T) {
	c := simpledi.NewContainer()
	c.Register("B", []string{"A"}, func() any { return "B" })

	err := c.Resolve()
	if err == nil || err.Error() == "" {
		t.Errorf("expected missing dependency error, got %v", err)
	}
}

func TestContainer_CyclicDependency(t *testing.T) {
	c := simpledi.NewContainer()
	c.Register("A", []string{"B"}, func() any { return "A" })
	c.Register("B", []string{"A"}, func() any { return "B" })

	err := c.Resolve()
	if err == nil || err.Error() == "" {
		t.Errorf("expected cyclic dependency error, got %v", err)
	}
}

func TestContainer_ConstructorMissing(t *testing.T) {
	c := simpledi.NewContainer()
	c.Register("A", nil, nil)

	err := c.Resolve()
	if err == nil || err.Error() == "" {
		t.Errorf("expected missing constructor error, got %v", err)
	}
}
