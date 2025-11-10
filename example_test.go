package simpledi_test

import (
	"errors"
	"fmt"

	"github.com/eerzho/simpledi"
)

type Config struct {
	CacheURL    string
	DatabaseURL string
}

type Cache struct {
	URL string
}

func (c *Cache) Close() error {
	return errors.New("cache close error")
}

type Database struct {
	URL string
}

func (d *Database) Close() error {
	return errors.New("database close error")
}

type Service struct {
	Cache    *Cache
	Database *Database
}

// Example demonstrates basic usage.
func Example() {
	// Close all definitions in reverse order at the end
	defer simpledi.Close()

	// Define configuration
	simpledi.Set(simpledi.Definition{
		ID: "config",
		New: func() any {
			fmt.Println("Creating config")
			return &Config{
				CacheURL:    "redis",
				DatabaseURL: "postgres",
			}
		},
	})

	// Define cache
	simpledi.Set(simpledi.Definition{
		ID:   "cache",
		Deps: []string{"config"},
		New: func() any {
			fmt.Println("Creating cache")
			config := simpledi.Get[*Config]("config")
			return &Cache{URL: config.CacheURL}
		},
		Close: func() error {
			fmt.Println("Closing cache")
			cache := simpledi.Get[*Cache]("cache")
			return cache.Close()
		},
	})

	// Define database
	simpledi.Set(simpledi.Definition{
		ID:   "database",
		Deps: []string{"config"},
		New: func() any {
			fmt.Println("Creating database")
			config := simpledi.Get[*Config]("config")
			return &Database{URL: config.DatabaseURL}
		},
		Close: func() error {
			fmt.Println("Closing database")
			database := simpledi.Get[*Database]("database")
			return database.Close()
		},
	})

	// Define service
	simpledi.Set(simpledi.Definition{
		ID:   "service",
		Deps: []string{"cache", "database"},
		New: func() any {
			fmt.Println("Creating service")
			cache := simpledi.Get[*Cache]("cache")
			database := simpledi.Get[*Database]("database")
			return &Service{Cache: cache, Database: database}
		},
	})

	// Resolve all dependencies in correct order
	simpledi.Resolve()

	// Use the service
	service := simpledi.Get[*Service]("service")
	fmt.Printf("Service uses cache: %s\n", service.Cache.URL)
	fmt.Printf("Service uses database: %s\n", service.Database.URL)

	// Output:
	// Creating config
	// Creating cache
	// Creating database
	// Creating service
	// Service uses cache: redis
	// Service uses database: postgres
	// Closing database
	// Closing cache
}
