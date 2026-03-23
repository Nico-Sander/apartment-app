package db

import (
	"context"
	"time"

	"apartment-app/models"

	"github.com/google/uuid"
)

// calculateNextDueDate finds the actual calendar date for the next occurrence of a weekday (0=Sun, 6=Sat)
func calculateNextDueDate(targetWeekday int) time.Time {
	today := time.Now()
	currentWeekday := int(today.Weekday())

	// Calculate how many days until the target weekday
	daysUntil := (targetWeekday - currentWeekday + 7) % 7

	// If it's today, schedule it for today. Otherwise, push forward
	return today.AddDate(0, 0, daysUntil)
}

// CreateChore inserts a new chore into the database
func CreateChore(groupID uuid.UUID, title, description string, assignedTo *uuid.UUID, isRecurring bool, intervalUnit string, deadlineWeekday int) (models.Chore, error) {
	var chore models.Chore

	// Calculate the concrete due date
	dueDate := calculateNextDueDate(deadlineWeekday)

	query := `
		INSERT INTO chores (group_id, title, description, assigned_to, is_recurring, interval_unit, deadline_weekday, due_date) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id, group_id, title, description, is_recurring, interval_unit, deadline_weekday, due_date, status;`

	err := Pool.QueryRow(context.Background(), query,
		groupID, title, description, assignedTo, isRecurring, intervalUnit, deadlineWeekday, dueDate,
	).Scan(
		&chore.ID, &chore.GroupID, &chore.Title, &chore.Description,
		&chore.IsRecurring, &chore.IntervalUnit, &chore.DeadlineWeekday, &chore.DueDate, &chore.Status,
	)

	return chore, err
}

// GetChoresByGroup fetches all chores for a specific group, sorted by due date
func GetChoresByGroup(groupID uuid.UUID) ([]models.Chore, error) {
	var chores []models.Chore

	// Order by due_date ASC so the closest deadline are at the top
	query := `
		SELECT id, title, description, is_recurring, due_date, status 
		FROM chores 
		WHERE group_id = $1 AND status != 'completed' 
		ORDER BY due_date ASC;`

	rows, err := Pool.Query(context.Background(), query, groupID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c models.Chore
		// Scan the row into the chore struct
		err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.IsRecurring, &c.DueDate, &c.Status)
		if err != nil {
			return nil, err
		}
		chores = append(chores, c)
	}
	return chores, nil
}

// CompleteChore marks a chore as done, or pushes its due date forward if it is a recurring chore
func CompleteChore(choreID uuid.UUID) error {
	var isRecurring bool

	// 1. Check if the chore is recurring
	err := Pool.QueryRow(context.Background(), "SELECT is_recurring FROM chores WHERE id = $1", choreID).Scan(&isRecurring)
	if err != nil {
		return err
	}

	if isRecurring {
		// 2. If it is recurring, simply push the due date forward by 1 week using Postgres date math
		query := `UPDATE chores SET due_date = (due_date + interval '1 week')::date WHERE id = $1;`
		_, err := Pool.Exec(context.Background(), query, choreID)
		return err
	}

	// 3. If it's a one-time chore, mark it as completed
	query := `UPDATE chores SET status = 'completed' WHERE id = $1;`
	_, err = Pool.Exec(context.Background(), query, choreID)
	return err
}
