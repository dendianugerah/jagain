package jagain

import (
	"errors"
	"testing"
)

func TestResult(t *testing.T) {
	// Create Ok with a value
	ok := Ok(42)
	if !ok.IsOk() {
		t.Errorf("Expected Ok to be Ok, got Err")
	}
	if ok.IsErr() {
		t.Errorf("Expected Ok not to be Err")
	}
	if ok.Unwrap() != 42 {
		t.Errorf("Expected Ok.Unwrap() to be 42, got %v", ok.Unwrap())
	}

	// Create Err with an error
	testErr := errors.New("test error")
	err := Err[int](testErr)
	if err.IsOk() {
		t.Errorf("Expected Err not to be Ok")
	}
	if !err.IsErr() {
		t.Errorf("Expected Err to be Err")
	}
	if err.UnwrapErr() != testErr {
		t.Errorf("Expected Err.UnwrapErr() to be the original error")
	}

	// Test UnwrapOr
	if err.UnwrapOr(10) != 10 {
		t.Errorf("Expected UnwrapOr to return default value")
	}
	if ok.UnwrapOr(10) != 42 {
		t.Errorf("Expected UnwrapOr to return contained value")
	}

	// Test UnwrapOrElse
	if err.UnwrapOrElse(func(e error) int { return 20 }) != 20 {
		t.Errorf("Expected UnwrapOrElse to execute function")
	}
	if ok.UnwrapOrElse(func(e error) int { return 20 }) != 42 {
		t.Errorf("Expected UnwrapOrElse to return contained value")
	}

	// Test Map
	mapped := ok.Map(func(i int) int { return i * 2 })
	if mapped.Unwrap() != 84 {
		t.Errorf("Expected Map to transform value")
	}
	// Mapping an error should be a no-op
	mappedErr := err.Map(func(i int) int { return i * 2 })
	if !mappedErr.IsErr() {
		t.Errorf("Expected Map on Err to be a no-op")
	}

	// Test MapErr
	mappedErr = err.MapErr(func(e error) error { return errors.New("new error") })
	if mappedErr.UnwrapErr().Error() != "new error" {
		t.Errorf("Expected MapErr to transform error")
	}
	// MapErr on an Ok should be a no-op
	mappedOk := ok.MapErr(func(e error) error { return errors.New("new error") })
	if !mappedOk.IsOk() || mappedOk.Unwrap() != 42 {
		t.Errorf("Expected MapErr on Ok to be a no-op")
	}

	// Test FlatMap
	flatMapped := ok.FlatMap(func(i int) Result[int] {
		if i > 0 {
			return Ok(i * 3)
		}
		return Err[int](errors.New("negative"))
	})
	if flatMapped.Unwrap() != 126 {
		t.Errorf("Expected FlatMap to transform value")
	}

	// Test Match
	matchResult := ok.Match(
		func(i int) int { return i + 1 },
		func(e error) int { return 0 },
	)
	if matchResult != 43 {
		t.Errorf("Expected Match to apply 'ok' function")
	}

	errMatchResult := err.Match(
		func(i int) int { return i + 1 },
		func(e error) int { return 0 },
	)
	if errMatchResult != 0 {
		t.Errorf("Expected Match to apply 'err' function")
	}

	// Test ToOption
	optFromOk := ok.ToOption()
	if !optFromOk.IsSome() || optFromOk.Unwrap() != 42 {
		t.Errorf("Expected ToOption on Ok to return Some")
	}

	optFromErr := err.ToOption()
	if !optFromErr.IsNone() {
		t.Errorf("Expected ToOption on Err to return None")
	}

	// Test String
	if ok.String() != "Ok(42)" {
		t.Errorf("Expected ok.String() to be 'Ok(42)', got '%s'", ok.String())
	}
	if err.String() != "Err(test error)" {
		t.Errorf("Expected err.String() to be 'Err(test error)', got '%s'", err.String())
	}
}
