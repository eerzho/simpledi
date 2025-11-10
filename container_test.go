package simpledi_test

import (
	"errors"
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

func Test_Resolve_Single_Recipe_All_Ingredients_Available(t *testing.T) {
	defer simpledi.Close()
	order := make([]string, 0)

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "yeast",
			New: func() any {
				order = append(order, "yeast")
				return "yeast"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "flour",
			New: func() any {
				order = append(order, "flour")
				return "flour"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "bread",
			Deps: []string{"yeast", "flour"},
			New: func() any {
				order = append(order, "bread")
				return "bread"
			},
		})
		simpledi.Resolve()
	})

	assertOrder(t, order, []string{"yeast", "flour", "bread"})

	assertNoPanic(t, func() {
		assertSameValue(t, simpledi.Get[string]("yeast"), "yeast")
		assertSameValue(t, simpledi.Get[string]("flour"), "flour")
		assertSameValue(t, simpledi.Get[string]("bread"), "bread")
	})
}

func Test_Resolve_Single_Recipe_Missing_One_Ingredient(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID:  "yeast",
			New: func() any { return "yeast" },
		})
		simpledi.Set(simpledi.Definition{
			ID:   "bread",
			Deps: []string{"yeast", "flour"},
			New:  func() any { return "bread" },
		})
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyNotFound)
}

func Test_Resolve_Chain_Two_Levels_Complete(t *testing.T) {
	defer simpledi.Close()
	order := make([]string, 0)

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "yeast",
			New: func() any {
				order = append(order, "yeast")
				return "yeast"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "flour",
			New: func() any {
				order = append(order, "flour")
				return "flour"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "meat",
			New: func() any {
				order = append(order, "meat")
				return "meat"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "bread",
			Deps: []string{"yeast", "flour"},
			New: func() any {
				order = append(order, "bread")
				return "bread"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "sandwich",
			Deps: []string{"bread", "meat"},
			New: func() any {
				order = append(order, "sandwich")
				return "sandwich"
			},
		})
		simpledi.Resolve()
	})

	assertOrder(t, order, []string{"yeast", "flour", "meat", "bread", "sandwich"})

	assertNoPanic(t, func() {
		assertSameValue(t, simpledi.Get[string]("yeast"), "yeast")
		assertSameValue(t, simpledi.Get[string]("flour"), "flour")
		assertSameValue(t, simpledi.Get[string]("meat"), "meat")
		assertSameValue(t, simpledi.Get[string]("bread"), "bread")
		assertSameValue(t, simpledi.Get[string]("sandwich"), "sandwich")
	})
}

func Test_Resolve_Chain_Two_Levels_Broken(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID:   "bread",
			Deps: []string{"yeast", "flour"},
			New:  func() any { return "bread" },
		})
		simpledi.Set(simpledi.Definition{
			ID:   "sandwich",
			Deps: []string{"bread", "meat"},
			New:  func() any { return "sandwich" },
		})
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyNotFound)
}

