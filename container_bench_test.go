package simpledi_test

import (
	"fmt"
	"testing"

	"github.com/eerzho/simpledi"
)

func BenchmarkRegister(b *testing.B) {
	type TestStruct struct {
		value string
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := simpledi.NewContainer()

		b.StartTimer()
		c.MustRegister(simpledi.Option{
			Key: "test",
			Ctor: func() any {
				return &TestStruct{value: "test"}
			},
		})
		b.StopTimer()
	}
}

func BenchmarkResolve(b *testing.B) {
	type TestStruct struct {
		value string
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := simpledi.NewContainer()
		c.MustRegister(simpledi.Option{
			Key: "test",
			Ctor: func() any {
				return &TestStruct{value: "test"}
			},
		})

		b.StartTimer()
		c.MustResolve()
		b.StopTimer()
	}
}

func BenchmarkGet(b *testing.B) {
	type TestStruct struct {
		value string
	}

	c := simpledi.NewContainer()
	c.MustRegister(simpledi.Option{
		Key: "test",
		Ctor: func() any {
			return &TestStruct{value: "test"}
		},
	})
	c.MustResolve()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.MustGet("test")
	}
}

func BenchmarkReset(b *testing.B) {
	type TestStruct struct {
		value string
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := simpledi.NewContainer()
		c.MustRegister(simpledi.Option{
			Key: "test",
			Ctor: func() any {
				return &TestStruct{value: "test"}
			},
			Dtor: func() error {
				return nil
			},
		})
		c.MustResolve()

		b.StartTimer()
		c.MustReset()
		b.StopTimer()
	}
}

func BenchmarkFullPipeline(b *testing.B) {
	type TestStruct struct {
		value string
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := simpledi.NewContainer()
		c.MustRegister(simpledi.Option{
			Key: "test",
			Ctor: func() any {
				return &TestStruct{value: "test"}
			},
		})
		c.MustResolve()
		_ = c.MustGet("test")
		c.MustReset()
	}
}

func BenchmarkFullPipelineWith10Dependencies(b *testing.B) {
	type TestStruct struct {
		value string
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := simpledi.NewContainer()

		// Register 10 dependencies
		for j := 0; j < 10; j++ {
			key := fmt.Sprintf("test-%d", j)
			c.MustRegister(simpledi.Option{
				Key: key,
				Ctor: func() any {
					return &TestStruct{value: key}
				},
			})
		}

		c.MustResolve()

		// Get all dependencies
		for j := 0; j < 10; j++ {
			key := fmt.Sprintf("test-%d", j)
			_ = c.MustGet(key)
		}

		c.MustReset()
	}
}

func BenchmarkFullPipelineWith100Dependencies(b *testing.B) {
	type TestStruct struct {
		value string
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := simpledi.NewContainer()

		// Register 100 dependencies
		for j := 0; j < 100; j++ {
			key := fmt.Sprintf("test-%d", j)
			c.MustRegister(simpledi.Option{
				Key: key,
				Ctor: func() any {
					return &TestStruct{value: key}
				},
			})
		}

		c.MustResolve()

		// Get all dependencies
		for j := 0; j < 100; j++ {
			key := fmt.Sprintf("test-%d", j)
			_ = c.MustGet(key)
		}

		c.MustReset()
	}
}

func BenchmarkFullPipelineWithDependencyChain(b *testing.B) {
	type TestStruct struct {
		value string
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := simpledi.NewContainer()

		// Register dependencies with chain A -> B -> C -> D -> E
		c.MustRegister(simpledi.Option{
			Key: "a",
			Ctor: func() any {
				return &TestStruct{value: "a"}
			},
		})
		c.MustRegister(simpledi.Option{
			Key:  "b",
			Deps: []string{"a"},
			Ctor: func() any {
				return &TestStruct{value: "b"}
			},
		})
		c.MustRegister(simpledi.Option{
			Key:  "c",
			Deps: []string{"b"},
			Ctor: func() any {
				return &TestStruct{value: "c"}
			},
		})
		c.MustRegister(simpledi.Option{
			Key:  "d",
			Deps: []string{"c"},
			Ctor: func() any {
				return &TestStruct{value: "d"}
			},
		})
		c.MustRegister(simpledi.Option{
			Key:  "e",
			Deps: []string{"d"},
			Ctor: func() any {
				return &TestStruct{value: "e"}
			},
		})

		c.MustResolve()

		// Get all dependencies
		_ = c.MustGet("a")
		_ = c.MustGet("b")
		_ = c.MustGet("c")
		_ = c.MustGet("d")
		_ = c.MustGet("e")

		c.MustReset()
	}
}
