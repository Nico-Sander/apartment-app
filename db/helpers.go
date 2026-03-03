package db

import (
	"context"
	"log"
)

// ResetDatabase wipes all data from all tables.
// USE WITH CAUTION: Only run this locally in tests!
func ResetDatabase() {
	query := `TRUNCATE TABLE chores, expenses, group_members, groups, users CASCADE`

	_, err := Pool.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to reset database: %v\n", err)
	}
	log.Println("Database completely cleared")
}
