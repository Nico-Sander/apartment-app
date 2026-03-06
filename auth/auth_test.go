package auth

import (
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestPasswordHashing(t *testing.T) {
	password := "my_secure_apartment_password!123"

	// 1. Test Hashing
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Expected no error when hashing, got: %v", err)
	}

	if hash == password {
		t.Fatal("Security Alert: Hash matches the plain text password")
	}

	if hash == "" {
		t.Fatal("Generated hash is empty")
	}

	// 2. Test Correct Password Verification
	isValid := CheckPasswordHash(password, hash)
	if !isValid {
		t.Errorf("Expected true for correct password, got false")
	}

	// 3. Test Incorrect Password Verification
	isInvalid := CheckPasswordHash("wrong_password", hash)
	if isInvalid {
		t.Errorf("Security Alert: CheckPasswordHash returned true for the WRONG password")
	}
}

func TestJWTGenerationValidation(t *testing.T) {
	// Setup a temporary environment variable just for this test
	// This way the test doesn't rely on the actual .env file
	os.Setenv("JWT_SECRET", "temporary_test_secret")
	defer os.Unsetenv("JWT_SECRET")

	// Generate a fake user ID
	originalUserID := uuid.New()

	// 1. Test Token Generation
	tokenString, err := GenerateJWT(originalUserID)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	if tokenString == "" {
		t.Fatalf("Generated JWT is empty")
	}

	// 2. Test Token Validation
	parsedUserID, err := ValidateJWT(tokenString)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	// 3. Ensure the IDs match
	if parsedUserID != originalUserID {
		t.Errorf("Expected userID %s, got %s", originalUserID, parsedUserID)
	}

	// 4. Test an invalid token
	_, err = ValidateJWT("this.is.not.a.valid.token")
	if err == nil {
		t.Errorf("Security Alert: ValidateJWT accepted garbage data as a valid token")
	}
}
