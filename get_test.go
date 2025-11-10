package simpledi_test

import (
	"testing"

	"github.com/eerzho/simpledi"
)

func Test_Get_Err_ID_Required(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "service_1",
			New: func() any {
				return &testServiceImpl1{}
			},
		})
		simpledi.Resolve()
	})

	assertPanic(t, func() {
		_ = simpledi.Get[*testServiceImpl1]("")
	}, simpledi.ErrIDRequired)
}

func Test_Get_Err_ID_NotFound(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "service_1",
			New: func() any {
				return &testServiceImpl1{}
			},
		})
		simpledi.Resolve()
	})

	assertPanic(t, func() {
		_ = simpledi.Get[*testServiceImpl1]("not_found")
	}, simpledi.ErrIDNotFound)
}

func Test_Get_Err_Type_Mismatch(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "service_1",
			New: func() any {
				return &testServiceImpl1{}
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "service_2",
			New: func() any {
				return &testServiceImpl2{}
			},
		})
		simpledi.Resolve()
	})

	assertPanic(t, func() {
		_ = simpledi.Get[*testServiceImpl2]("service_1")
	}, simpledi.ErrTypeMismatch)
}

func Test_Get_Generic_Type_With_Interface(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "service_1",
			New: func() any {
				return &testServiceImpl1{}
			},
		})
		simpledi.Resolve()
	})

	assertNoPanic(t, func() {
		_ = simpledi.Get[testService1]("service_1")
	})
}

func Test_Get_Returns_Same_Instance_Value(t *testing.T) {
	defer simpledi.Close()
	someData := "some data"

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "service_2",
			New: func() any {
				return &testServiceImpl2{data: someData}
			},
		})
		simpledi.Resolve()
	})

	assertNoPanic(t, func() {
		v1 := simpledi.Get[*testServiceImpl2]("service_2")
		v2 := simpledi.Get[*testServiceImpl2]("service_2")
		assertSameValue(t, v1.data, v2.data)
	})
}

func Test_Get_Returns_Same_Instance_Pointer(t *testing.T) {
	defer simpledi.Close()
	someData := "someData"

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "service_2",
			New: func() any {
				return &testServiceImpl2{data: someData}
			},
		})
		simpledi.Resolve()
	})

	assertNoPanic(t, func() {
		v1 := simpledi.Get[*testServiceImpl2]("service_2")
		v2 := simpledi.Get[*testServiceImpl2]("service_2")
		assertSamePointer(t, v1, v2)
	})
}

func Test_Get_Err_Before_Resolve(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "service_1",
			New: func() any {
				return &testServiceImpl1{}
			},
		})
	})

	assertPanic(t, func() {
		_ = simpledi.Get[*testServiceImpl1]("service_1")
	}, simpledi.ErrIDNotFound)
}

func Test_Get_With_Value_Type(t *testing.T) {
	defer simpledi.Close()
	someData := "some data"

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "service_2",
			New: func() any {
				return testServiceImpl2{data: someData}
			},
		})
		simpledi.Resolve()
	})

	assertNoPanic(t, func() {
		service2 := simpledi.Get[testServiceImpl2]("service_2")
		assertSameValue(t, service2.data, someData)
	})
}

func Test_Get_With_Nil_Value(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "nil_val",
			New: func() any {
				return nil
			},
		})
		simpledi.Resolve()
	})

	assertNoPanic(t, func() {
		nilVal := simpledi.Get[any]("nil_val")
		assertSameValue(t, nilVal, nil)
	})
}

func Test_Get_From_New(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "service_1",
			New: func() any {
				return &testServiceImpl1{}
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "service_3",
			Deps: []string{"service_1"},
			New: func() any {
				service1 := simpledi.Get[*testServiceImpl1]("service_1")
				return &testServiceImpl3{service1: service1}
			},
		})
		simpledi.Resolve()
	})

	assertNoPanic(t, func() {
		service1 := simpledi.Get[*testServiceImpl1]("service_1")
		service3 := simpledi.Get[*testServiceImpl3]("service_3")
		assertSamePointer(t, service1, service3.service1)
	})
}

func Test_Get_From_New_Err_ID_Not_Found(t *testing.T) {
	defer simpledi.Close()

	assertPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "service_3",
			New: func() any {
				simpledi.Get[*testServiceImpl1]("service_1")
				return &testServiceImpl3{}
			},
		})
		simpledi.Resolve()
	}, simpledi.ErrIDNotFound)
}

func Test_Get_From_New_Err_Type_Mismatch(t *testing.T) {
	defer simpledi.Close()

	assertPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "service_1",
			New: func() any {
				return &testServiceImpl1{}
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "service_3",
			New: func() any {
				_ = simpledi.Get[*testServiceImpl2]("service_1")
				return &testServiceImpl3{}
			},
		})
		simpledi.Resolve()
	}, simpledi.ErrTypeMismatch)
}
