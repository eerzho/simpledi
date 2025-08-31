# SimpleDI

[![Release](https://img.shields.io/github/release/eerzho/simpledi.svg)](https://github.com/eerzho/simpledi/releases/latest)
[![License](https://img.shields.io/github/license/eerzho/simpledi.svg)](https://github.com/eerzho/simpledi/blob/main/LICENSE)
[![Go Reference](https://img.shields.io/badge/go-reference-blue.svg)](https://pkg.go.dev/github.com/eerzho/simpledi)
[![Go Report](https://goreportcard.com/badge/github.com/eerzho/simpledi)](https://goreportcard.com/report/github.com/eerzho/simpledi)
[![Codecov](https://codecov.io/gh/eerzho/simpledi/branch/main/graph/badge.svg)](https://codecov.io/gh/eerzho/simpledi)

A simple dependency injection container for Go — zero dependencies, no reflection, no code generation.

###### Installation

```bash
go get github.com/eerzho/simpledi
```

###### Getting started

```go
type Database struct {
	url string
}
type UserService struct {
	db *Database
}

// registration of dependencies
simpledi.MustRegister(simpledi.Def{
	Key: "database",
	Ctor: func() any {
		return &Database{url: "real_url"}
	},
})
simpledi.MustRegister(simpledi.Def{
	Key:  "user_service",
	Deps: []string{"database"},
	Ctor: func() any {
		db := simpledi.MustGetAs[*Database]("database")
		return &UserService{db: db}
	},
})

// resolving dependencies
simpledi.MustResolve()

// getting dependencies
userService := simpledi.MustGetAs[*UserService]("user_service")
```

You can see the full documentation and list of examples at [pkg.go.dev](https://pkg.go.dev/github.com/eerzho/simpledi#pkg-examples).

---

## Usage

### Key Features

* **Zero dependencies** - Pure Go, no external packages
* **No reflection** - Type-safe with generics support
* **Dependency ordering** - Register in any order, automatic resolution
* **Lifecycle management** - Constructor and destructor support
* **Thread-safe** - Protected by mutex for concurrent access
* **Global container** - Use package-level functions without creating container

### Best Practices

* **Use descriptive keys**: `userService` not `service`
* **Register then resolve**: Complete all registrations before calling `Resolve()`
* **Leverage type safety**: Use `MustGetAs[T]()`/`GetAs[T]()` instead of type assertions
* **Handle cleanup**: Implement `Dtor` for resources that need cleanup
* **Test with mocks**: Override dependencies for testing

### Documentation and examples

Examples are live in [pkg.go.dev](https://pkg.go.dev/github.com/eerzho/simpledi#pkg-examples)
and also in the [example file](./container_example_test.go).

## Current state

`simpledi` provides the core features intended.

Further improvements or new features may be introduced as good ideas come up.

Suggestions and feedback are always welcome.

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=eerzho/simpledi&type=Timeline)](https://www.star-history.com/#eerzho/simpledi&Timeline)
