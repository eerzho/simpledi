# SimpleDI

[English](README.md) | [Русский](README.ru.md)

SimpleDI is a lightweight dependency injection container for Go applications. It provides a simple way to manage dependencies and their lifecycle in your application.

### Features

- Simple and intuitive API
- Dependency resolution with automatic ordering
- Cyclic dependency detection
- Type-safe dependency injection
- No external dependencies

### Installation

```bash
go get github.com/eerzho/simpledi@latest
```

### Quick Start

```go
package simpledi

import "github.com/eerzho/simpledi"

func main() {
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
}
```

### API Reference

#### NewContainer()

Creates a new dependency injection container.

#### Register(name string, deps []string, constructor func() any)

Registers a new dependency with the container.
- `name`: Unique identifier for the dependency
- `deps`: List of dependency names this component depends on
- `constructor`: Function that creates the dependency instance

#### Get(name string) any

Retrieves a resolved dependency instance by its name.

#### Resolve() error

Resolves all registered dependencies in the correct order. Returns an error if there are cyclic dependencies or missing dependencies.

### License

MIT License - see [LICENSE](LICENSE) file for details
