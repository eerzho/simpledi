package simpledi_test

import (
	"fmt"
	"testing"

	"github.com/eerzho/simpledi"
)

func Benchmark(b *testing.B) {
	const count = 1000
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := simpledi.NewContainer()
		for j := 0; j < count; j++ {
			key := fmt.Sprintf("key-%d", j)
			c.Register(
				key,
				nil,
				func() any {
					return key
				},
			)
		}
		c.Resolve()
		for j := 0; j < count; j++ {
			key := fmt.Sprintf("key-%d", j)
			v := c.Get(key)
			_ = v
		}
	}
}

func BenchmarkWithDeps(b *testing.B) {
	const count = 1000
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := simpledi.NewContainer()
		prevKeys := []string{}
		for j := 0; j < count; j++ {
			key := fmt.Sprintf("key-%d", j)
			c.Register(
				key,
				prevKeys,
				func() any {
					return key
				},
			)
			prevKeys = append(prevKeys, key)
		}
		c.Resolve()
		for j := 0; j < count; j++ {
			key := fmt.Sprintf("key-%d", j)
			v := c.Get(key)
			_ = v
		}
	}
}

func BenchmarkWithRealisticDeps(b *testing.B) {
	const count = 1000
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := simpledi.NewContainer()
		for j := 0; j < count; j++ {
			key := fmt.Sprintf("key-%d", j)
			deps := []string{}
			if j < 100 {
				deps = nil
			} else if j < 500 {
				deps = []string{
					fmt.Sprintf("key-%d", j%100),
					fmt.Sprintf("key-%d", (j+1)%100),
				}
			} else {
				deps = []string{
					fmt.Sprintf("key-%d", 100+(j-500)%400),
					fmt.Sprintf("key-%d", 100+((j-500)+1)%400),
				}
			}
			c.Register(
				key,
				deps,
				func() any {
					return key
				},
			)
		}
		c.Resolve()
		for j := 0; j < count; j++ {
			key := fmt.Sprintf("key-%d", j)
			v := c.Get(key)
			_ = v
		}
	}
}
