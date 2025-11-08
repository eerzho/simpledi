package simpledi_test

import (
	"fmt"
	"testing"

	"github.com/eerzho/simpledi"
)

func Test_Basic_Single_Recipe_Available(t *testing.T) {
	defer simpledi.Close()
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

	assertNoPanic(t, func() {
		simpledi.Resolve()

		_ = simpledi.Get[string]("yeast")
		_ = simpledi.Get[string]("flour")
		_ = simpledi.Get[string]("bread")
	})
}

func Test_Single_Recipe_Missing_Ingredient(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "yeast",
		New: func() any {
			return "yeast"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "bread",
		Deps: []string{"yeast", "flour"},
		New: func() any {
			return "bread"
		},
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyNotFound)
}

func Test_Chain_Recipe_Two_Levels(t *testing.T) {
	defer simpledi.Close()
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
		ID: "meat",
		New: func() any {
			return "meat"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "bread",
		Deps: []string{"yeast", "flour"},
		New: func() any {
			return "bread"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "sandwich",
		Deps: []string{"bread", "meat"},
		New: func() any {
			return "sandwich"
		},
	})

	assertNoPanic(t, func() {
		simpledi.Resolve()

		_ = simpledi.Get[string]("yeast")
		_ = simpledi.Get[string]("flour")
		_ = simpledi.Get[string]("meat")
		_ = simpledi.Get[string]("bread")
		_ = simpledi.Get[string]("sandwich")
	})
}

