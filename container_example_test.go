package simpledi_test

import (
	"fmt"

	"github.com/eerzho/simpledi"
)

// Example of a simple service
type Database struct {
	connectionString string
}

func NewDatabase() *Database {
	return &Database{connectionString: "localhost:5432"}
}

// Example of a service that depends on another service
type UserService struct {
	db *Database
}

func NewUserService(db *Database) *UserService {
	return &UserService{db: db}
}

func ExampleContainer() {
	// Create a new container
	container := simpledi.NewContainer()

	// Register services
	container.Register("database", nil, func() any {
		return NewDatabase()
	})

	container.Register("userService", []string{"database"}, func() any {
		db := container.Get("database").(*Database)
		return NewUserService(db)
	})

	// Resolve all dependencies
	container.Resolve()

	// Get the resolved services
	db := container.Get("database").(*Database)
	userService := container.Get("userService").(*UserService)

	fmt.Printf("Database connection: %s\n", db.connectionString)
	fmt.Printf("UserService has database: %v\n", userService.db != nil)

	// Output:
	// Database connection: localhost:5432
	// UserService has database: true
}

func ExampleContainer_chainedRegistration() {
	container := simpledi.NewContainer()

	// Chain registration calls
	container.
		Register("service1", nil, func() any { return "service1" }).
		Register("service2", []string{"service1"}, func() any { return "service2" }).
		Register("service3", []string{"service2"}, func() any { return "service3" })

	container.Resolve()

	fmt.Println(container.Get("service1"))
	fmt.Println(container.Get("service2"))
	fmt.Println(container.Get("service3"))

	// Output:
	// service1
	// service2
	// service3
}
