package db

import (
	"testing"
)

func TestCreateGroupAndJoin(t *testing.T) {
	// 1. Setup and clear the database
	setupTestDB()
	ResetDatabase()

	// 2. We need a user to act as the creator, so let's create one first!
	user, _, err := CreateUser("Test Creator", "creator@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to setup test user: %v", err)
	}

	// 3. Execute our new group function
	groupName := "The Cool Apartment"
	group, err := CreateGroupAndJoin(groupName, user.ID)

	// 4. Check the results
	if err != nil {
		t.Fatalf("Expected no error creating group, got: %v", err)
	}

	if group.Name != groupName {
		t.Errorf("Expected group name %s, got %s", groupName, group.Name)
	}

	if len(group.InviteCode) != 6 {
		t.Errorf("Expected invite code to be 6 characters, got %d (%s)", len(group.InviteCode), group.InviteCode)
	}

	if group.ID.String() == "" {
		t.Errorf("Expected valid group UUID, got empty string")
	}
}