func Test_Chain_Recipe_Broken(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "yeast",
		New: func() any {
			return "yeast"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "meat",
		New: func() any {
			return "meat"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "bread",
		Deps: []string{"yeast", "flour"},
		New: func() any {
			return "bread"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "sandwich",
		Deps: []string{"bread", "meat"},
		New: func() any {
			return "sandwich"
		},
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyNotFound)
}

func Test_Multiple_Independent_Recipes(t *testing.T) {
	defer simpledi.Close()
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
		ID: "orange",
		New: func() any {
			return "orange"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "water",
		New: func() any {
			return "water"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "tomato",
		New: func() any {
			return "tomato"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "cucumber",
		New: func() any {
			return "cucumber"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "bread",
		Deps: []string{"yeast", "flour"},
		New: func() any {
			return "bread"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "juice",
		Deps: []string{"orange", "water"},
		New: func() any {
			return "juice"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "salad",
		Deps: []string{"tomato", "cucumber"},
		New: func() any {
			return "salad"
		},
	})

	assertNoPanic(t, func() {
		simpledi.Resolve()

		_ = simpledi.Get[string]("yeast")
		_ = simpledi.Get[string]("flour")
		_ = simpledi.Get[string]("orange")
		_ = simpledi.Get[string]("water")
		_ = simpledi.Get[string]("tomato")
		_ = simpledi.Get[string]("cucumber")
		_ = simpledi.Get[string]("bread")
		_ = simpledi.Get[string]("juice")
		_ = simpledi.Get[string]("salad")
	})
}

func Test_Multiple_Recipes_Partial_Available(t *testing.T) {
	defer simpledi.Close()
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
		ID: "tomato",
		New: func() any {
			return "tomato"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "cucumber",
		New: func() any {
			return "cucumber"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "bread",
		Deps: []string{"yeast", "flour"},
		New: func() any {
			return "bread"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "juice",
		Deps: []string{"orange", "water"},
		New: func() any {
			return "juice"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "salad",
		Deps: []string{"tomato", "cucumber"},
		New: func() any {
			return "salad"
		},
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyNotFound)
}

func Test_Empty_Supplies(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID:   "bread",
		Deps: []string{"yeast", "flour"},
		New: func() any {
			return "bread"
		},
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyNotFound)
}

func Test_Empty_Recipes(t *testing.T) {
	defer simpledi.Close()
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

	assertNoPanic(t, func() {
		simpledi.Resolve()

		_ = simpledi.Get[string]("yeast")
		_ = simpledi.Get[string]("flour")
	})
}

func Test_Circular_Dependency_Simple(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "salt",
		New: func() any {
			return "salt"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "pepper",
		New: func() any {
			return "pepper"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "dish_a",
		Deps: []string{"dish_b", "salt"},
		New: func() any {
			return "dish_a"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "dish_b",
		Deps: []string{"dish_a", "pepper"},
		New: func() any {
			return "dish_b"
		},
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyCycle)
}

func Test_Complex_Chain_Three_Levels(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "flour",
		New: func() any {
			return "flour"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "sugar",
		New: func() any {
			return "sugar"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "milk",
		New: func() any {
			return "milk"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "dish_a",
		Deps: []string{"flour"},
		New: func() any {
			return "dish_a"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "dish_b",
		Deps: []string{"dish_a", "sugar"},
		New: func() any {
			return "dish_b"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "dish_c",
		Deps: []string{"dish_b", "milk"},
		New: func() any {
			return "dish_c"
		},
	})

	assertNoPanic(t, func() {
		simpledi.Resolve()

		_ = simpledi.Get[string]("dish_c")
		_ = simpledi.Get[string]("sugar")
		_ = simpledi.Get[string]("milk")
		_ = simpledi.Get[string]("dish_a")
		_ = simpledi.Get[string]("dish_b")
		_ = simpledi.Get[string]("dish_c")
	})
}

func Test_Diamond_Dependency(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "flour",
		New: func() any {
			return "flour"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "sugar",
		New: func() any {
			return "sugar"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "salt",
		New: func() any {
			return "salt"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "dish_a",
		Deps: []string{"flour"},
		New: func() any {
			return "dish_a"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "dish_b",
		Deps: []string{"dish_a", "sugar"},
		New: func() any {
			return "dish_b"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "dish_c",
		Deps: []string{"dish_a", "salt"},
		New: func() any {
			return "dish_c"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "dish_d",
		Deps: []string{"dish_b", "dish_c"},
		New: func() any {
			return "dish_d"
		},
	})

	assertNoPanic(t, func() {
		simpledi.Resolve()

		_ = simpledi.Get[string]("flour")
		_ = simpledi.Get[string]("sugar")
		_ = simpledi.Get[string]("salt")
		_ = simpledi.Get[string]("dish_a")
		_ = simpledi.Get[string]("dish_b")
		_ = simpledi.Get[string]("dish_c")
		_ = simpledi.Get[string]("dish_d")
	})
}

func Test_Multiple_Uses_Of_Same_Ingredient(t *testing.T) {
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
	simpledi.Set(simpledi.Definition{
		ID: "tomato",
		New: func() any {
			return "tomato"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "egg",
		New: func() any {
			return "egg"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "bread",
		Deps: []string{"flour", "water"},
		New: func() any {
			return "bread"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "pizza",
		Deps: []string{"flour", "tomato"},
		New: func() any {
			return "pizza"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "pasta",
		Deps: []string{"flour", "egg"},
		New: func() any {
			return "pasta"
		},
	})

	assertNoPanic(t, func() {
		simpledi.Resolve()

		_ = simpledi.Get[string]("flour")
		_ = simpledi.Get[string]("water")
		_ = simpledi.Get[string]("tomato")
		_ = simpledi.Get[string]("egg")
		_ = simpledi.Get[string]("bread")
		_ = simpledi.Get[string]("pizza")
		_ = simpledi.Get[string]("pasta")
	})
}

func Test_Recipe_Requires_Another_Recipe_Multiple_Times(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "flour",
		New: func() any {
			return "flour"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "meat",
		New: func() any {
			return "meat"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "bread",
		Deps: []string{"flour"},
		New: func() any {
			return "bread"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "double_sandwich",
		Deps: []string{"bread", "bread", "meat"},
		New: func() any {
			return "double_sandwich"
		},
	})

	assertNoPanic(t, func() {
		simpledi.Resolve()

		_ = simpledi.Get[string]("flour")
		_ = simpledi.Get[string]("meat")
		_ = simpledi.Get[string]("bread")
		_ = simpledi.Get[string]("double_sandwich")
	})
}

func Test_Extra_Unused_Supplies(t *testing.T) {
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
	simpledi.Set(simpledi.Definition{
		ID: "sugar",
		New: func() any {
			return "sugar"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "salt",
		New: func() any {
			return "salt"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "pepper",
		New: func() any {
			return "pepper"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "cheese",
		New: func() any {
			return "cheese"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "bread",
		Deps: []string{"flour", "water"},
		New: func() any {
			return "bread"
		},
	})

	assertNoPanic(t, func() {
		simpledi.Resolve()

		_ = simpledi.Get[string]("flour")
		_ = simpledi.Get[string]("water")
		_ = simpledi.Get[string]("sugar")
		_ = simpledi.Get[string]("salt")
		_ = simpledi.Get[string]("pepper")
		_ = simpledi.Get[string]("cheese")
		_ = simpledi.Get[string]("bread")
	})
}

func Test_All_Recipes_Missing_One_Common_Ingredient(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "flour",
		New: func() any {
			return "flour"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "sugar",
		New: func() any {
			return "sugar"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "water",
		New: func() any {
			return "water"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "dish_a",
		Deps: []string{"salt", "flour"},
		New: func() any {
			return "dish_a"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "dish_b",
		Deps: []string{"salt", "sugar"},
		New: func() any {
			return "dish_b"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "dish_c",
		Deps: []string{"salt", "water"},
		New: func() any {
			return "dish_c"
		},
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyNotFound)
}

func Test_Long_Chain_Five_Levels(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "base",
		New: func() any {
			return "base"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "ing1",
		New: func() any {
			return "ing1"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "ing2",
		New: func() any {
			return "ing2"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "ing3",
		New: func() any {
			return "ing3"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "ing4",
		New: func() any {
			return "ing4"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "level1",
		Deps: []string{"base"},
		New: func() any {
			return "level1"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "level2",
		Deps: []string{"level1", "ing1"},
		New: func() any {
			return "level2"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "level3",
		Deps: []string{"level2", "ing2"},
		New: func() any {
			return "level3"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "level4",
		Deps: []string{"level3", "ing3"},
		New: func() any {
			return "level4"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "level5",
		Deps: []string{"level4", "ing4"},
		New: func() any {
			return "level5"
		},
	})

	assertNoPanic(t, func() {
		simpledi.Resolve()

		_ = simpledi.Get[string]("base")
		_ = simpledi.Get[string]("ing1")
		_ = simpledi.Get[string]("ing2")
		_ = simpledi.Get[string]("ing3")
		_ = simpledi.Get[string]("ing4")
		_ = simpledi.Get[string]("level1")
		_ = simpledi.Get[string]("level2")
		_ = simpledi.Get[string]("level3")
		_ = simpledi.Get[string]("level4")
		_ = simpledi.Get[string]("level5")
	})
}

func Test_Parallel_Chains(t *testing.T) {
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
	simpledi.Set(simpledi.Definition{
		ID: "sugar",
		New: func() any {
			return "sugar"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "milk",
		New: func() any {
			return "milk"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "chain1_a",
		Deps: []string{"flour"},
		New: func() any {
			return "chain1_a"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "chain1_b",
		Deps: []string{"chain1_a", "water"},
		New: func() any {
			return "chain1_b"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "chain2_a",
		Deps: []string{"sugar"},
		New: func() any {
			return "chain2_a"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "chain2_b",
		Deps: []string{"chain2_a", "milk"},
		New: func() any {
			return "chain2_b"
		},
	})

	assertNoPanic(t, func() {
		simpledi.Resolve()

		_ = simpledi.Get[string]("flour")
		_ = simpledi.Get[string]("water")
		_ = simpledi.Get[string]("sugar")
		_ = simpledi.Get[string]("milk")
		_ = simpledi.Get[string]("chain1_a")
		_ = simpledi.Get[string]("chain1_b")
		_ = simpledi.Get[string]("chain2_a")
		_ = simpledi.Get[string]("chain2_b")
	})
}

func Test_Recipe_With_Single_Ingredient(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "orange",
		New: func() any {
			return "orange"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "juice",
		Deps: []string{"orange"},
		New: func() any {
			return "juice"
		},
	})

	assertNoPanic(t, func() {
		simpledi.Resolve()

		_ = simpledi.Get[string]("orange")
		_ = simpledi.Get[string]("juice")
	})
}

func Test_Recipe_With_Many_Ingredients(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "ing1",
		New: func() any {
			return "ing1"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "ing2",
		New: func() any {
			return "ing2"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "ing3",
		New: func() any {
			return "ing3"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "ing4",
		New: func() any {
			return "ing4"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "ing5",
		New: func() any {
			return "ing5"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "ing6",
		New: func() any {
			return "ing6"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "ing7",
		New: func() any {
			return "ing7"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "ing8",
		New: func() any {
			return "ing8"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "complex_dish",
		Deps: []string{"ing1", "ing2", "ing3", "ing4", "ing5", "ing6", "ing7", "ing8"},
		New: func() any {
			return "complex_dish"
		},
	})

	assertNoPanic(t, func() {
		simpledi.Resolve()

		_ = simpledi.Get[string]("ing1")
		_ = simpledi.Get[string]("ing2")
		_ = simpledi.Get[string]("ing3")
		_ = simpledi.Get[string]("ing4")
		_ = simpledi.Get[string]("ing5")
		_ = simpledi.Get[string]("ing6")
		_ = simpledi.Get[string]("ing7")
		_ = simpledi.Get[string]("ing8")
		_ = simpledi.Get[string]("complex_dish")
	})
}

func Test_Complex_Graph_Mixed_Dependencies(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "a",
		New: func() any {
			return "a"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "b",
		New: func() any {
			return "b"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "c",
		New: func() any {
			return "c"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "d",
		New: func() any {
			return "d"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "e",
		New: func() any {
			return "e"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "r1",
		Deps: []string{"a", "b"},
		New: func() any {
			return "r1"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "r2",
		Deps: []string{"c", "r1"},
		New: func() any {
			return "r2"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "r3",
		Deps: []string{"r1", "d"},
		New: func() any {
			return "r3"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "r4",
		Deps: []string{"r2", "r3"},
		New: func() any {
			return "r4"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "r5",
		Deps: []string{"e"},
		New: func() any {
			return "r5"
		},
	})

	assertNoPanic(t, func() {
		simpledi.Resolve()

		_ = simpledi.Get[string]("a")
		_ = simpledi.Get[string]("b")
		_ = simpledi.Get[string]("c")
		_ = simpledi.Get[string]("d")
		_ = simpledi.Get[string]("e")
		_ = simpledi.Get[string]("r1")
		_ = simpledi.Get[string]("r2")
		_ = simpledi.Get[string]("r3")
		_ = simpledi.Get[string]("r4")
		_ = simpledi.Get[string]("r5")
	})
}

func Test_Circular_Dependency_Complex(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "a",
		New: func() any {
			return "a"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "b",
		New: func() any {
			return "b"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "c",
		New: func() any {
			return "c"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "r1",
		Deps: []string{"r3", "a"},
		New: func() any {
			return "r1"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "r2",
		Deps: []string{"r1", "b"},
		New: func() any {
			return "r2"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "r3",
		Deps: []string{"r2", "c"},
		New: func() any {
			return "r3"
		},
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyCycle)
}

func Test_Self_Dependency(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "flour",
		New: func() any {
			return "flour"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "infinite_recipe",
		Deps: []string{"infinite_recipe", "flour"},
		New: func() any {
			return "infinite_recipe"
		},
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrDependencyCycle)
}

func Test_Recipe_Already_In_Supplies(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "flour",
		New: func() any {
			return "flour"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "meat",
		New: func() any {
			return "meat"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "bread",
		New: func() any {
			return "bread"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "bread_recipe",
		Deps: []string{"flour"},
		New: func() any {
			return "bread"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "sandwich",
		Deps: []string{"bread", "meat"},
		New: func() any {
			return "sandwich"
		},
	})

	assertNoPanic(t, func() {
		simpledi.Resolve()

		_ = simpledi.Get[string]("flour")
		_ = simpledi.Get[string]("meat")
		_ = simpledi.Get[string]("bread")
		_ = simpledi.Get[string]("sandwich")
	})
}

func Test_All_Supplies_No_Recipes_Needed(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "bread",
		New: func() any {
			return "bread"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID: "juice",
		New: func() any {
			return "juice"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "bread",
		Deps: []string{"flour"},
		New: func() any {
			return "bread"
		},
	})
	simpledi.Set(simpledi.Definition{
		ID:   "juice",
		Deps: []string{"orange"},
		New: func() any {
			return "juice"
		},
	})

	assertPanic(t, func() {
		simpledi.Resolve()
	}, simpledi.ErrIDDuplicate)
}

func Test_Resolve_Large_Number_Of_Definitions(t *testing.T) {
	defer simpledi.Close()

	for i := 0; i < 1000; i++ {
		id := fmt.Sprintf("ingredient_%d", i)
		simpledi.Set(simpledi.Definition{
			ID: id,
			New: func() any {
				return id
			},
		})
	}
	simpledi.Resolve()

	assertNoPanic(t, func() {
		_ = simpledi.Get[string]("ingredient_0")
		_ = simpledi.Get[string]("ingredient_500")
		_ = simpledi.Get[string]("ingredient_999")
	})
}

func Test_Resolve_Deep_Dependency_Chain(t *testing.T) {
	defer simpledi.Close()

	simpledi.Set(simpledi.Definition{
		ID: "level0",
		New: func() any {
			return "level0"
		},
	})
	for i := 1; i < 50; i++ {
		prevID := fmt.Sprintf("level%d", i-1)
		currID := fmt.Sprintf("level%d", i)
		prevIDCopy := prevID
		currIDCopy := currID
		simpledi.Set(simpledi.Definition{
			ID:   currIDCopy,
			Deps: []string{prevIDCopy},
			New: func() any {
				return currIDCopy
			},
		})
	}
	simpledi.Resolve()

	assertNoPanic(t, func() {
		_ = simpledi.Get[string]("level0")
		_ = simpledi.Get[string]("level25")
		_ = simpledi.Get[string]("level49")
	})
}

func Test_Resolve_Wide_Dependency_Graph(t *testing.T) {
	defer simpledi.Close()

	baseCount := 100
	for i := 0; i < baseCount; i++ {
		id := fmt.Sprintf("base_%d", i)
		simpledi.Set(simpledi.Definition{
			ID: id,
			New: func() any {
				return id
			},
		})
	}
	deps := make([]string, baseCount)
	for i := 0; i < baseCount; i++ {
		deps[i] = fmt.Sprintf("base_%d", i)
	}
	simpledi.Set(simpledi.Definition{
		ID:   "mega_dish",
		Deps: deps,
		New: func() any {
			return "mega_dish"
		},
	})
	simpledi.Resolve()

	assertNoPanic(t, func() {
		_ = simpledi.Get[string]("base_0")
		_ = simpledi.Get[string]("base_50")
		_ = simpledi.Get[string]("base_99")
		_ = simpledi.Get[string]("mega_dish")
	})
}
