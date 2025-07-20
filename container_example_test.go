package simpledi_test

import (
	"errors"
	"fmt"

	"github.com/eerzho/simpledi"
)

func Example() {
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

	// creating a container
	c := simpledi.NewContainer()

	// registration of dependencies
	c.MustRegister(simpledi.Option{
		Key: "database",
		Ctor: func() any {
			fmt.Println("creating database...")
			return &Database{url: "real_url"}
		},
	})
	c.MustRegister(simpledi.Option{
		Key:  "user_service",
		Deps: []string{"database"},
		Ctor: func() any {
			fmt.Println("creating userService...")
			db := c.MustGet("database").(*Database)
			return &UserService{db: db}
		},
	})
	c.MustRegister(simpledi.Option{
		Key:  "order_service",
		Deps: []string{"database", "user_service"},
		Ctor: func() any {
			fmt.Println("creating orderService...")
			db := c.MustGet("database").(*Database)
			userService := c.MustGet("user_service").(*UserService)
			return &OrderService{db: db, userService: userService}
		},
	})

	// resolving dependencies
	c.MustResolve()

	// getting dependencies
	userService := c.MustGet("user_service").(*UserService)
	orderService := c.MustGet("order_service").(*OrderService)

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

func Example_anyOrder() {
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

	// creating a container
	c := simpledi.NewContainer()

	// registration of dependencies
	c.MustRegister(simpledi.Option{
		Key:  "order_service",
		Deps: []string{"database", "user_service"},
		Ctor: func() any {
			fmt.Println("creating orderService...")
			db := c.MustGet("database").(*Database)
			userService := c.MustGet("user_service").(*UserService)
			return &OrderService{db: db, userService: userService}
		},
	})
	c.MustRegister(simpledi.Option{
		Key:  "user_service",
		Deps: []string{"database"},
		Ctor: func() any {
			fmt.Println("creating userService...")
			db := c.MustGet("database").(*Database)
			return &UserService{db: db}
		},
	})
	c.MustRegister(simpledi.Option{
		Key: "database",
		Ctor: func() any {
			fmt.Println("creating database...")
			return &Database{url: "real_url"}
		},
	})

	// resolving dependencies
	c.MustResolve()

	// getting dependencies
	userService := c.MustGet("user_service").(*UserService)
	orderService := c.MustGet("order_service").(*OrderService)

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

func Example_override() {
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

	// creating a container
	c := simpledi.NewContainer()

	// registration of dependencies
	c.MustRegister(simpledi.Option{
		Key: "database",
		Ctor: func() any {
			fmt.Println("creating database...")
			return &Database{url: "real_url"}
		},
	})
	c.MustRegister(simpledi.Option{
		Key:  "user_service",
		Deps: []string{"database"},
		Ctor: func() any {
			fmt.Println("creating userService...")
			db := c.MustGet("database").(*Database)
			return &UserService{db: db}
		},
	})
	c.MustRegister(simpledi.Option{
		Key:  "order_service",
		Deps: []string{"database", "user_service"},
		Ctor: func() any {
			fmt.Println("creating orderService...")
			db := c.MustGet("database").(*Database)
			userService := c.MustGet("user_service").(*UserService)
			return &OrderService{db: db, userService: userService}
		},
	})

	// resolving dependencies
	c.MustResolve()

	// getting dependencies
	userService := c.MustGet("user_service").(*UserService)
	orderService := c.MustGet("order_service").(*OrderService)

	fmt.Printf("userService db url: %s\n", userService.db.url)
	fmt.Printf("orderService db url: %s\n", orderService.db.url)
	fmt.Printf("same database instance used: %t\n", userService.db == orderService.db)
	fmt.Printf("orderService has userService: %t\n", orderService.userService == userService)

	// overriding database
	c.Register(simpledi.Option{
		Key: "database",
		Ctor: func() any {
			fmt.Println("creating fake database...")
			return &Database{url: "fake_url"}
		},
	})

	// resolving dependencies
	c.MustResolve()

	// getting dependencies
	userService = c.MustGet("user_service").(*UserService)
	orderService = c.MustGet("order_service").(*OrderService)

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
	// creating fake database...
	// creating userService...
	// creating orderService...
	// userService db url: fake_url
	// orderService db url: fake_url
	// same database instance used: true
	// orderService has userService: true
}

func Example_withDestructor() {
	// example of structures
	type Database struct {
		connected bool
	}

	// creating a container
	c := simpledi.NewContainer()

	// registration of dependencies
	c.MustRegister(simpledi.Option{
		Key: "database",
		Ctor: func() any {
			fmt.Println("connecting to database...")
			return &Database{connected: true}
		},
		Dtor: func() error {
			fmt.Println("disconnecting from database...")
			db := c.MustGet("database").(*Database)
			db.connected = false
			return nil
		},
	})

	// resolving dependencies
	c.MustResolve()

	// getting dependencies
	db := c.MustGet("database").(*Database)

	fmt.Printf("database connected: %t\n", db.connected)

	// resetting container
	c.MustReset()

	fmt.Printf("database connected: %t\n", db.connected)

	// Output:
	// connecting to database...
	// database connected: true
	// disconnecting from database...
	// database connected: false
}

func Example_errorHandling() {
	// example of structures
	type Database struct {
		url string
	}

	// creating a container
	c := simpledi.NewContainer()

	// registration of dependencies
	err := c.Register(simpledi.Option{
		Ctor: func() any {
			return &Database{}
		},
	})

	// parsing error using simpledi.Error structure
	var diErr *simpledi.Error
	if errors.As(err, &diErr) {
		fmt.Printf("registration failed: %s\n", diErr.Error())
		fmt.Printf("type ErrEmptyKey: %t\n", diErr.Type == simpledi.ErrEmptyKey)
	}

	// Output:
	// registration failed: dependency key cannot be empty
	// type ErrEmptyKey: true
}
