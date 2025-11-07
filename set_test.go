package simpledi_test

import (
	"testing"

	"github.com/eerzho/simpledi"
)

func Test_Set_Success(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "yeast",
			New: func() any {
				return "yeast"
			},
		})
	})
}

func Test_Set_ID_Required(t *testing.T) {
	defer simpledi.Close()

	assertPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			New: func() any {
				return "yeast"
			},
		})
	}, simpledi.ErrIDRequired)
}

func Test_Set_New_Required(t *testing.T) {
	defer simpledi.Close()

	assertPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "yeast",
		})
	}, simpledi.ErrNewRequired)
}

func Test_Set_Container_Resolved(t *testing.T) {
	defer simpledi.Close()
	simpledi.Resolve()

	assertPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "yeast",
			New: func() any {
				return "yeast"
			},
		})
	}, simpledi.ErrContainerResolved)
}
