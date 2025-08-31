package simpledi

import "fmt"

// ErrorType represents different types of dependency injection errors.
type ErrorType int

const (
	// ErrEmptyKey indicates that dependency key is empty.
	ErrEmptyKey ErrorType = iota
	// ErrNilCtor indicates that constructor function is nil.
	ErrNilCtor
	// ErrMissingDep indicates that required dependency is not registered.
	ErrMissingDep
	// ErrCyclicDeps indicates that circular dependency detected.
	ErrCyclicDeps
	// ErrNotFound indicates that dependency is not found in container.
	ErrNotFound
	// ErrWrongType indicates that dependency cannot be cast to expected type.
	ErrWrongType
)

// Error represents a dependency injection error with specific type and message.
type Error struct {
	errType ErrorType
	message string
}

// Type returns the error type.
func (e *Error) Type() ErrorType {
	return e.errType
}

// Error returns the error message.
func (e *Error) Error() string {
	return e.message
}

func errEmptyKey() *Error {
	return &Error{
		errType: ErrEmptyKey,
		message: "dependency key cannot be empty",
	}
}

func errNilCtor(key string) *Error {
	return &Error{
		errType: ErrNilCtor,
		message: fmt.Sprintf("constructor for dependency [%s] cannot be nil", key),
	}
}

func errMissingDep(key, dep string) *Error {
	return &Error{
		errType: ErrMissingDep,
		message: fmt.Sprintf("dependency [%s] required by [%s] is not registered", dep, key),
	}
}

func errCyclicDeps(deps []string) *Error {
	formatted := make([]string, len(deps))
	for i, dep := range deps {
		formatted[i] = fmt.Sprintf("[%s]", dep)
	}
	return &Error{
		errType: ErrCyclicDeps,
		message: fmt.Sprintf("cyclic dependency detected among: %v", formatted),
	}
}

func errNotFound(key string) *Error {
	return &Error{
		errType: ErrNotFound,
		message: fmt.Sprintf("dependency [%s] not found", key),
	}
}

func errWrongType(key string, expected, actual any) *Error {
	return &Error{
		errType: ErrWrongType,
		message: fmt.Sprintf("dependency [%s] cannot be cast from [%T] to [%T]", key, actual, expected),
	}
}
