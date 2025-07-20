package simpledi_test

import (
	"fmt"

	"github.com/eerzho/simpledi"
)

func Example_default() {
	// example of structures
	type Database struct {
		url string
	}
	type UserService struct {
		db *Database
	}
	type OrderService struct {
		db          *Database
		userService *UserService
	}

	// registration of dependencies
	simpledi.MustRegister(simpledi.Option{
		Key: "database",
		Ctor: func() any {
			fmt.Println("creating database...")
			return &Database{url: "real_url"}
		},
	})
	simpledi.MustRegister(simpledi.Option{
		Key:  "user_service",
		Deps: []string{"database"},
		Ctor: func() any {
			fmt.Println("creating userService...")
			db := simpledi.MustGetAs[*Database]("database")
			return &UserService{db: db}
		},
	})
	simpledi.MustRegister(simpledi.Option{
		Key:  "order_service",
		Deps: []string{"database", "user_service"},
		Ctor: func() any {
			fmt.Println("creating orderService...")
			db := simpledi.MustGetAs[*Database]("database")
			userService := simpledi.MustGetAs[*UserService]("user_service")
			return &OrderService{db: db, userService: userService}
		},
	})

	// resolving dependencies
	simpledi.MustResolve()

	// getting dependencies
	userService := simpledi.MustGetAs[*UserService]("user_service")
	orderService := simpledi.MustGetAs[*OrderService]("order_service")

	fmt.Printf("userService db url: %s\n", userService.db.url)
	fmt.Printf("orderService db url: %s\n", orderService.db.url)
	fmt.Printf("same database instance used: %t\n", userService.db == orderService.db)
	fmt.Printf("orderService has userService: %t\n", orderService.userService == userService)

	// Output:
	// creating database...
	// creating userService...
	// creating orderService...
	// userService db url: real_url
	// orderService db url: real_url
	// same database instance used: true
	// orderService has userService: true
}