func Test_Resolve_Multiple_Independent_All_Available(t *testing.T) {
	defer simpledi.Close()
	order := make([]string, 0)

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "yeast",
			New: func() any {
				order = append(order, "yeast")
				return "yeast"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "flour",
			New: func() any {
				order = append(order, "flour")
				return "flour"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "orange",
			New: func() any {
				order = append(order, "orange")
				return "orange"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "water",
			New: func() any {
				order = append(order, "water")
				return "water"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "tomato",
			New: func() any {
				order = append(order, "tomato")
				return "tomato"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "cucumber",
			New: func() any {
				order = append(order, "cucumber")
				return "cucumber"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "bread",
			Deps: []string{"yeast", "flour"},
			New: func() any {
				order = append(order, "bread")
				return "bread"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "juice",
			Deps: []string{"orange", "water"},
			New: func() any {
				order = append(order, "juice")
				return "juice"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "salad",
			Deps: []string{"tomato", "cucumber"},
			New: func() any {
				order = append(order, "salad")
				return "salad"
			},
		})
		simpledi.Resolve()
	})

	assertOrder(t, order, []string{"yeast", "flour", "orange", "water", "tomato", "cucumber", "bread", "juice", "salad"})

	assertNoPanic(t, func() {
		assertSameValue(t, simpledi.Get[string]("yeast"), "yeast")
		assertSameValue(t, simpledi.Get[string]("flour"), "flour")
		assertSameValue(t, simpledi.Get[string]("orange"), "orange")
		assertSameValue(t, simpledi.Get[string]("water"), "water")
		assertSameValue(t, simpledi.Get[string]("tomato"), "tomato")
		assertSameValue(t, simpledi.Get[string]("cucumber"), "cucumber")
		assertSameValue(t, simpledi.Get[string]("bread"), "bread")
		assertSameValue(t, simpledi.Get[string]("juice"), "juice")
		assertSameValue(t, simpledi.Get[string]("salad"), "salad")
	})
}

func Test_Resolve_Multiple_Independent_Partial_Available(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID:  "yeast",
			New: func() any { return "yeast" },
		})
		simpledi.Set(simpledi.Definition{
			ID:  "flour",
			New: func() any { return "flour" },
		})
		simpledi.Set(simpledi.Definition{
			ID:   "bread",
			Deps: []string{"yeast", "flour"},
			New:  func() any { return "bread" },
		})
		simpledi.Set(simpledi.Definition{
			ID:   "juice",
			Deps: []string{"orange", "water"},
			New:  func() any { return "juice" },
		})
		simpledi.Set(simpledi.Definition{
			ID:  "tomato",
			New: func() any { return "tomato" },
		})
		simpledi.Set(simpledi.Definition{
			ID:  "cucumber",
			New: func() any { return "cucumber" },
		})
		simpledi.Set(simpledi.Definition{
			ID:   "salad",
			Deps: []string{"tomato", "cucumber"},
			New:  func() any { return "salad" },
		})
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyNotFound)
}

func Test_Resolve_Empty_Supplies_List(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID:   "bread",
			Deps: []string{"yeast", "flour"},
			New:  func() any { return "bread" },
		})
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyNotFound)
}

func Test_Resolve_Empty_Recipes_List(t *testing.T) {
	defer simpledi.Close()
	order := make([]string, 0)

	assertNoPanic(t, func() {
		simpledi.Resolve()
	})

	assertOrder(t, order, []string{})
}

func Test_Resolve_Circular_Dependency_Two_Recipes(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID:  "salt",
			New: func() any { return "salt" },
		})
		simpledi.Set(simpledi.Definition{
			ID:  "pepper",
			New: func() any { return "pepper" },
		})
		simpledi.Set(simpledi.Definition{
			ID:   "dish_a",
			Deps: []string{"dish_b", "salt"},
			New:  func() any { return "dish_a" },
		})
		simpledi.Set(simpledi.Definition{
			ID:   "dish_b",
			Deps: []string{"dish_a", "pepper"},
			New:  func() any { return "dish_b" },
		})
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyCycle)
}

