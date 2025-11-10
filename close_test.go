package simpledi_test

import (
	"errors"
	"testing"

	"github.com/eerzho/simpledi"
)

func Test_Close_Without_Close_Functions(t *testing.T) {
	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "yeast",
			New: func() any {
				return "yeast"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "flour",
			New: func() any {
				return "flour"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "bread",
			Deps: []string{"yeast", "flour"},
			New: func() any {
				return "bread"
			},
		})
		simpledi.Resolve()
	})

	assertNoError(t, simpledi.Close)
}

func Test_Close_Error(t *testing.T) {
	order := make([]string, 0)
	someError1 := errors.New("some error 1")
	someError2 := errors.New("some error 2")
	someError3 := errors.New("some error 3")

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "yeast",
			New: func() any {
				return "yeast"
			},
			Close: func() error {
				order = append(order, "yeast")
				return someError1
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "flour",
			New: func() any {
				return "flour"
			},
			Close: func() error {
				order = append(order, "flour")
				return someError2
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "bread",
			Deps: []string{"yeast", "flour"},
			New: func() any {
				return "bread"
			},
			Close: func() error {
				order = append(order, "bread")
				return someError3
			},
		})
		simpledi.Resolve()
	})

	assertError(t, simpledi.Close, someError1, someError2, someError3)
	assertOrder(t, order, []string{"bread", "flour", "yeast"})
}

func Test_Close_Multiple_Times(t *testing.T) {
	someError := errors.New("some error")

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "yeast",
			New: func() any {
				return "yeast"
			},
			Close: func() error {
				return someError
			},
		})
		simpledi.Resolve()
	})

	assertError(t, simpledi.Close, someError)
	assertNoError(t, simpledi.Close)
}
