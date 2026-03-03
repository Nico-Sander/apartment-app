package db

import (
	"log"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

// setupTestDB ensures the connection to the database is successful before running any other tests
func setupTestDB() {
	// Test run in the directory of the package, so first the .env file needs to located in the root
	envPath := filepath.Join("..", ".env")

	err := godotenv.Overload(envPath)

	if err != nil {
		log.Printf("Warning: Could not explicitly load .env file from %s: %v", envPath, err)
	} else {
		log.Println("Successfully forced load of .env file for tests")
	}

	// Only initialize if the pools isn't already set up
	if Pool == nil {
		InitDB()
	}
}

func TestCreateUser(t *testing.T) {
	// 1. Setup and clear the database
	setupTestDB()
	ResetDatabase()

	// 2. Execute the function we want to test
	name := "Alex"
	email := "alex@example.com"
	user, err := CraeteUser(name, email)

	// 3. Assertions (Check if it worked)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Name != name {
		t.Errorf("Expected name %s, got %s", name, user.Name)
	}

	if user.Email != email {
		t.Errorf("Expected email %s, got %s", email, user.Email)
	}

	if user.ID.String() == "" {
		t.Errorf("Expeced a valid UUID, got an empty string")
	}
}
