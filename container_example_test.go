package simpledi_test

import (
	"fmt"

	"github.com/eerzho/simpledi"
)

func ExampleContainer_basic() {
	type (
		db struct {
			dsn string
		}
		repo struct {
			db *db
		}
		service struct {
			repo *repo
		}
	)

	c := simpledi.NewContainer()
	dsn := "example://example:example"

	c.Register("db", nil, func() any {
		fmt.Println("creating: db")
		return &db{dsn: dsn}
	})
	c.Register("repo", []string{"db"}, func() any {
		fmt.Println("creating: repo")
		db := c.Get("db").(*db)
		return &repo{db: db}
	})
	c.Register("service", []string{"repo"}, func() any {
		fmt.Println("creating: service")
		repo := c.Get("repo").(*repo)
		return &service{repo: repo}
	})

	c.Resolve()

	fmt.Println(c.Get("service").(*service).repo.db.dsn)

	// Output:
	// creating: db
	// creating: repo
	// creating: service
	// example://example:example
}

func ExampleContainer_complex() {
	type (
		db struct {
			dsn string
		}
		repo1 struct {
			db *db
		}
		repo2 struct {
			db *db
		}
		repo3 struct {
			db *db
		}
		service1 struct {
			repo1 *repo1
			repo2 *repo2
		}
		service2 struct {
			service1 *service1
			repo3    *repo3
		}
	)

	c := simpledi.NewContainer()
	dsn := "example://example:example"

	c.Register("db", nil, func() any {
		fmt.Println("creating: db")
		return &db{dsn: dsn}
	})
	c.Register("repo1", []string{"db"}, func() any {
		db := c.Get("db").(*db)
		return &repo1{db: db}
	})
	c.Register("repo2", []string{"db"}, func() any {
		db := c.Get("db").(*db)
		return &repo2{db: db}
	})
	c.Register("repo3", []string{"db"}, func() any {
		db := c.Get("db").(*db)
		return &repo3{db: db}
	})
	c.Register("service1", []string{"repo1", "repo2"}, func() any {
		fmt.Println("creating: service1")
		repo1 := c.Get("repo1").(*repo1)
		repo2 := c.Get("repo2").(*repo2)
		return &service1{repo1: repo1, repo2: repo2}
	})
	c.Register("service2", []string{"service1", "repo3"}, func() any {
		fmt.Println("creating: service2")
		service1 := c.Get("service1").(*service1)
		repo3 := c.Get("repo3").(*repo3)
		return &service2{service1: service1, repo3: repo3}
	})

	c.Resolve()

	fmt.Println(c.Get("service2").(*service2).service1.repo1.db.dsn)
	fmt.Println(c.Get("service2").(*service2).service1.repo2.db.dsn)
	fmt.Println(c.Get("service2").(*service2).repo3.db.dsn)

	// Output:
	// creating: db
	// creating: service1
	// creating: service2
	// example://example:example
	// example://example:example
	// example://example:example
}