func Test_Resolve_Chain_Three_Levels_Complete(t *testing.T) {
	defer simpledi.Close()
	order := make([]string, 0)

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "flour",
			New: func() any {
				order = append(order, "flour")
				return "flour"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "sugar",
			New: func() any {
				order = append(order, "sugar")
				return "sugar"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "milk",
			New: func() any {
				order = append(order, "milk")
				return "milk"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "dish_a",
			Deps: []string{"flour"},
			New: func() any {
				order = append(order, "dish_a")
				return "dish_a"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "dish_b",
			Deps: []string{"dish_a", "sugar"},
			New: func() any {
				order = append(order, "dish_b")
				return "dish_b"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "dish_c",
			Deps: []string{"dish_b", "milk"},
			New: func() any {
				order = append(order, "dish_c")
				return "dish_c"
			},
		})
		simpledi.Resolve()
	})

	assertOrder(t, order, []string{"flour", "sugar", "milk", "dish_a", "dish_b", "dish_c"})

	assertNoPanic(t, func() {
		assertSameValue(t, simpledi.Get[string]("flour"), "flour")
		assertSameValue(t, simpledi.Get[string]("sugar"), "sugar")
		assertSameValue(t, simpledi.Get[string]("milk"), "milk")
		assertSameValue(t, simpledi.Get[string]("dish_a"), "dish_a")
		assertSameValue(t, simpledi.Get[string]("dish_b"), "dish_b")
		assertSameValue(t, simpledi.Get[string]("dish_c"), "dish_c")
	})
}

func Test_Resolve_Diamond_Dependency_Pattern(t *testing.T) {
	defer simpledi.Close()
	order := make([]string, 0)

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "flour",
			New: func() any {
				order = append(order, "flour")
				return "flour"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "sugar",
			New: func() any {
				order = append(order, "sugar")
				return "sugar"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "salt",
			New: func() any {
				order = append(order, "salt")
				return "salt"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "dish_a",
			Deps: []string{"flour"},
			New: func() any {
				order = append(order, "dish_a")
				return "dish_a"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "dish_b",
			Deps: []string{"dish_a", "sugar"},
			New: func() any {
				order = append(order, "dish_b")
				return "dish_b"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "dish_c",
			Deps: []string{"dish_a", "salt"},
			New: func() any {
				order = append(order, "dish_c")
				return "dish_c"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "dish_d",
			Deps: []string{"dish_b", "dish_c"},
			New: func() any {
				order = append(order, "dish_d")
				return "dish_d"
			},
		})
		simpledi.Resolve()
	})

	assertOrder(t, order, []string{"flour", "sugar", "salt", "dish_a", "dish_b", "dish_c", "dish_d"})

	assertNoPanic(t, func() {
		assertSameValue(t, simpledi.Get[string]("flour"), "flour")
		assertSameValue(t, simpledi.Get[string]("sugar"), "sugar")
		assertSameValue(t, simpledi.Get[string]("salt"), "salt")
		assertSameValue(t, simpledi.Get[string]("dish_a"), "dish_a")
		assertSameValue(t, simpledi.Get[string]("dish_b"), "dish_b")
		assertSameValue(t, simpledi.Get[string]("dish_c"), "dish_c")
		assertSameValue(t, simpledi.Get[string]("dish_d"), "dish_d")
	})
}

func Test_Resolve_Same_Ingredient_Multiple_Recipes(t *testing.T) {
	defer simpledi.Close()
	order := make([]string, 0)

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "flour",
			New: func() any {
				order = append(order, "flour")
				return "flour"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "water",
			New: func() any {
				order = append(order, "water")
				return "water"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "tomato",
			New: func() any {
				order = append(order, "tomato")
				return "tomato"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "egg",
			New: func() any {
				order = append(order, "egg")
				return "egg"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "bread",
			Deps: []string{"flour", "water"},
			New: func() any {
				order = append(order, "bread")
				return "bread"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "pizza",
			Deps: []string{"flour", "tomato"},
			New: func() any {
				order = append(order, "pizza")
				return "pizza"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "pasta",
			Deps: []string{"flour", "egg"},
			New: func() any {
				order = append(order, "pasta")
				return "pasta"
			},
		})
		simpledi.Resolve()
	})

	assertOrder(t, order, []string{"flour", "water", "tomato", "egg", "bread", "pizza", "pasta"})

	assertNoPanic(t, func() {
		assertSameValue(t, simpledi.Get[string]("flour"), "flour")
		assertSameValue(t, simpledi.Get[string]("water"), "water")
		assertSameValue(t, simpledi.Get[string]("tomato"), "tomato")
		assertSameValue(t, simpledi.Get[string]("egg"), "egg")
		assertSameValue(t, simpledi.Get[string]("bread"), "bread")
		assertSameValue(t, simpledi.Get[string]("pizza"), "pizza")
		assertSameValue(t, simpledi.Get[string]("pasta"), "pasta")
	})
}

func Test_Resolve_Recipe_Used_Multiple_Times_In_Another(t *testing.T) {
	defer simpledi.Close()
	order := make([]string, 0)

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "flour",
			New: func() any {
				order = append(order, "flour")
				return "flour"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "meat",
			New: func() any {
				order = append(order, "meat")
				return "meat"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "bread",
			Deps: []string{"flour"},
			New: func() any {
				order = append(order, "bread")
				return "bread"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "double_sandwich",
			Deps: []string{"bread", "bread", "meat"},
			New: func() any {
				order = append(order, "double_sandwich")
				return "double_sandwich"
			},
		})
		simpledi.Resolve()
	})

	assertOrder(t, order, []string{"flour", "meat", "bread", "double_sandwich"})

	assertNoPanic(t, func() {
		assertSameValue(t, simpledi.Get[string]("flour"), "flour")
		assertSameValue(t, simpledi.Get[string]("meat"), "meat")
		assertSameValue(t, simpledi.Get[string]("bread"), "bread")
		assertSameValue(t, simpledi.Get[string]("double_sandwich"), "double_sandwich")
	})
}

func Test_Resolve_Extra_Unused_Supplies_Present(t *testing.T) {
	defer simpledi.Close()
	order := make([]string, 0)

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "flour",
			New: func() any {
				order = append(order, "flour")
				return "flour"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "water",
			New: func() any {
				order = append(order, "water")
				return "water"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:  "sugar",
			New: func() any { return "sugar" },
		})
		simpledi.Set(simpledi.Definition{
			ID:  "salt",
			New: func() any { return "salt" },
		})
		simpledi.Set(simpledi.Definition{
			ID:  "pepper",
			New: func() any { return "pepper" },
		})
		simpledi.Set(simpledi.Definition{
			ID:  "cheese",
			New: func() any { return "cheese" },
		})
		simpledi.Set(simpledi.Definition{
			ID:   "bread",
			Deps: []string{"flour", "water"},
			New: func() any {
				order = append(order, "bread")
				return "bread"
			},
		})
		simpledi.Resolve()
	})

	assertOrder(t, order, []string{"flour", "water", "bread"})

	assertNoPanic(t, func() {
		assertSameValue(t, simpledi.Get[string]("flour"), "flour")
		assertSameValue(t, simpledi.Get[string]("water"), "water")
		assertSameValue(t, simpledi.Get[string]("sugar"), "sugar")
		assertSameValue(t, simpledi.Get[string]("salt"), "salt")
		assertSameValue(t, simpledi.Get[string]("pepper"), "pepper")
		assertSameValue(t, simpledi.Get[string]("cheese"), "cheese")
		assertSameValue(t, simpledi.Get[string]("bread"), "bread")
	})
}

func Test_Resolve_All_Recipes_Missing_Common_Ingredient(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID:  "flour",
			New: func() any { return "flour" },
		})
		simpledi.Set(simpledi.Definition{
			ID:  "sugar",
			New: func() any { return "sugar" },
		})
		simpledi.Set(simpledi.Definition{
			ID:  "water",
			New: func() any { return "water" },
		})
		simpledi.Set(simpledi.Definition{
			ID:   "dish_a",
			Deps: []string{"salt", "flour"},
			New:  func() any { return "dish_a" },
		})
		simpledi.Set(simpledi.Definition{
			ID:   "dish_b",
			Deps: []string{"salt", "sugar"},
			New:  func() any { return "dish_b" },
		})
		simpledi.Set(simpledi.Definition{
			ID:   "dish_c",
			Deps: []string{"salt", "water"},
			New:  func() any { return "dish_c" },
		})
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyNotFound)
}

