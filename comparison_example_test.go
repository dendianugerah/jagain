package jagain

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
)

// This file demonstrates the differences between standard Go error handling
// and using the Result/Option types from the jagain library.

// MockDatabase simulates a database for our examples
type MockDatabase struct {
	users map[int]*UserRecord
}

// UserRecord simulates a database record
type UserRecord struct {
	ID       int
	Username string
	Email    string // Can be empty
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		users: map[int]*UserRecord{
			1: {ID: 1, Username: "alice", Email: "alice@example.com"},
			2: {ID: 2, Username: "bob", Email: ""}, // No email
		},
	}
}

// GetUserByID returns a user by ID using the standard Go error approach
func (db *MockDatabase) GetUserByID(id int) (*UserRecord, error) {
	user, exists := db.users[id]
	if !exists {
		return nil, fmt.Errorf("user with ID %d not found", id)
	}
	return user, nil
}

// GetUserByIDResult returns a user by ID using our Result type
func (db *MockDatabase) GetUserByIDResult(id int) Result[UserRecord] {
	user, exists := db.users[id]
	if !exists {
		return Err[UserRecord](fmt.Errorf("user with ID %d not found", id))
	}
	return Ok(*user)
}

// GetUserEmailByID returns a user's email using standard Go approach
func GetUserEmailByID(db *MockDatabase, id int) (string, error) {
	// Get the user
	user, err := db.GetUserByID(id)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// Check if email exists (empty string means no email)
	if user.Email == "" {
		return "", errors.New("user has no email")
	}

	return user.Email, nil
}

// GetUserEmailByIDSafe returns a user's email with pointer safety check
func GetUserEmailByIDSafe(db *MockDatabase, id int) (string, error) {
	// Get the user
	user, err := db.GetUserByID(id)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// Nil check for user (defensive programming)
	if user == nil {
		return "", errors.New("user is nil")
	}

	// Check if email exists (empty string means no email)
	if user.Email == "" {
		return "", errors.New("user has no email")
	}

	return user.Email, nil
}

// GetUserEmailByIDResult returns a user's email using our Result type
func GetUserEmailByIDResult(db *MockDatabase, id int) Result[string] {
	// Get the user and chain the operations
	return FlatMapTo(
		db.GetUserByIDResult(id),
		func(user UserRecord) Result[string] {
			if user.Email == "" {
				return Err[string](errors.New("user has no email"))
			}
			return Ok(user.Email)
		},
	)
}

// StandardErrorHandlingDemo demonstrates the standard Go error handling approach
func StandardErrorHandlingDemo(t *testing.T) {
	db := NewMockDatabase()

	// Example 1: User exists with email
	email, err := GetUserEmailByID(db, 1)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	fmt.Println("User 1 email:", email)

	// Example 2: User exists but has no email
	email, err = GetUserEmailByID(db, 2)
	if err == nil {
		t.Fatalf("Expected error for missing email")
	}
	fmt.Println("User 2 email error:", err)

	// Example 3: User does not exist
	email, err = GetUserEmailByID(db, 999)
	if err == nil {
		t.Fatalf("Expected error for non-existent user")
	}
	fmt.Println("User 999 error:", err)

	// Potential runtime panic: nil pointer dereference
	// Uncomment to see the panic:
	/*
		db.users[3] = nil // Nil user pointer
		email, err = GetUserEmailByID(db, 3) // This would panic
	*/

	// Safer version with nil check:
	db.users[3] = nil
	email, err = GetUserEmailByIDSafe(db, 3)
	if err == nil {
		t.Fatalf("Expected error for nil user")
	}
	fmt.Println("Nil user error:", err)
}

// JagainErrorHandlingDemo demonstrates error handling with jagain Result type
func JagainErrorHandlingDemo(t *testing.T) {
	db := NewMockDatabase()

	// Example 1: User exists with email
	emailResult := GetUserEmailByIDResult(db, 1)
	if !emailResult.IsOk() {
		t.Fatalf("Expected success, got: %v", emailResult.UnwrapErr())
	}
	fmt.Println("User 1 email:", emailResult.Unwrap())

	// Example 2: User exists but has no email
	emailResult = GetUserEmailByIDResult(db, 2)
	if !emailResult.IsErr() {
		t.Fatalf("Expected error for missing email")
	}
	fmt.Println("User 2 email error:", emailResult.UnwrapErr())

	// Example 3: User does not exist
	emailResult = GetUserEmailByIDResult(db, 999)
	if !emailResult.IsErr() {
		t.Fatalf("Expected error for non-existent user")
	}
	fmt.Println("User 999 error:", emailResult.UnwrapErr())

	// No panic risk with nil values due to proper encapsulation
	// We can use pattern matching for elegant handling:
	message := MatchTo(emailResult,
		func(email string) string { return fmt.Sprintf("Email is: %s", email) },
		func(err error) string { return fmt.Sprintf("Error occurred: %v", err) },
	)
	fmt.Println("Pattern matching result:", message)
}

// This example demonstrates a real-world use case comparing traditional
// nullable SQL fields with Option-based handling
func SQLNullableFieldsExample(t *testing.T) {
	// Traditional Go approach with sql.NullString
	type UserTraditional struct {
		ID    int
		Name  string
		Email sql.NullString // Nullable email
	}

	// Function to process a traditional user
	processTraditionalUser := func(user UserTraditional) (string, error) {
		// Have to check Valid field every time
		if user.Email.Valid {
			return user.Email.String, nil
		}
		return "", errors.New("email not provided")
	}

	// Jagain approach using Option
	type UserWithOption struct {
		ID    int
		Name  string
		Email Option[string] // Optional email
	}

	// Function to process a user with Option
	processOptionUser := func(user UserWithOption) Result[string] {
		// Convert Option to Result
		return user.Email.ToResult(errors.New("email not provided"))
	}

	// Creating test users
	traditionalUser := UserTraditional{
		ID:    1,
		Name:  "Alice",
		Email: sql.NullString{String: "alice@example.com", Valid: true},
	}

	optionUser := UserWithOption{
		ID:    1,
		Name:  "Alice",
		Email: Some("alice@example.com"),
	}

	// Using the traditional approach
	email, err := processTraditionalUser(traditionalUser)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	fmt.Println("Traditional approach email:", email)

	// Using the Option approach
	emailResult := processOptionUser(optionUser)
	if !emailResult.IsOk() {
		t.Fatalf("Expected success, got: %v", emailResult.UnwrapErr())
	}
	fmt.Println("Option approach email:", emailResult.Unwrap())

	// Handling the missing email case
	traditionalUser.Email.Valid = false
	email, err = processTraditionalUser(traditionalUser)
	if err == nil {
		t.Fatalf("Expected error for invalid email")
	}
	fmt.Println("Traditional approach missing email error:", err)

	optionUser.Email = None[string]()
	emailResult = processOptionUser(optionUser)
	if !emailResult.IsErr() {
		t.Fatalf("Expected error for missing email")
	}
	fmt.Println("Option approach missing email error:", emailResult.UnwrapErr())
}

// TestComparisonExamples runs all the comparison examples
func TestComparisonExamples(t *testing.T) {
	t.Run("StandardErrorHandling", func(t *testing.T) {
		StandardErrorHandlingDemo(t)
	})

	t.Run("JagainErrorHandling", func(t *testing.T) {
		JagainErrorHandlingDemo(t)
	})

	t.Run("SQLNullableFields", func(t *testing.T) {
		SQLNullableFieldsExample(t)
	})
}
