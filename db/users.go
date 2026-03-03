package db

import (
	"apartment-app/models"
	"context"
)

// CraeteUser inserts a new user into the database and returns the created user (with their new UUID)
func CreateUser(name string, email string) (models.User, error) {
	var user models.User

	query := `
			INSERT INTO users (name, email)
			VALUES ($1, $2)
			RETURNING id, name, email;
	`

	// QueryRow executes the query and scans the returned row into the struct
	err := Pool.QueryRow(context.Background(), query, name, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)

	return user, err
}
