package simpledi_test

import (
	"fmt"

	"github.com/eerzho/simpledi"
)

func ExampleContainer() {
	c := simpledi.NewContainer()

	// Register dependencies
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
		fmt.Printf(
			"bread created using: [%s, %s]\n",
			c.Get("yeast"),
			c.Get("flour"),
		)
		return "bread"
	})

	c.Register("sandwich", []string{"bread", "meat"}, func() any {
		fmt.Printf(
			"sandwich created using: [%s, %s]\n",
			c.Get("bread"),
			c.Get("meat"),
		)
		return "sandwich"
	})

	c.Register("burger", []string{"sandwich", "meat", "bread"}, func() any {
		fmt.Printf(
			"burger created using: [%s, %s, %s]\n",
			c.Get("sandwich"),
			c.Get("meat"),
			c.Get("bread"),
		)
		return "burger"
	})

	// Resolve all dependencies
	if err := c.Resolve(); err != nil {
		panic(err)
	}

	fmt.Println("resolved")

	// Output:
	// yeast created
	// flour created
	// meat created
	// bread created using: [yeast, flour]
	// sandwich created using: [bread, meat]
	// burger created using: [sandwich, meat, bread]
	// resolved
}

func ExampleContainer_withServices() {
	type DB struct {
		DSN string
	}

	type Repo1 struct {
		DB *DB
	}

	type Repo2 struct {
		DB *DB
	}

	type Service struct {
		Repo1 *Repo1
		Repo2 *Repo2
	}

	type UseCase struct {
		DB      *DB
		Service *Service
	}

	c := simpledi.NewContainer()

	// Register dependencies
	c.Register("db", nil, func() any {
		fmt.Println("db created")
		return &DB{DSN: "example"}
	})

	c.Register("repo1", []string{"db"}, func() any {
		fmt.Println("repo1 created using: [db]")
		return &Repo1{
			DB: c.Get("db").(*DB),
		}
	})

	c.Register("repo2", []string{"db"}, func() any {
		fmt.Println("repo2 created using: [db]")
		return &Repo2{
			DB: c.Get("db").(*DB),
		}
	})

	c.Register("service", []string{"repo1", "repo2"}, func() any {
		fmt.Println("service created using: [repo1, repo2]")
		return &Service{
			Repo1: c.Get("repo1").(*Repo1),
			Repo2: c.Get("repo2").(*Repo2),
		}
	})

	c.Register("usecase", []string{"db", "service"}, func() any {
		fmt.Println("usecase created using: [db, service]")
		return &UseCase{
			DB:      c.Get("db").(*DB),
			Service: c.Get("service").(*Service),
		}
	})

	// Resolve all dependencies
	if err := c.Resolve(); err != nil {
		panic(err)
	}

	fmt.Println("resolved")

	// Output:
	// db created
	// repo1 created using: [db]
	// repo2 created using: [db]
	// service created using: [repo1, repo2]
	// usecase created using: [db, service]
	// resolved
}