func Test_Resolve_Chain_Five_Levels_Complete(t *testing.T) {
	defer simpledi.Close()
	order := make([]string, 0)

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "base",
			New: func() any {
				order = append(order, "base")
				return "base"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "ing1",
			New: func() any {
				order = append(order, "ing1")
				return "ing1"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "ing2",
			New: func() any {
				order = append(order, "ing2")
				return "ing2"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "ing3",
			New: func() any {
				order = append(order, "ing3")
				return "ing3"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "ing4",
			New: func() any {
				order = append(order, "ing4")
				return "ing4"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "level1",
			Deps: []string{"base"},
			New: func() any {
				order = append(order, "level1")
				return "level1"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "level2",
			Deps: []string{"level1", "ing1"},
			New: func() any {
				order = append(order, "level2")
				return "level2"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "level3",
			Deps: []string{"level2", "ing2"},
			New: func() any {
				order = append(order, "level3")
				return "level3"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "level4",
			Deps: []string{"level3", "ing3"},
			New: func() any {
				order = append(order, "level4")
				return "level4"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "level5",
			Deps: []string{"level4", "ing4"},
			New: func() any {
				order = append(order, "level5")
				return "level5"
			},
		})
		simpledi.Resolve()
	})

	assertOrder(t, order, []string{"base", "ing1", "ing2", "ing3", "ing4", "level1", "level2", "level3", "level4", "level5"})

	assertNoPanic(t, func() {
		assertSameValue(t, simpledi.Get[string]("base"), "base")
		assertSameValue(t, simpledi.Get[string]("ing1"), "ing1")
		assertSameValue(t, simpledi.Get[string]("ing2"), "ing2")
		assertSameValue(t, simpledi.Get[string]("ing3"), "ing3")
		assertSameValue(t, simpledi.Get[string]("ing4"), "ing4")
		assertSameValue(t, simpledi.Get[string]("level1"), "level1")
		assertSameValue(t, simpledi.Get[string]("level2"), "level2")
		assertSameValue(t, simpledi.Get[string]("level3"), "level3")
		assertSameValue(t, simpledi.Get[string]("level4"), "level4")
		assertSameValue(t, simpledi.Get[string]("level5"), "level5")
	})
}

func Test_Resolve_Two_Parallel_Independent_Chains(t *testing.T) {
	defer simpledi.Close()
	order := make([]string, 0)

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "flour",
			New: func() any {
				order = append(order, "flour")
				return "flour"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "water",
			New: func() any {
				order = append(order, "water")
				return "water"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "sugar",
			New: func() any {
				order = append(order, "sugar")
				return "sugar"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "milk",
			New: func() any {
				order = append(order, "milk")
				return "milk"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "chain1_a",
			Deps: []string{"flour"},
			New: func() any {
				order = append(order, "chain1_a")
				return "chain1_a"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "chain1_b",
			Deps: []string{"chain1_a", "water"},
			New: func() any {
				order = append(order, "chain1_b")
				return "chain1_b"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "chain2_a",
			Deps: []string{"sugar"},
			New: func() any {
				order = append(order, "chain2_a")
				return "chain2_a"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "chain2_b",
			Deps: []string{"chain2_a", "milk"},
			New: func() any {
				order = append(order, "chain2_b")
				return "chain2_b"
			},
		})
		simpledi.Resolve()
	})

	assertOrder(t, order, []string{"flour", "water", "sugar", "milk", "chain1_a", "chain2_a", "chain1_b", "chain2_b"})

	assertNoPanic(t, func() {
		assertSameValue(t, simpledi.Get[string]("flour"), "flour")
		assertSameValue(t, simpledi.Get[string]("water"), "water")
		assertSameValue(t, simpledi.Get[string]("sugar"), "sugar")
		assertSameValue(t, simpledi.Get[string]("milk"), "milk")
		assertSameValue(t, simpledi.Get[string]("chain1_a"), "chain1_a")
		assertSameValue(t, simpledi.Get[string]("chain1_b"), "chain1_b")
		assertSameValue(t, simpledi.Get[string]("chain2_a"), "chain2_a")
		assertSameValue(t, simpledi.Get[string]("chain2_b"), "chain2_b")
	})
}

func Test_Resolve_Recipe_Requires_Single_Ingredient_Only(t *testing.T) {
	defer simpledi.Close()
	order := make([]string, 0)

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "orange",
			New: func() any {
				order = append(order, "orange")
				return "orange"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "juice",
			Deps: []string{"orange"},
			New: func() any {
				order = append(order, "juice")
				return "juice"
			},
		})
		simpledi.Resolve()
	})

	assertOrder(t, order, []string{"orange", "juice"})

	assertNoPanic(t, func() {
		assertSameValue(t, simpledi.Get[string]("orange"), "orange")
		assertSameValue(t, simpledi.Get[string]("juice"), "juice")
	})
}

func Test_Resolve_Recipe_Requires_Many_Ingredients(t *testing.T) {
	defer simpledi.Close()
	order := make([]string, 0)

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "ing1",
			New: func() any {
				order = append(order, "ing1")
				return "ing1"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "ing2",
			New: func() any {
				order = append(order, "ing2")
				return "ing2"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "ing3",
			New: func() any {
				order = append(order, "ing3")
				return "ing3"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "ing4",
			New: func() any {
				order = append(order, "ing4")
				return "ing4"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "ing5",
			New: func() any {
				order = append(order, "ing5")
				return "ing5"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "ing6",
			New: func() any {
				order = append(order, "ing6")
				return "ing6"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "ing7",
			New: func() any {
				order = append(order, "ing7")
				return "ing7"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "ing8",
			New: func() any {
				order = append(order, "ing8")
				return "ing8"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "complex_dish",
			Deps: []string{"ing1", "ing2", "ing3", "ing4", "ing5", "ing6", "ing7", "ing8"},
			New: func() any {
				order = append(order, "complex_dish")
				return "complex_dish"
			},
		})
		simpledi.Resolve()
	})

	assertOrder(t, order, []string{"ing1", "ing2", "ing3", "ing4", "ing5", "ing6", "ing7", "ing8", "complex_dish"})

	assertNoPanic(t, func() {
		assertSameValue(t, simpledi.Get[string]("ing1"), "ing1")
		assertSameValue(t, simpledi.Get[string]("ing2"), "ing2")
		assertSameValue(t, simpledi.Get[string]("ing3"), "ing3")
		assertSameValue(t, simpledi.Get[string]("ing4"), "ing4")
		assertSameValue(t, simpledi.Get[string]("ing5"), "ing5")
		assertSameValue(t, simpledi.Get[string]("ing6"), "ing6")
		assertSameValue(t, simpledi.Get[string]("ing7"), "ing7")
		assertSameValue(t, simpledi.Get[string]("ing8"), "ing8")
		assertSameValue(t, simpledi.Get[string]("complex_dish"), "complex_dish")
	})
}

func Test_Resolve_Complex_Graph_Mixed_Dependencies(t *testing.T) {
	defer simpledi.Close()
	order := make([]string, 0)

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "a",
			New: func() any {
				order = append(order, "a")
				return "a"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "b",
			New: func() any {
				order = append(order, "b")
				return "b"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "c",
			New: func() any {
				order = append(order, "c")
				return "c"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "d",
			New: func() any {
				order = append(order, "d")
				return "d"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "e",
			New: func() any {
				order = append(order, "e")
				return "e"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "r1",
			Deps: []string{"a", "b"},
			New: func() any {
				order = append(order, "r1")
				return "r1"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "r2",
			Deps: []string{"c", "r1"},
			New: func() any {
				order = append(order, "r2")
				return "r2"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "r3",
			Deps: []string{"r1", "d"},
			New: func() any {
				order = append(order, "r3")
				return "r3"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "r4",
			Deps: []string{"r2", "r3"},
			New: func() any {
				order = append(order, "r4")
				return "r4"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "r5",
			Deps: []string{"e"},
			New: func() any {
				order = append(order, "r5")
				return "r5"
			},
		})
		simpledi.Resolve()
	})

	assertOrder(t, order, []string{"a", "b", "c", "d", "e", "r1", "r5", "r2", "r3", "r4"})

	assertNoPanic(t, func() {
		assertSameValue(t, simpledi.Get[string]("a"), "a")
		assertSameValue(t, simpledi.Get[string]("b"), "b")
		assertSameValue(t, simpledi.Get[string]("c"), "c")
		assertSameValue(t, simpledi.Get[string]("d"), "d")
		assertSameValue(t, simpledi.Get[string]("e"), "e")
		assertSameValue(t, simpledi.Get[string]("r1"), "r1")
		assertSameValue(t, simpledi.Get[string]("r2"), "r2")
		assertSameValue(t, simpledi.Get[string]("r3"), "r3")
		assertSameValue(t, simpledi.Get[string]("r4"), "r4")
		assertSameValue(t, simpledi.Get[string]("r5"), "r5")
	})
}

func Test_Resolve_Circular_Dependency_Three_Recipes(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID:  "a",
			New: func() any { return "a" },
		})
		simpledi.Set(simpledi.Definition{
			ID:  "b",
			New: func() any { return "b" },
		})
		simpledi.Set(simpledi.Definition{
			ID:  "c",
			New: func() any { return "c" },
		})
		simpledi.Set(simpledi.Definition{
			ID:   "r1",
			Deps: []string{"r3", "a"},
			New:  func() any { return "r1" },
		})
		simpledi.Set(simpledi.Definition{
			ID:   "r2",
			Deps: []string{"r1", "b"},
			New:  func() any { return "r2" },
		})
		simpledi.Set(simpledi.Definition{
			ID:   "r3",
			Deps: []string{"r2", "c"},
			New:  func() any { return "r3" },
		})
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyCycle)
}

func Test_Resolve_Recipe_Requires_Itself(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID:  "flour",
			New: func() any { return "flour" },
		})
		simpledi.Set(simpledi.Definition{
			ID:   "infinite_recipe",
			Deps: []string{"infinite_recipe", "flour"},
			New:  func() any { return "infinite_recipe" },
		})
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyCycle)
}

func Test_Resolve_Recipe_Result_Already_In_Supplies(t *testing.T) {
	defer simpledi.Close()
	order := make([]string, 0)

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "flour",
			New: func() any {
				order = append(order, "flour")
				return "flour"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "meat",
			New: func() any {
				order = append(order, "meat")
				return "meat"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "bread",
			New: func() any {
				order = append(order, "bread")
				return "bread"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "sandwich",
			Deps: []string{"bread", "meat"},
			New: func() any {
				order = append(order, "sandwich")
				return "sandwich"
			},
		})
		simpledi.Resolve()
	})

	assertOrder(t, order, []string{"flour", "meat", "bread", "sandwich"})

	assertNoPanic(t, func() {
		assertSameValue(t, simpledi.Get[string]("flour"), "flour")
		assertSameValue(t, simpledi.Get[string]("meat"), "meat")
		assertSameValue(t, simpledi.Get[string]("bread"), "bread")
		assertSameValue(t, simpledi.Get[string]("sandwich"), "sandwich")
	})
}

func Test_Resolve_All_Recipe_Results_In_Supplies(t *testing.T) {
	defer simpledi.Close()
	order := make([]string, 0)

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "bread",
			New: func() any {
				order = append(order, "bread")
				return "bread"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "juice",
			New: func() any {
				order = append(order, "juice")
				return "juice"
			},
		})
		simpledi.Resolve()
	})

	assertOrder(t, order, []string{"bread", "juice"})

	assertNoPanic(t, func() {
		assertSameValue(t, simpledi.Get[string]("bread"), "bread")
		assertSameValue(t, simpledi.Get[string]("juice"), "juice")
	})
}

func Test_Resolve_Recipe_Missing_Last_Ingredient(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID:  "flour",
			New: func() any { return "flour" },
		})
		simpledi.Set(simpledi.Definition{
			ID:  "sugar",
			New: func() any { return "sugar" },
		})
		simpledi.Set(simpledi.Definition{
			ID:  "eggs",
			New: func() any { return "eggs" },
		})
		simpledi.Set(simpledi.Definition{
			ID:  "butter",
			New: func() any { return "butter" },
		})
		simpledi.Set(simpledi.Definition{
			ID:   "cake",
			Deps: []string{"flour", "sugar", "eggs", "butter", "vanilla"},
			New:  func() any { return "cake" },
		})
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyNotFound)
}

func Test_Resolve_Multiple_Chains_One_Complete_One_Broken(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID:  "a",
			New: func() any { return "a" },
		})
		simpledi.Set(simpledi.Definition{
			ID:  "b",
			New: func() any { return "b" },
		})
		simpledi.Set(simpledi.Definition{
			ID:  "c",
			New: func() any { return "c" },
		})
		simpledi.Set(simpledi.Definition{
			ID:   "r1",
			Deps: []string{"a"},
			New:  func() any { return "r1" },
		})
		simpledi.Set(simpledi.Definition{
			ID:   "r2",
			Deps: []string{"r1", "b"},
			New:  func() any { return "r2" },
		})
		simpledi.Set(simpledi.Definition{
			ID:   "r3",
			Deps: []string{"c"},
			New:  func() any { return "r3" },
		})
		simpledi.Set(simpledi.Definition{
			ID:   "r4",
			Deps: []string{"r3", "d"},
			New:  func() any { return "r4" },
		})
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyNotFound)
}

