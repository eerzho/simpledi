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

func Test_Get_Empty_String_ID(t *testing.T) {
	defer simpledi.Close()
	simpledi.Resolve()

	assertPanic(t, func() {
		simpledi.Get[string]("")
	}, simpledi.ErrIDRequired)
}

func Test_Get_Multiple_Types_Same_ID(t *testing.T) {

}

func Test_Get_After_Close(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "yeast",
		New: func() any {
			return "yeast"
		},
	})
	simpledi.Resolve()
	simpledi.Close()

	assertPanic(t, func() {
		simpledi.Get[string]("yeast")
	}, simpledi.ErrContainerNotResolved)
}

func Test_Get_Same_Instance_Returned(t *testing.T) {
	defer simpledi.Close()
	type service struct {
		data string
	}
	simpledi.Set(simpledi.Definition{
		ID: "service",
		New: func() any {
			return &service{data: "some data"}
		},
	})
	simpledi.Resolve()

	first := simpledi.Get[*service]("service")
	second := simpledi.Get[*service]("service")
	assertSameInstance(t, first, second)
}

func Test_Get_Nil_Value(t *testing.T) {

}
