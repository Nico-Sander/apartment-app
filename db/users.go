package db

import (
	"apartment-app/auth"
	"apartment-app/models"
	"context"
)

// CraeteUser inserts a new user into the database and returns the created user (with their new UUID)
func CreateUser(name, email, password string) (models.User, string, error) {
	var user models.User

	// 1. Hash the incoming password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return user, "", err
	}

	// 2. Insert into the DB
	query := `
			INSERT INTO users (name, email, password_hash)
			VALUES ($1, $2, $3)
			RETURNING id, name, email;
	`

	// QueryRow executes the query and scans the returned row into the struct
	err = Pool.QueryRow(context.Background(), query, name, email, hashedPassword).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)
	if err != nil {
		return user, "", err
	}

	// 3. Create the JWT Token for the new user
	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return user, "", err
	}

	// Return the user struct, the token string, and no error
	return user, token, err
}

// GetUserByEmail fetches a user by their email, including their pasoword hash
func GetUserByEmail(email string) (models.User, error) {
	var user models.User

	query := `
	SELECT id, name, email, password_hash
	FROM users
	WHERE email = $1;`

	err := Pool.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
	)

	return user, err
}
