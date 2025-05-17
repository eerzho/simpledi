# SimpleDI

[![Release](https://img.shields.io/github/release/eerzho/simpledi.svg?style=for-the-badge)](https://github.com/eerzho/simpledi/releases/latest)
[![Software License](https://img.shields.io/github/license/eerzho/simpledi.svg?style=for-the-badge)](https://github.com/eerzho/simpledi/blob/main/LICENSE)
[![Go docs](https://img.shields.io/badge/go-reference-blue.svg?style=for-the-badge)](https://pkg.go.dev/github.com/eerzho/simpledi)

A simple dependency injection container for Go — zero dependencies, no reflection, no code generation.

###### Installation

```bash
go get github.com/eerzho/simpledi
```

###### Getting started

```go
type repo struct {
    DSN string
}
type service struct {
    repo *repo
}

// create container
c := simpledi.NewContainer()

// register dependencies
c.Register("repo", nil, func() any {
    return &repo{DSN: "example"}
})
c.Register("service", []string{"repo"}, func() any {
    return &service{repo: c.Get("repo").(*repo)}
})

// resolve all dependencies
err := c.Resolve()
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
* `NewContainer`: creates and returns a new DI container
* `Register`: register a dependency by key
* `Resolve`: resolve all dependencies
* `Get`: get a dependency by key

### Documentation and examples

Examples are live in [pkg.go.dev](https://pkg.go.dev/github.com/eerzho/simpledi)
and also in the [example file](./example_test.go).

## Current state

`simpledi` provides the core features intended.

Further improvements or new features may be introduced as good ideas come up.

Suggestions and feedback are always welcome.

## Stargazers over time

[![Stargazers over time](https://starchart.cc/eerzho/simpledi.svg?background=%23FFFFFF&axis=%23333333&line=%236b63ff)](https://starchart.cc/eerzho/simpledi)
