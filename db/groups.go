package db

import (
	"context"
	"crypto/rand"
	"math/big"

	"apartment-app/models"

	"github.com/google/uuid"
)

// generateInviteCode creates a secure, random 6-character string (A-Z, 0.9)
func generateInviteCode() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var code []byte
	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code = append(code, charset[n.Int64()])
	}
	return string(code), nil
}

// CreateGroupAndJoin creates a new apartment group and immediately adds the creator.
func CreateGroupAndJoin(groupName string, creatorID uuid.UUID) (models.Group, error) {
	var group models.Group
	ctx := context.Background()

	// 1. Start the SQL Transaction
	tx, err := Pool.Begin(ctx)
	if err != nil {
		return group, err
	}
	defer tx.Rollback(ctx)

	// 2. The Retry Loop (Try up to 3 times to get a unique code)
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		inviteCode, err := generateInviteCode()
		if err != nil {
			return group, err
		}

		// Try to insert the group
		groupQuery := `
			INSERT INTO groups (name, invite_code) 
			VALUES ($1, $2) 
			RETURNING id, name, invite_code;`

		err = tx.QueryRow(ctx, groupQuery, groupName, inviteCode).Scan(
			&group.ID,
			&group.Name,
			&group.InviteCode,
		)

		// If there is no error, the code was unique! Break out of the retry loop.
		if err == nil {
			break
		}

		// If we hit our max retries and still have an error, fail completely
		if i == maxRetries-1 {
			return group, err
		}
		// Otherwise, the loop restarts, generates a new code, and tries again!
	}

	// 3. Link the creator to the group
	memberQuery := `
		INSERT INTO group_members (group_id, user_id) 
		VALUES ($1, $2);`

	_, err = tx.Exec(ctx, memberQuery, group.ID, creatorID)
	if err != nil {
		return group, err
	}

	// 4. Commit the transaction
	err = tx.Commit(ctx)
	return group, err
}
