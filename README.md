# simpledi

[![Release](https://img.shields.io/github/v/release/eerzho/simpledi?sort=semver)](https://github.com/eerzho/simpledi/releases/latest)
[![License](https://img.shields.io/github/license/eerzho/simpledi)](https://github.com/eerzho/simpledi/blob/main/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/eerzho/simpledi.svg)](https://pkg.go.dev/github.com/eerzho/simpledi@latest)
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

type UserService struct {
	DB *Database
}

func main() {
	defer simpledi.Close()

	simpledi.Set(simpledi.Definition{
		ID: "database",
		New: func() any {
			return &Database{URL: "real_url"}
		},
	})

	simpledi.Set(simpledi.Definition{
		ID:   "user_service",
		Deps: []string{"database"},
		New: func() any {
			db := simpledi.Get[*Database]("database")
			return &UserService{DB: db}
		},
	})

	simpledi.Resolve()

	userService := simpledi.Get[*UserService]("user_service")
	_ = userService
}
```

## Features

* Zero dependencies
* No reflection
* Type-safe with generics
* Automatic dependency ordering
* Optional cleanup with `Close`
* Global container with `Set`, `Get`, `Resolve`, `Close`

## Docs

Full documentation and examples:
[https://pkg.go.dev/github.com/eerzho/simpledi](https://pkg.go.dev/github.com/eerzho/simpledi)
