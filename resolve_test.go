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

func Test_Resolve_Panicking_New_Function(t *testing.T) {

}
