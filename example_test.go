package simpledi_test

import (
	"fmt"

	"github.com/eerzho/simpledi"
)

func Example() {
	type repo struct {
		DSN string
	}
	type service struct {
		repo *repo
	}

	// create container
	c := simpledi.NewContainer()

	// register dependencies
	c.Register("repo", nil, func() any {
		fmt.Println("creating repo")
		return &repo{DSN: "example"}
	})
	c.Register("service", []string{"repo"}, func() any {
		fmt.Println("creating service")
		return &service{repo: c.Get("repo").(*repo)}
	})

	// resolve all dependencies
	c.Resolve()

	// get resolved dependency
	s := c.Get("service").(*service)
	fmt.Println(s.repo.DSN)

	// Output:
	// creating repo
	// creating service
	// example
}

func Example_registerInAnyOrder() {
	type repo struct {
		DSN string
	}
	type service struct {
		repo *repo
	}

	// create container
	c := simpledi.NewContainer()

	// register dependencies in any order
	c.Register("service", []string{"repo"}, func() any {
		fmt.Println("creating service")
		return &service{repo: c.Get("repo").(*repo)}
	})
	c.Register("repo", nil, func() any {
		fmt.Println("creating repo")
		return &repo{DSN: "example"}
	})

	// resolve all dependencies
	c.Resolve()

	// get resolved dependency
	s := c.Get("service").(*service)
	fmt.Println(s.repo.DSN)

	// Output:
	// creating repo
	// creating service
	// example
}

func Example_overrideDependency() {
	type repo struct {
		DSN string
	}
	type service struct {
		repo *repo
	}

	// create container
	c := simpledi.NewContainer()

	// register dependencies
	c.Register("repo", nil, func() any {
		fmt.Println("creating repo")
		return &repo{DSN: "example"}
	})
	c.Register("service", []string{"repo"}, func() any {
		fmt.Println("creating service")
		return &service{repo: c.Get("repo").(*repo)}
	})

	// resolve all dependencies
	c.Resolve()

	// get resolved dependency
	s := c.Get("service").(*service)
	fmt.Println(s.repo.DSN)

	// override the "repo"
	c.Register("repo", nil, func() any {
		fmt.Println("creating mock")
		return &repo{DSN: "mock-example"}
	})

	// resolve all dependencies again
	c.Resolve()

	// get resolved dependency
	s = c.Get("service").(*service)
	fmt.Println(s.repo.DSN)

	// Output:
	// creating repo
	// creating service
	// example
	// creating mock
	// creating service
	// mock-example
}

func Example_defaultContainer() {
	type repo struct {
		DSN string
	}
	type service struct {
		repo *repo
	}

	// register dependencies
	simpledi.Register("repo", nil, func() any {
		fmt.Println("creating repo")
		return &repo{DSN: "example"}
	})
	simpledi.Register("service", []string{"repo"}, func() any {
		fmt.Println("creating service")
		return &service{repo: simpledi.Get("repo").(*repo)}
	})

	// resolve all dependencies
	simpledi.Resolve()

	// get resolved dependency
	s := simpledi.Get("service").(*service)
	fmt.Println(s.repo.DSN)

	// Output:
	// creating repo
	// creating service
	// example
}
