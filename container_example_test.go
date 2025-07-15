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
	c.Register("repo", nil, func() any {
		fmt.Println("creating repo")
		return &repo{dsn: "example"}
	})
	c.Register("service", []string{"repo"}, func() any {
		fmt.Println("creating service")
		return &service{repo: c.MustGet("repo").(*repo)}
	})
	// resolve all dependencies
	c.Resolve()
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
	simpledi.Register("repo", nil, func() any {
		fmt.Println("creating repo")
		return &repo{dsn: "example"}
	})
	simpledi.Register("service", []string{"repo"}, func() any {
		fmt.Println("creating service")
		return &service{repo: simpledi.MustGetAs[*repo]("repo")}
	})
	// resolve all dependencies
	simpledi.Resolve()
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
	c.Register("service", []string{"repo"}, func() any {
		fmt.Println("creating service")
		return &service{repo: c.MustGet("repo").(*repo)}
	})
	c.Register("repo", nil, func() any {
		fmt.Println("creating repo")
		return &repo{dsn: "example"}
	})
	// resolve all dependencies
	c.Resolve()
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
	c.Register("repo", nil, func() any {
		fmt.Println("creating repo")
		return &repo{dsn: "example"}
	})
	c.Register("service", []string{"repo"}, func() any {
		fmt.Println("creating service")
		return &service{repo: c.MustGet("repo").(*repo)}
	})
	// resolve all dependencies
	c.Resolve()
	// get resolved dependency
	fmt.Println(c.MustGet("service").(*service).repo.dsn)
	// override the "repo"
	c.Register("repo", nil, func() any {
		fmt.Println("creating mock")
		return &repo{dsn: "example-2"}
	})
	// resolve all dependencies again
	c.Resolve()
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
