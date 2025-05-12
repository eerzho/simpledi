package simpledi_test

import (
	"fmt"

	"github.com/eerzho/simpledi"
)

func ExampleContainer() {
	c := simpledi.NewContainer()

	c.Register("yeast", nil, func() any {
		fmt.Println("yeast created")
		return "yeast"
	})
	c.Register("flour", nil, func() any {
		fmt.Println("flour created")
		return "flour"
	})
	c.Register("meat", nil, func() any {
		fmt.Println("meat created")
		return "meat"
	})
	c.Register("bread", []string{"yeast", "flour"}, func() any {
		fmt.Println("bread created using:", c.Get("yeast"), "and", c.Get("flour"))
		return "bread"
	})
	c.Register("sandwich", []string{"bread", "meat"}, func() any {
		fmt.Println("sandwich created using:", c.Get("bread"), "and", c.Get("meat"))
		return "sandwich"
	})
	c.Register("burger", []string{"sandwich", "meat", "bread"}, func() any {
		fmt.Println("burger created using:", c.Get("sandwich"), c.Get("meat"), "and", c.Get("bread"))
		return "burger"
	})

	if err := c.Resolve(); err != nil {
		panic(err)
	}

	fmt.Println("final product:", c.Get("burger"))

	// Output:
	// yeast created
	// flour created
	// meat created
	// bread created using: yeast and flour
	// sandwich created using: bread and meat
	// burger created using: sandwich meat and bread
	// final product: burger
}

func ExampleContainer_withServices() {
	type db struct {
		DSN string
	}

	type repo struct {
		DB *db
	}

	type service struct {
		Repo *repo
	}

	c := simpledi.NewContainer()

	c.Register("db", nil, func() any {
		fmt.Println("db created using dsn")
		return &db{DSN: "example"}
	})
	c.Register("repo", []string{"db"}, func() any {
		db := c.Get("db").(*db)
		fmt.Printf("repo created using db(dsn: %s)\n", db.DSN)
		return &repo{DB: db}
	})
	c.Register("service", []string{"repo"}, func() any {
		repo := c.Get("repo").(*repo)
		fmt.Println("service created using repo")
		return &service{Repo: repo}
	})

	if err := c.Resolve(); err != nil {
		panic(err)
	}

	fmt.Println("resolved")

	// Output:
	// db created using dsn
	// repo created using db(dsn: example)
	// service created using repo
	// resolved
}
