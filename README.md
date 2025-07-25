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
type repo struct {
	dsn string
}
type service struct {
	repo *repo
}

// create container
c := simpledi.NewContainer()

// register dependencies
c.Register("repo", nil, func() any {
	return &repo{dsn: "example"}
})
c.Register("service", []string{"repo"}, func() any {
	return &service{repo: c.Get("repo").(*repo)}
})

// resolve all dependencies
c.Resolve()
```

You can see the full documentation and list of examples at [pkg.go.dev](https://pkg.go.dev/github.com/eerzho/simpledi).

---

## Usage

### Notes

* You can register dependencies in any order.
* Call `Resolve` only after all registrations are done.
* To recreate all dependencies, call `Resolve` again.
* To override an implementation, register with the same key and call `Resolve` again.

### Functions
* `NewContainer`: creates a new DI container
* `Register`: register a dependency by key
* `Get`: get a dependency by key
* `GetAs`: get a dependency by key and casts it to the specified type
* `MustGet`: get a dependency by key or panics
* `MustGetAs`: get a dependency by key and casts it to the specified type or panics
* `Resolve`: resolve all dependencies
* `MustResolve`: resolve all dependencies or panic

### Documentation and examples

Examples are live in [pkg.go.dev](https://pkg.go.dev/github.com/eerzho/simpledi)
and also in the [example file](./container_example_test.go).

## Current state

`simpledi` provides the core features intended.

Further improvements or new features may be introduced as good ideas come up.

Suggestions and feedback are always welcome.

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=eerzho/simpledi&type=Timeline)](https://www.star-history.com/#eerzho/simpledi&Timeline)

## Alternatives

- [goioc/di](https://github.com/goioc/di)
- [sarulabs/dingo](https://github.com/sarulabs/dingo)