func Test_Resolve_Wide_Dependency_Tree(t *testing.T) {
	defer simpledi.Close()
	order := make([]string, 0)

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "a",
			New: func() any {
				order = append(order, "a")
				return "a"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "b",
			New: func() any {
				order = append(order, "b")
				return "b"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "c",
			New: func() any {
				order = append(order, "c")
				return "c"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "d",
			New: func() any {
				order = append(order, "d")
				return "d"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "r1",
			Deps: []string{"a"},
			New: func() any {
				order = append(order, "r1")
				return "r1"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "r2",
			Deps: []string{"b"},
			New: func() any {
				order = append(order, "r2")
				return "r2"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "r3",
			Deps: []string{"c"},
			New: func() any {
				order = append(order, "r3")
				return "r3"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "r4",
			Deps: []string{"d"},
			New: func() any {
				order = append(order, "r4")
				return "r4"
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "final",
			Deps: []string{"r1", "r2", "r3", "r4"},
			New: func() any {
				order = append(order, "final")
				return "final"
			},
		})
		simpledi.Resolve()
	})

	assertOrder(t, order, []string{"a", "b", "c", "d", "r1", "r2", "r3", "r4", "final"})

	assertNoPanic(t, func() {
		assertSameValue(t, simpledi.Get[string]("a"), "a")
		assertSameValue(t, simpledi.Get[string]("b"), "b")
		assertSameValue(t, simpledi.Get[string]("c"), "c")
		assertSameValue(t, simpledi.Get[string]("d"), "d")
		assertSameValue(t, simpledi.Get[string]("r1"), "r1")
		assertSameValue(t, simpledi.Get[string]("r2"), "r2")
		assertSameValue(t, simpledi.Get[string]("r3"), "r3")
		assertSameValue(t, simpledi.Get[string]("r4"), "r4")
		assertSameValue(t, simpledi.Get[string]("final"), "final")
	})
}

func Test_Resolve_Err_Container_Resolved(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Resolve()
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrContainerResolved)
}

func Test_Resolve_Err_ID_Duplicate(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "service_1",
			New: func() any {
				return &testServiceImpl1{}
			},
		})
		simpledi.Set(simpledi.Definition{
			ID: "service_1",
			New: func() any {
				return &testServiceImpl1{}
			},
		})
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrIDDuplicate)
}

