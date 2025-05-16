package jagain

import (
	"encoding/json"
	"testing"
)

func TestOption(t *testing.T) {
	// Create Some with a value
	s := Some(42)
	if !s.IsSome() {
		t.Errorf("Expected Some to be Some, got None")
	}
	if s.IsNone() {
		t.Errorf("Expected Some not to be None")
	}
	if s.Unwrap() != 42 {
		t.Errorf("Expected Some.Unwrap() to be 42, got %v", s.Unwrap())
	}

	// Create None
	n := None[int]()
	if n.IsSome() {
		t.Errorf("Expected None not to be Some")
	}
	if !n.IsNone() {
		t.Errorf("Expected None to be None")
	}

	// Test UnwrapOr
	if n.UnwrapOr(10) != 10 {
		t.Errorf("Expected UnwrapOr to return default value")
	}
	if s.UnwrapOr(10) != 42 {
		t.Errorf("Expected UnwrapOr to return contained value")
	}

	// Test UnwrapOrElse
	if n.UnwrapOrElse(func() int { return 20 }) != 20 {
		t.Errorf("Expected UnwrapOrElse to execute function")
	}
	if s.UnwrapOrElse(func() int { return 20 }) != 42 {
		t.Errorf("Expected UnwrapOrElse to return contained value")
	}

	// Test Map
	mapped := s.Map(func(i int) int { return i * 2 })
	if mapped.Unwrap() != 84 {
		t.Errorf("Expected Map to transform value")
	}

	// Test FlatMap
	flatMapped := s.FlatMap(func(i int) Option[int] {
		if i > 0 {
			return Some(i * 3)
		}
		return None[int]()
	})
	if flatMapped.Unwrap() != 126 {
		t.Errorf("Expected FlatMap to transform value")
	}

	// Test Match
	matchResult := s.Match(
		func(i int) int { return i + 1 },
		func() int { return 0 },
	)
	if matchResult != 43 {
		t.Errorf("Expected Match to apply 'some' function")
	}

	noneMatchResult := n.Match(
		func(i int) int { return i + 1 },
		func() int { return 0 },
	)
	if noneMatchResult != 0 {
		t.Errorf("Expected Match to apply 'none' function")
	}

	// Test FromPtr
	var ptr *int
	ptrOption := FromPtr(ptr)
	if !ptrOption.IsNone() {
		t.Errorf("Expected FromPtr with nil to be None")
	}

	val := 42
	ptr = &val
	ptrOption = FromPtr(ptr)
	if !ptrOption.IsSome() || ptrOption.Unwrap() != 42 {
		t.Errorf("Expected FromPtr with non-nil to be Some")
	}

	// Test ToPtr
	nonePtr := n.ToPtr()
	if nonePtr != nil {
		t.Errorf("Expected ToPtr on None to return nil")
	}

	somePtr := s.ToPtr()
	if somePtr == nil || *somePtr != 42 {
		t.Errorf("Expected ToPtr on Some to return non-nil pointer to value")
	}

	// Test String
	if s.String() != "Some(42)" {
		t.Errorf("Expected s.String() to be 'Some(42)', got '%s'", s.String())
	}
	if n.String() != "None" {
		t.Errorf("Expected n.String() to be 'None', got '%s'", n.String())
	}
}

func TestOptionJSON(t *testing.T) {
	// Test marshaling Some
	s := Some("hello")
	bytes, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("Failed to marshal Some: %v", err)
	}
	if string(bytes) != `"hello"` {
		t.Errorf("Expected marshaled Some to be '\"hello\"', got '%s'", string(bytes))
	}

	// Test marshaling None
	n := None[string]()
	bytes, err = json.Marshal(n)
	if err != nil {
		t.Fatalf("Failed to marshal None: %v", err)
	}
	if string(bytes) != "null" {
		t.Errorf("Expected marshaled None to be 'null', got '%s'", string(bytes))
	}

	// Test unmarshaling to Some
	var opt Option[string]
	err = json.Unmarshal([]byte(`"world"`), &opt)
	if err != nil {
		t.Fatalf("Failed to unmarshal to Some: %v", err)
	}
	if !opt.IsSome() || opt.Unwrap() != "world" {
		t.Errorf("Expected unmarshaled value to be Some(\"world\")")
	}

	// Test unmarshaling to None
	err = json.Unmarshal([]byte("null"), &opt)
	if err != nil {
		t.Fatalf("Failed to unmarshal to None: %v", err)
	}
	if !opt.IsNone() {
		t.Errorf("Expected unmarshaled value to be None")
	}
}
