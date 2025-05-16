package jagain

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
)

// UserRepository represents a hypothetical data access layer
type UserRepository struct {
	users map[int]User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: map[int]User{
			1: {
				ID:    1,
				Name:  "John Doe",
				Email: Some("john@example.com"),
				Age:   Some(30),
			},
			2: {
				ID:    2,
				Name:  "Jane Smith",
				Email: Some("jane@example.com"),
				Age:   None[int](),
			},
		},
	}
}

// FindUser returns a Result with either a User or an error
func (r *UserRepository) FindUser(id int) Result[User] {
	if user, ok := r.users[id]; ok {
		return Ok(user)
	}
	return Err[User](fmt.Errorf("user with ID %d not found", id))
}

// ParseUserID attempts to parse a string ID into an integer
func ParseUserID(id string) Result[int] {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return Err[int](fmt.Errorf("invalid user ID format: %w", err))
	}
	if userID <= 0 {
		return Err[int](errors.New("user ID must be positive"))
	}
	return Ok(userID)
}

// UserService provides business logic for user operations
type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetUserEmailByStringID demonstrates chaining Result operations
func (s *UserService) GetUserEmailByStringID(id string) Result[string] {
	// Use FlatMapTo to chain operations with different types
	idResult := ParseUserID(id)

	// Chain: Result[int] -> Result[User]
	userResult := FlatMapTo(idResult, func(userID int) Result[User] {
		return s.repo.FindUser(userID)
	})

	// Chain: Result[User] -> Result[string]
	return FlatMapTo(userResult, func(user User) Result[string] {
		return user.Email.ToResult(errors.New("user has no email"))
	})
}

func TestResultExample(t *testing.T) {
	repo := NewUserRepository()
	service := NewUserService(repo)

	// Test successful case
	emailResult := service.GetUserEmailByStringID("1")
	if !emailResult.IsOk() {
		t.Errorf("Expected success for valid user ID")
	}
	if emailResult.Unwrap() != "john@example.com" {
		t.Errorf("Expected email to be john@example.com, got %s", emailResult.Unwrap())
	}

	// Test user not found
	nonExistentResult := service.GetUserEmailByStringID("999")
	if !nonExistentResult.IsErr() {
		t.Errorf("Expected error for non-existent user ID")
	}
	if nonExistentResult.UnwrapErr().Error() != "user with ID 999 not found" {
		t.Errorf("Unexpected error: %v", nonExistentResult.UnwrapErr())
	}

	// Test invalid ID format
	invalidIDResult := service.GetUserEmailByStringID("abc")
	if !invalidIDResult.IsErr() {
		t.Errorf("Expected error for invalid user ID format")
	}
	if _, ok := errors.Unwrap(invalidIDResult.UnwrapErr()).(*strconv.NumError); !ok {
		t.Errorf("Expected strconv.NumError, got: %v", invalidIDResult.UnwrapErr())
	}

	// Using pattern matching for better error handling
	message := MatchTo(emailResult,
		func(email string) string {
			return fmt.Sprintf("Email found: %s", email)
		},
		func(err error) string {
			return fmt.Sprintf("Error occurred: %v", err)
		},
	)
	if message != "Email found: john@example.com" {
		t.Errorf("Unexpected message: %s", message)
	}

	// Using error mapping
	errorWithContext := nonExistentResult.MapErr(func(err error) error {
		return fmt.Errorf("while getting user email: %w", err)
	})
	expectedError := "while getting user email: user with ID 999 not found"
	if errorWithContext.UnwrapErr().Error() != expectedError {
		t.Errorf("Expected error to be '%s', got '%s'", expectedError, errorWithContext.UnwrapErr().Error())
	}
}
