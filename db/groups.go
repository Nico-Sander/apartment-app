package db

import (
	"context"
	"crypto/rand"
	"errors"
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

// JoinGroupByCode looks up a group by its invite code and adds the user as a member
func JoinGroupByCode(inviteCode string, userID uuid.UUID) (models.Group, error) {
	var group models.Group
	ctx := context.Background()

	// 1. Find the group using the invite code
	groupQuery := `
		SELECT id, name, invite_code
		FROM groups
		WHERE invite_code = $1;`

	err := Pool.QueryRow(ctx, groupQuery, inviteCode).Scan(
		&group.ID,
		&group.Name,
		&group.InviteCode,
	)
	if err != nil {
		return group, errors.New("invalid invite code")
	}

	// 2. Link the user to the group
	memberQuery := `
		INSERT INTO group_members (group_id, user_id)
		VALUES ($1, $2);`

	_, err = Pool.Exec(ctx, memberQuery, group.ID, userID)
	if err != nil {
		return group, errors.New("you are already a member of this group")
	}

	return group, nil
}

// GetUserGroups fetches all groups a specific user is a member of
func GetUserGroups(userID uuid.UUID) ([]models.Group, error) {
	var groups []models.Group

	query := `
		SELECT g.id, g.name, g.invite_code 
		FROM groups g
		JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.user_id = $1;`

	rows, err := Pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.ID, &group.Name, &group.InviteCode); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}
