package simpledi_test

import (
	"fmt"

	"github.com/eerzho/simpledi"
)

func Example() {
	type repo struct {
		dsn string
	}
	type service struct {
		repo *repo
	}
	// create container
	c := simpledi.NewContainer()
	// register dependencies
	c.MustRegister(simpledi.Option{
		Key:          "repo",
		Dependencies: nil,
		Constructor: func() any {
			fmt.Println("creating repo")
			return &repo{dsn: "example"}
		},
	})
	c.MustRegister(simpledi.Option{
		Key:          "service",
		Dependencies: []string{"repo"},
		Constructor: func() any {
			fmt.Println("creating service")
			return &service{repo: c.MustGet("repo").(*repo)}
		},
	})
	// resolve all dependencies
	c.MustResolve()
	// get resolved dependency
	fmt.Println(c.MustGet("service").(*service).repo.dsn)
	// Output:
	// creating repo
	// creating service
	// example
}

func Example_defaultContainer() {
	type repo struct {
		dsn string
	}
	type service struct {
		repo *repo
	}
	// register dependencies
	simpledi.MustRegister(simpledi.Option{
		Key: "repo",
		Constructor: func() any {
			fmt.Println("creating repo")
			return &repo{dsn: "example"}
		},
	})
	simpledi.MustRegister(simpledi.Option{
		Key:          "service",
		Dependencies: []string{"repo"},
		Constructor: func() any {
			fmt.Println("creating service")
			return &service{repo: simpledi.MustGetAs[*repo]("repo")}
		},
	})
	// resolve all dependencies
	simpledi.MustResolve()
	// get resolved dependency
	fmt.Println(simpledi.MustGetAs[*service]("service").repo.dsn)
	// Output:
	// creating repo
	// creating service
	// example
}

func Example_registerInAnyOrder() {
	type repo struct {
		dsn string
	}
	type service struct {
		repo *repo
	}
	// create container
	c := simpledi.NewContainer()
	// register dependencies in any order
	c.MustRegister(simpledi.Option{
		Key:          "service",
		Dependencies: []string{"repo"},
		Constructor: func() any {
			fmt.Println("creating service")
			return &service{repo: c.MustGet("repo").(*repo)}
		},
	})
	c.MustRegister(simpledi.Option{
		Key: "repo",
		Constructor: func() any {
			fmt.Println("creating repo")
			return &repo{dsn: "example"}
		},
	})
	// resolve all dependencies
	c.MustResolve()
	// get resolved dependency
	fmt.Println(c.MustGet("service").(*service).repo.dsn)
	// Output:
	// creating repo
	// creating service
	// example
}

func Example_overrideDependency() {
	type repo struct {
		dsn string
	}
	type service struct {
		repo *repo
	}
	// create container
	c := simpledi.NewContainer()
	// register dependencies
	c.MustRegister(simpledi.Option{
		Key: "repo",
		Constructor: func() any {
			fmt.Println("creating repo")
			return &repo{dsn: "example"}
		},
	})
	c.MustRegister(simpledi.Option{
		Key:          "service",
		Dependencies: []string{"repo"},
		Constructor: func() any {
			fmt.Println("creating service")
			return &service{repo: c.MustGet("repo").(*repo)}
		},
	})
	// resolve all dependencies
	c.MustResolve()
	// get resolved dependency
	fmt.Println(c.MustGet("service").(*service).repo.dsn)
	// override the "repo"
	c.MustRegister(simpledi.Option{
		Key: "repo",
		Constructor: func() any {
			fmt.Println("creating mock")
			return &repo{dsn: "example-2"}
		},
	})
	// resolve all dependencies again
	c.MustResolve()
	// get resolved dependency
	fmt.Println(c.MustGet("service").(*service).repo.dsn)
	// Output:
	// creating repo
	// creating service
	// example
	// creating mock
	// creating service
	// example-2
}
