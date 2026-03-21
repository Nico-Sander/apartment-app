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
