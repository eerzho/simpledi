# simpledi

[![Release](https://img.shields.io/github/v/release/eerzho/simpledi?sort=semver)](https://github.com/eerzho/simpledi/releases/latest)
[![License](https://img.shields.io/github/license/eerzho/simpledi)](https://github.com/eerzho/simpledi/blob/main/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/eerzho/simpledi.svg)](https://pkg.go.dev/github.com/eerzho/simpledi)
[![Go Report](https://goreportcard.com/badge/github.com/eerzho/simpledi)](https://goreportcard.com/report/github.com/eerzho/simpledi)
[![Codecov](https://img.shields.io/codecov/c/github/eerzho/simpledi?token=YOUR_TOKEN)](https://codecov.io/gh/eerzho/simpledi)

A simple dependency injection container for Go.

Zero dependencies, no reflection, no code generation.

## Install

```bash
go get github.com/eerzho/simpledi
```

## Example

```go
package main

import "github.com/eerzho/simpledi"

type Database struct {
	URL string
}

func (d *Database) Close() error {
	return errors.New("database close error")
}

type Service struct {
	Database *Database
}

func main() {
	// Close all definitions in reverse order at the end
	defer simpledi.Close()

	// Define database
	simpledi.Set(simpledi.Definition{
		ID:   "database",
		New: func() any {
			return &Database{URL: "postgres"}
		},
		Close: func() error {
			database := simpledi.Get[*Database]("database")
			return database.Close()
		},
	})

	// Define service
	simpledi.Set(simpledi.Definition{
		ID:   "service",
		Deps: []string{"database"},
		New: func() any {
			database := simpledi.Get[*Database]("database")
			return &Service{Database: database}
		},
	})

	// Resolve all dependencies in correct order
	simpledi.Resolve()

	// Use the service
	service := simpledi.Get[*Service]("service")
	_ = service
}
```

## Features

* Zero dependencies
* No reflection
* Type-safe with generics
* Automatic dependency ordering

## Docs

Full documentation and examples: [pkg.go.dev](https://pkg.go.dev/github.com/eerzho/simpledi#pkg-examples).

## Current state

`simpledi` provides the core features intended.

Further improvements or new features may be introduced as good ideas come up.

Suggestions and feedback are always welcome.

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=eerzho/simpledi&type=Timeline)](https://www.star-history.com/#eerzho/simpledi&Timeline)
