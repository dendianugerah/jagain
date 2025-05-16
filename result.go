package jagain

import (
	"fmt"
)

// Result represents either a success value or an error.
// It's similar to Rust's Result type.
type Result[T any] struct {
	value *T
	err   error
	valid bool
}

// Ok creates a Result containing a success value.
func Ok[T any](value T) Result[T] {
	return Result[T]{
		value: &value,
		err:   nil,
		valid: true,
	}
}

// Err creates a Result containing an error.
func Err[T any](err error) Result[T] {
	return Result[T]{
		value: nil,
		err:   err,
		valid: false,
	}
}

// IsOk returns true if the Result contains a success value.
func (r Result[T]) IsOk() bool {
	return r.valid
}

// IsErr returns true if the Result contains an error.
func (r Result[T]) IsErr() bool {
	return !r.valid
}

// Unwrap returns the contained success value or panics if the Result contains an error.
func (r Result[T]) Unwrap() T {
	if !r.valid {
		panic(fmt.Sprintf("called unwrap on an error result: %v", r.err))
	}
	return *r.value
}

// UnwrapOr returns the contained success value or the provided default if the Result contains an error.
func (r Result[T]) UnwrapOr(defaultValue T) T {
	if !r.valid {
		return defaultValue
	}
	return *r.value
}

// UnwrapOrElse returns the contained success value or computes a value from the provided function.
func (r Result[T]) UnwrapOrElse(f func(error) T) T {
	if !r.valid {
		return f(r.err)
	}
	return *r.value
}

// UnwrapErr returns the contained error or panics if the Result contains a success value.
func (r Result[T]) UnwrapErr() error {
	if r.valid {
		panic("called unwrap_err on an ok result")
	}
	return r.err
}

// Map transforms the Result's success value using the provided function.
// If the Result contains an error, it is returned unchanged.
func (r Result[T]) Map(f func(T) T) Result[T] {
	if !r.valid {
		return r
	}
	return Ok(f(*r.value))
}

// MapTo transforms the Result's success value into a different type using the provided function.
// If the Result contains an error, an error result of the new type is returned.
func MapTo[T, U any](r Result[T], f func(T) U) Result[U] {
	if !r.valid {
		return Err[U](r.err)
	}
	return Ok(f(*r.value))
}

// MapErr transforms the Result's error using the provided function.
// If the Result contains a success value, it is returned unchanged.
func (r Result[T]) MapErr(f func(error) error) Result[T] {
	if r.valid {
		return r
	}
	return Err[T](f(r.err))
}

// FlatMap transforms the Result's success value into another Result of the same type using the provided function.
// If the Result contains an error, it is returned unchanged.
func (r Result[T]) FlatMap(f func(T) Result[T]) Result[T] {
	if !r.valid {
		return r
	}
	return f(*r.value)
}

// FlatMapTo transforms the Result's success value into a Result of a different type.
// If the Result contains an error, an error result of the new type is returned.
func FlatMapTo[T, U any](r Result[T], f func(T) Result[U]) Result[U] {
	if !r.valid {
		return Err[U](r.err)
	}
	return f(*r.value)
}

// Match pattern-matches on the Result, applying one of two functions.
func (r Result[T]) Match(ok func(T) T, err func(error) T) T {
	if r.valid {
		return ok(*r.value)
	}
	return err(r.err)
}

// MatchTo pattern-matches on the Result, applying one of two functions that return a different type.
func MatchTo[T, U any](r Result[T], ok func(T) U, err func(error) U) U {
	if r.valid {
		return ok(*r.value)
	}
	return err(r.err)
}

// ToOption converts a Result to an Option.
// If the Result contains a success value, Some is returned.
// If the Result contains an error, None is returned.
func (r Result[T]) ToOption() Option[T] {
	if !r.valid {
		return None[T]()
	}
	return Some(*r.value)
}

// String implements the fmt.Stringer interface.
func (r Result[T]) String() string {
	if !r.valid {
		return fmt.Sprintf("Err(%v)", r.err)
	}
	return fmt.Sprintf("Ok(%v)", *r.value)
}
