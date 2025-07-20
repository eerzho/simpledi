package simpledi

import "fmt"

type ErrorType int

const (
	ErrEmptyKey ErrorType = iota
	ErrNilCtor
	ErrMissingDep
	ErrCyclicDeps
	ErrNotFound
	ErrWrongType
)

type Error struct {
	Type    ErrorType
	message string
}

func (e *Error) Error() string {
	return e.message
}

func errEmptyKey() *Error {
	return &Error{
		Type:    ErrEmptyKey,
		message: "dependency key cannot be empty",
	}
}

func errNilCtor(key string) *Error {
	return &Error{
		Type:    ErrNilCtor,
		message: fmt.Sprintf("constructor for dependency [%s] cannot be nil", key),
	}
}

func errMissingDep(key, dep string) *Error {
	return &Error{
		Type:    ErrMissingDep,
		message: fmt.Sprintf("dependency [%s] required by [%s] is not registered", dep, key),
	}
}

func errCyclicDeps(deps []string) *Error {
	formatted := make([]string, len(deps))
	for i, dep := range deps {
		formatted[i] = fmt.Sprintf("[%s]", dep)
	}
	return &Error{
		Type:    ErrCyclicDeps,
		message: fmt.Sprintf("cyclic dependency detected among: %v", formatted),
	}
}

func errNotFound(key string) *Error {
	return &Error{
		Type:    ErrNotFound,
		message: fmt.Sprintf("dependency [%s] not found", key),
	}
}

func errWrongType(key string, expected, actual any) *Error {
	return &Error{
		Type:    ErrWrongType,
		message: fmt.Sprintf("dependency [%s] cannot be cast from [%T] to [%T]", key, actual, expected),
	}
}
