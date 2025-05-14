package simpledi_test

import (
	"testing"

	"github.com/eerzho/simpledi"
)

func TestContainer_Register(t *testing.T) {
	c := simpledi.NewContainer()

	// Test valid registration
	c.Register("service1", nil, func() any { return "service1" })
	c.Register("service2", []string{"service1"}, func() any { return "service2" })

	// Test nil builder
	defer func() {
		if r := recover(); r == nil {
			t.Error("Register with nil builder should panic")
		}
	}()
	c.Register("service3", nil, nil)
}

func TestContainer_Get(t *testing.T) {
	c := simpledi.NewContainer()
	expected := "test"
	c.Register("test", nil, func() any { return expected })

	// Test before resolution
	defer func() {
		if r := recover(); r == nil {
			t.Error("Get before resolution should panic")
		}
	}()
	c.Get("test")

	// Test after resolution
	c.Resolve()
	if got := c.Get("test"); got != expected {
		t.Errorf("Get() = %v, want %v", got, expected)
	}

	// Test non-existent service
	defer func() {
		if r := recover(); r == nil {
			t.Error("Get non-existent service should panic")
		}
	}()
	c.Get("non-existent")
}

func TestContainer_Resolve(t *testing.T) {
	c := simpledi.NewContainer()

	// Test simple resolution
	c.Register("service1", nil, func() any { return "service1" })
	c.Register("service2", []string{"service1"}, func() any { return "service2" })
	c.Resolve()

	if got := c.Get("service1"); got != "service1" {
		t.Errorf("Get(service1) = %v, want %v", got, "service1")
	}
	if got := c.Get("service2"); got != "service2" {
		t.Errorf("Get(service2) = %v, want %v", got, "service2")
	}

	// Test circular dependency
	c = simpledi.NewContainer()
	c.Register("service1", []string{"service2"}, func() any { return "service1" })
	c.Register("service2", []string{"service1"}, func() any { return "service2" })

	defer func() {
		if r := recover(); r == nil {
			t.Error("Resolve with circular dependency should panic")
		}
	}()
	c.Resolve()
}

func TestContainer_RegisterDuplicate(t *testing.T) {
	c := simpledi.NewContainer()

	// First registration should succeed
	c.Register("service", nil, func() any { return "first" })

	// Second registration should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Register duplicate service should panic")
		}
	}()
	c.Register("service", nil, func() any { return "second" })
}

func TestContainer_NonExistentDependency(t *testing.T) {
	c := simpledi.NewContainer()

	// Register service with non-existent dependency
	c.Register("service", []string{"non-existent"}, func() any { return "service" })

	// Panic should occur during resolve
	defer func() {
		if r := recover(); r == nil {
			t.Error("Resolve with non-existent dependency should panic")
		}
	}()
	c.Resolve()
}

func TestContainer_ResolutionOrder(t *testing.T) {
	c := simpledi.NewContainer()

	// Create a chain of dependencies: service1 <- service2 <- service3
	order := make([]string, 0)

	c.Register("service1", nil, func() any {
		order = append(order, "service1")
		return "service1"
	})

	c.Register("service2", []string{"service1"}, func() any {
		order = append(order, "service2")
		return "service2"
	})

	c.Register("service3", []string{"service2"}, func() any {
		order = append(order, "service3")
		return "service3"
	})

	c.Resolve()

	// Verify resolution order
	expected := []string{"service1", "service2", "service3"}
	if len(order) != len(expected) {
		t.Errorf("Expected %d services to be resolved, got %d", len(expected), len(order))
	}
	for i, service := range expected {
		if order[i] != service {
			t.Errorf("Expected service %s at position %d, got %s", service, i, order[i])
		}
	}
}

func TestContainer_GetBeforeResolve(t *testing.T) {
	c := simpledi.NewContainer()

	// Register a service
	c.Register("service", nil, func() any { return "service" })

	// Try to get service before resolve
	defer func() {
		if r := recover(); r == nil {
			t.Error("Get before resolve should panic")
		}
	}()
	c.Get("service")
}

func TestContainer_EmptyNeeds(t *testing.T) {
	c := simpledi.NewContainer()

	// Test with empty needs slice
	c.Register("service", []string{}, func() any { return "service" })
	c.Resolve()

	if got := c.Get("service"); got != "service" {
		t.Errorf("Get() = %v, want %v", got, "service")
	}
}

func TestContainer_MultipleDependencies(t *testing.T) {
	c := simpledi.NewContainer()

	// Create a service with multiple dependencies
	c.Register("dep1", nil, func() any { return "dep1" })
	c.Register("dep2", nil, func() any { return "dep2" })
	c.Register("service", []string{"dep1", "dep2"}, func() any { return "service" })

	c.Resolve()

	if got := c.Get("service"); got != "service" {
		t.Errorf("Get() = %v, want %v", got, "service")
	}
}
