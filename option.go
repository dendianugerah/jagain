// Package jagain provides null safety utilities for Go.
// It helps express optionality through the type system rather than using nil pointers.
package jagain

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ErrNoValue is returned when attempting to access a value that is not present.
var ErrNoValue = errors.New("option contains no value")

// Option represents a value that may or may not be present.
type Option[T any] struct {
	value *T
	valid bool
}

// Some creates an Option containing a value.
func Some[T any](value T) Option[T] {
	return Option[T]{
		value: &value,
		valid: true,
	}
}

// None creates an Option with no value.
func None[T any]() Option[T] {
	return Option[T]{
		value: nil,
		valid: false,
	}
}

// FromPtr creates an Option from a pointer.
// If the pointer is nil, None is returned.
// Otherwise, Some is returned with the dereferenced value.
func FromPtr[T any](ptr *T) Option[T] {
	if ptr == nil {
		return None[T]()
	}
	return Some(*ptr)
}

// ToPtr converts an Option to a pointer.
// If the Option has no value, nil is returned.
// Otherwise, a pointer to the value is returned.
func (o Option[T]) ToPtr() *T {
	if !o.valid {
		return nil
	}
	v := *o.value
	return &v
}

// IsSome returns true if the Option contains a value.
func (o Option[T]) IsSome() bool {
	return o.valid
}

// IsNone returns true if the Option does not contain a value.
func (o Option[T]) IsNone() bool {
	return !o.valid
}

// Unwrap returns the contained value or panics if no value is present.
// This should be used only when you are confident a value is present.
func (o Option[T]) Unwrap() T {
	if !o.valid {
		panic(ErrNoValue)
	}
	return *o.value
}

// UnwrapOr returns the contained value or the provided default if no value is present.
func (o Option[T]) UnwrapOr(defaultValue T) T {
	if !o.valid {
		return defaultValue
	}
	return *o.value
}

// UnwrapOrElse returns the contained value or computes a value from the provided function.
func (o Option[T]) UnwrapOrElse(f func() T) T {
	if !o.valid {
		return f()
	}
	return *o.value
}

// Map transforms the Option's value using the provided function if a value is present.
func (o Option[T]) Map(f func(T) T) Option[T] {
	if !o.valid {
		return o
	}
	return Some(f(*o.value))
}

// FlatMap transforms the Option's value into another Option using the provided function.
func (o Option[T]) FlatMap(f func(T) Option[T]) Option[T] {
	if !o.valid {
		return o
	}
	return f(*o.value)
}

// Match pattern-matches on the Option, applying one of two functions.
func (o Option[T]) Match(some func(T) T, none func() T) T {
	if o.valid {
		return some(*o.value)
	}
	return none()
}

// ToResult converts an Option to a Result.
// If the Option contains a value, Ok is returned.
// If the Option does not contain a value, Err is returned with the provided error.
func (o Option[T]) ToResult(err error) Result[T] {
	if !o.valid {
		return Err[T](err)
	}
	return Ok(*o.value)
}

// MarshalJSON implements the json.Marshaler interface.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if !o.valid {
		return []byte("null"), nil
	}
	return json.Marshal(o.value)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*o = None[T]()
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	*o = Some(value)
	return nil
}

// String implements the fmt.Stringer interface.
func (o Option[T]) String() string {
	if !o.valid {
		return "None"
	}
	return fmt.Sprintf("Some(%v)", *o.value)
}