func Test_Resolve_Err_Dependency_Not_Found(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID:   "service_1",
			Deps: []string{"service_2"},
			New: func() any {
				return &testServiceImpl1{}
			},
		})
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyNotFound)
}

func Test_Resolve_Err_Dependency_Cycle(t *testing.T) {
	defer simpledi.Close()

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID:   "service_1",
			Deps: []string{"service_2"},
			New: func() any {
				return &testServiceImpl1{}
			},
		})
		simpledi.Set(simpledi.Definition{
			ID:   "service_2",
			Deps: []string{"service_1"},
			New: func() any {
				return &testServiceImpl2{}
			},
		})
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyCycle)
}

func Test_Resolve_All_New_Functions_Invoked_Once(t *testing.T) {
	defer simpledi.Close()
	callCount := 0

	assertNoPanic(t, func() {
		simpledi.Set(simpledi.Definition{
			ID: "service_1",
			New: func() any {
				callCount++
				return &testServiceImpl1{}
			},
		})
		simpledi.Resolve()
	})

	assertSameValue(t, callCount, 1)
}

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

type testService1 interface{ doSomething1() }
type testServiceImpl1 struct{}

func (t *testServiceImpl1) doSomething1() {}

type testServiceImpl2 struct{ data string }

type testServiceImpl3 struct{ service1 *testServiceImpl1 }

