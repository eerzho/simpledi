package simpledi_test

import (
	"testing"

	"github.com/eerzho/simpledi"
)

func Test_Get_Success(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "yeast",
		New: func() any {
			return "yeast"
		},
	})
	simpledi.Resolve()

	assertNoPanic(t, func() {
		_ = simpledi.Get[string]("yeast")
	})
}

func Test_Get_ID_Not_Found(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "yeast",
		New: func() any {
			return "yeast"
		},
	})
	simpledi.Resolve()

	assertPanic(t, func() {
		_ = simpledi.Get[string]("bread")
	}, simpledi.ErrIDNotFound)
}

func Test_Get_Type_Mismatch(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "yeast",
		New: func() any {
			return "yeast"
		},
	})
	simpledi.Resolve()

	assertPanic(t, func() {
		_ = simpledi.Get[int]("yeast")
	}, simpledi.ErrTypeMismatch)
}

func Test_Get_Container_Not_Resolved(t *testing.T) {
	defer simpledi.Close()

	assertPanic(t, func() {
		_ = simpledi.Get[string]("yeast")
	}, simpledi.ErrContainerNotResolved)
}
