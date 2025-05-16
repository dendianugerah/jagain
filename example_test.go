package jagain

import (
	"encoding/json"
	"fmt"
	"testing"
)

// User represents a user in a system
type User struct {
	ID        int
	Name      string
	Email     Option[string] // Email is optional
	Age       Option[int]    // Age is optional
	Addresses []Address
}

// Address represents a user's address
type Address struct {
	Street      string
	City        string
	State       string
	ZipCode     string
	CountryCode Option[string] // CountryCode is optional
}

// UserDTO is a data transfer object for User
type UserDTO struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     *string   `json:"email"` // Traditional nullable field
	Age       *int      `json:"age"`   // Traditional nullable field
	Addresses []Address `json:"addresses"`
}

func (u User) ToDTO() UserDTO {
	return UserDTO{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email.ToPtr(),
		Age:       u.Age.ToPtr(),
		Addresses: u.Addresses,
	}
}

func UserFromDTO(dto UserDTO) User {
	return User{
		ID:        dto.ID,
		Name:      dto.Name,
		Email:     FromPtr(dto.Email),
		Age:       FromPtr(dto.Age),
		Addresses: dto.Addresses,
	}
}

// ExampleEmail demonstrates using Option for email validation
func ExampleEmail() {
	// Valid email
	email := "user@example.com"
	validEmailOpt := Some(email)

	// Invalid email handled with None
	invalidEmail := ""
	var invalidEmailOpt Option[string]
	if isValidEmail(invalidEmail) {
		invalidEmailOpt = Some(invalidEmail)
	} else {
		invalidEmailOpt = None[string]()
	}

	// Use pattern matching to handle both cases
	printEmail := func(emailOpt Option[string]) {
		emailOpt.Match(
			func(e string) string {
				fmt.Printf("Email: %s\n", e)
				return e
			},
			func() string {
				fmt.Println("No email provided")
				return ""
			},
		)
	}

	printEmail(validEmailOpt)   // Prints: Email: user@example.com
	printEmail(invalidEmailOpt) // Prints: No email provided

	// Output:
	// Email: user@example.com
	// No email provided
}

// Simple email validation for demonstration
func isValidEmail(email string) bool {
	return email != ""
}

func TestUserSerialization(t *testing.T) {
	// Create a user with some fields present and some not
	user := User{
		ID:    1,
		Name:  "John Doe",
		Email: Some("john@example.com"),
		Age:   None[int](), // Age not provided
		Addresses: []Address{
			{
				Street:      "123 Main St",
				City:        "Anytown",
				State:       "CA",
				ZipCode:     "12345",
				CountryCode: Some("US"),
			},
			{
				Street:      "456 High St",
				City:        "Othertown",
				State:       "NY",
				ZipCode:     "67890",
				CountryCode: None[string](), // CountryCode not provided
			},
		},
	}

	// Convert to DTO and marshal to JSON
	dto := user.ToDTO()
	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}

	// Parse back from JSON
	var parsedDTO UserDTO
	err = json.Unmarshal(jsonBytes, &parsedDTO)
	if err != nil {
		t.Fatalf("Failed to unmarshal user: %v", err)
	}

	// Convert back to domain model
	parsedUser := UserFromDTO(parsedDTO)

	// Validate
	if parsedUser.ID != user.ID {
		t.Errorf("ID mismatch: expected %d, got %d", user.ID, parsedUser.ID)
	}
	if parsedUser.Name != user.Name {
		t.Errorf("Name mismatch: expected %s, got %s", user.Name, parsedUser.Name)
	}
	if parsedUser.Email.Unwrap() != user.Email.Unwrap() {
		t.Errorf("Email mismatch: expected %s, got %s", user.Email.Unwrap(), parsedUser.Email.Unwrap())
	}
	if !parsedUser.Age.IsNone() {
		t.Errorf("Expected Age to be None")
	}

	// Using UnwrapOr for default values
	age := parsedUser.Age.UnwrapOr(0)
	if age != 0 {
		t.Errorf("Expected default age to be 0, got %d", age)
	}

	// Check the addresses
	if len(parsedUser.Addresses) != len(user.Addresses) {
		t.Fatalf("Address count mismatch: expected %d, got %d", len(user.Addresses), len(parsedUser.Addresses))
	}

	firstAddress := parsedUser.Addresses[0]
	if firstAddress.CountryCode.Unwrap() != "US" {
		t.Errorf("First address country code mismatch")
	}

	secondAddress := parsedUser.Addresses[1]
	if !secondAddress.CountryCode.IsNone() {
		t.Errorf("Second address should have no country code")
	}
}
