package simpledi_test

import (
	"testing"

	"github.com/eerzho/simpledi"
)

func Test_Set_Err_Container_Resolved(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Resolve()
	})

	assertPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "service_1",
			New: func() any {
				return &testServiceImpl1{}
			},
		})
	}, simpledi.ErrContainerResolved)
}

func Test_Set_Err_ID_Required(t *testing.T) {
	defer simpledi.Close()

	assertPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			New: func() any {
				return &testServiceImpl1{}
			},
		})
	}, simpledi.ErrIDRequired)
}

func Test_Set_Err_New_Required(t *testing.T) {
	defer simpledi.Close()

	assertPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "service_2",
		})
	}, simpledi.ErrNewRequired)
}