func assertOrder[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	gotCount, wantCount := len(got), len(want)
	if len(got) != len(want) {
		t.Errorf("got: %d count, want: %d count", gotCount, wantCount)
		return
	}
	for i := 0; i < gotCount; i++ {
		if got[i] != want[i] {
			t.Errorf("[%d] got: %v, want: %v", i, got[i], want[i])
			t.Errorf("got: %#v, want: %#v", got, want)
			return
		}
	}
}

func assertSameValue[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func assertSamePointer[T comparable](t *testing.T, got, want *T) {
	t.Helper()
	if got != want {
		t.Errorf("got: %p, want: %p", got, want)
	}
}

func assertError(t *testing.T, fn func() error, wants ...error) {
	t.Helper()
	err := fn()
	for _, want := range wants {
		if !errors.Is(err, want) {
			t.Errorf("got: %v, want: %v", err, want)
		}
	}
}

func assertNoError(t *testing.T, fn func() error) {
	t.Helper()
	if err := fn(); err != nil {
		t.Errorf("got: %v, want: no error", err)
	}
}

func assertPanic(t *testing.T, fn func(), want error) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("got: no panic, want: panic(%v)", want)
			return
		}

		err, ok := r.(error)
		if !ok {
			t.Errorf("got: %T, want: error", r)
			return
		}

		if !errors.Is(err, want) {
			t.Errorf("got: %v, want: %v", err, want)
		}
	}()
	fn()
}

func assertNoPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("got: panic(%v), want: no panic", r)
		}
	}()
	fn()
}
