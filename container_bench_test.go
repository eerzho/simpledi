package simpledi_test

import (
	"fmt"
	"testing"

	"github.com/eerzho/simpledi"
)

func Benchmark(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			id := fmt.Sprintf("base_%d", j)
			simpledi.Set(simpledi.Definition{
				ID: id,
				New: func() any {
					return fmt.Sprintf("instance_%s", id)
				},
			})
		}

		for j := 0; j < 15; j++ {
			id := fmt.Sprintf("level1_%d", j)
			dep1 := fmt.Sprintf("base_%d", j%10)
			dep2 := fmt.Sprintf("base_%d", (j+1)%10)
			simpledi.Set(simpledi.Definition{
				ID:   id,
				Deps: []string{dep1, dep2},
				New: func() any {
					return fmt.Sprintf("instance_%s", id)
				},
			})
		}

		for j := 0; j < 20; j++ {
			id := fmt.Sprintf("level2_%d", j)
			dep1 := fmt.Sprintf("level1_%d", j%15)
			dep2 := fmt.Sprintf("level1_%d", (j+1)%15)
			dep3 := fmt.Sprintf("level1_%d", (j+2)%15)
			simpledi.Set(simpledi.Definition{
				ID:   id,
				Deps: []string{dep1, dep2, dep3},
				New: func() any {
					return fmt.Sprintf("instance_%s", id)
				},
			})
		}

		for j := 0; j < 15; j++ {
			id := fmt.Sprintf("level3_%d", j)
			dep1 := fmt.Sprintf("level2_%d", j%20)
			dep2 := fmt.Sprintf("level2_%d", (j+1)%20)
			dep3 := fmt.Sprintf("level2_%d", (j+2)%20)
			dep4 := fmt.Sprintf("level2_%d", (j+3)%20)
			simpledi.Set(simpledi.Definition{
				ID:   id,
				Deps: []string{dep1, dep2, dep3, dep4},
				New: func() any {
					return fmt.Sprintf("instance_%s", id)
				},
			})
		}

		for j := 0; j < 10; j++ {
			id := fmt.Sprintf("level4_%d", j)
			dep1 := fmt.Sprintf("level3_%d", j%15)
			dep2 := fmt.Sprintf("level3_%d", (j+1)%15)
			dep3 := fmt.Sprintf("level3_%d", (j+2)%15)
			simpledi.Set(simpledi.Definition{
				ID:   id,
				Deps: []string{dep1, dep2, dep3},
				New: func() any {
					return fmt.Sprintf("instance_%s", id)
				},
			})
		}

		for j := 0; j < 5; j++ {
			id := fmt.Sprintf("final_%d", j)
			dep1 := fmt.Sprintf("level4_%d", j%10)
			dep2 := fmt.Sprintf("level4_%d", (j+1)%10)
			simpledi.Set(simpledi.Definition{
				ID:   id,
				Deps: []string{dep1, dep2},
				New: func() any {
					return fmt.Sprintf("instance_%s", id)
				},
			})
		}

		simpledi.Resolve()
		simpledi.Close()
	}
}
