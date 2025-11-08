package simpledi_test

import (
	"testing"

	"github.com/eerzho/simpledi"
)

func Test_Resolve_Container_Resolved(t *testing.T) {
	defer simpledi.Close()
	simpledi.Resolve()

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrContainerResolved)
}

func Test_Resolve_ID_Duplicate(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "yeast",
		New: func() any {
			return "yeast"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "yeast",
		New: func() any {
			return "yeast"
		},
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrIDDuplicate)
}

func Test_Resolve_Dependency_Not_Found(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID:   "yeast",
		Deps: []string{"bread"},
		New: func() any {
			return "yeast"
		},
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyNotFound)
}

func Test_Resolve_Dependency_Cycle(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID:   "yeast",
		Deps: []string{"bread"},
		New: func() any {
			return "yeast"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "bread",
		Deps: []string{"yeast"},
		New: func() any {
			return "bread"
		},
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyCycle)
}

func Test_Resolve_Empty_Container(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Resolve()
	})
}

func Test_Resolve_Execution_Order(t *testing.T) {
	defer simpledi.Close()

	var executionOrder []string
	simpledi.Set(simpledi.Definition{
		ID: "base",
		New: func() any {
			executionOrder = append(executionOrder, "base")
			return "base"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "level1",
		Deps: []string{"base"},
		New: func() any {
			executionOrder = append(executionOrder, "level1")
			return "level1"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "level2",
		Deps: []string{"level1"},
		New: func() any {
			executionOrder = append(executionOrder, "level2")
			return "level2"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "level3",
		Deps: []string{"level2"},
		New: func() any {
			executionOrder = append(executionOrder, "level3")
			return "level3"
		},
	})
	simpledi.Resolve()

	if len(executionOrder) != 4 {
		t.Errorf("got: %d executions, want: 4", len(executionOrder))
		return
	}
	if executionOrder[0] != "base" {
		t.Errorf("got: %s, want: base (first)", executionOrder[0])
	}
	if executionOrder[1] != "level1" {
		t.Errorf("got: %s, want: level1 (second)", executionOrder[1])
	}
	if executionOrder[2] != "level2" {
		t.Errorf("got: %s, want: level2 (third)", executionOrder[2])
	}
	if executionOrder[3] != "level3" {
		t.Errorf("got: %s, want: level3 (fourth)", executionOrder[3])
	}
}

func Test_Resolve_After_Close_And_Reset(t *testing.T) {
	defer simpledi.Close()

	simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "yeast",
		New: func() any {
			return "yeast"
		},
	})
	assertNoPanic(t, func() {
		simpledi.Resolve()
	})
}

func Test_Resolve_No_Dependencies(t *testing.T) {
	defer simpledi.Close()

	simpledi.Set(simpledi.Definition{
		ID: "flour",
		New: func() any {
			return "flour"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "water",
		New: func() any {
			return "water"
		},
	})

	assertNoPanic(t, func() {
		simpledi.Resolve()
	})
}

func Test_Multiple_Resolve_Calls_After_Close(t *testing.T) {
	defer simpledi.Close()

	simpledi.Close()
	simpledi.Resolve()
	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrContainerResolved)
}

func Test_Circular_Dependency_Three_Way_Cycle(t *testing.T) {
	defer simpledi.Close()

	simpledi.Set(simpledi.Definition{
		ID:   "A",
		Deps: []string{"B"},
		New: func() any {
			return "A"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "B",
		Deps: []string{"C"},
		New: func() any {
			return "B"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "C",
		Deps: []string{"A"},
		New: func() any {
			return "C"
		},
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyCycle)
}
