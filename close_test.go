package simpledi_test

import (
	"errors"
	"testing"

	"github.com/eerzho/simpledi"
)

func Test_Close_Some_Err(t *testing.T) {
	someError := errors.New("some error")
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

	assertError(t, simpledi.Close, someError)
}

func Test_Close_Empty_Container(t *testing.T) {
	assertNoError(t, simpledi.Close)
}

func Test_Close_Without_Resolve(t *testing.T) {

}

func Test_Close_Multiple_Errors(t *testing.T) {
	someError1 := errors.New("some error 1")
	simpledi.Set(simpledi.Definition{
		ID: "yeast 1",
		New: func() any {
			return "yeast 1"
		},
		Close: func() error {
			return someError1
		},
	})
	someError2 := errors.New("some error 2")
	simpledi.Set(simpledi.Definition{
		ID: "yeast 2",
		New: func() any {
			return "yeast 2"
		},
		Close: func() error {
			return someError2
		},
	})
	simpledi.Resolve()

	assertError(t, simpledi.Close, someError1, someError2)
}

func Test_Close_Mixed_Success_And_Error(t *testing.T) {
	someError := errors.New("some error")
	simpledi.Set(simpledi.Definition{
		ID: "yeast 1",
		New: func() any {
			return "yeast 1"
		},
		Close: func() error {
			return someError
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "yeast 2",
		New: func() any {
			return "yeast 2"
		},
		Close: func() error {
			return nil
		},
	})
	simpledi.Resolve()

	assertError(t, simpledi.Close, someError)
}

func Test_Close_Reverse_Order(t *testing.T) {
}

func Test_Close_Nil_Close_Function(t *testing.T) {
	simpledi.Set(simpledi.Definition{
		ID: "yeast",
		New: func() any {
			return "yeast"
		},
		Close: nil,
	})
	simpledi.Resolve()

	assertNoError(t, simpledi.Close)
}

func Test_Close_State_Reset(t *testing.T) {
}

func Test_Close_Multiple_Times(t *testing.T) {
	simpledi.Set(simpledi.Definition{
		ID: "yeast",
		New: func() any {
			return "yeast"
		},
	})
	simpledi.Resolve()

	assertNoError(t, simpledi.Close)
	assertNoError(t, simpledi.Close)
}
