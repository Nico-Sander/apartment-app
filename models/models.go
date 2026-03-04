package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // The "-" ensures the hash is NEVER sent to the frontend accidentally
}

type Group struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	InviteCode string    `json:"invite_code"`
}

type Expense struct {
	ID          uuid.UUID `json:"id"`
	GroupID     uuid.UUID `json:"group_id"`
	PaitByID    uuid.UUID `json:"paid_by_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

type Chore struct {
	ID          uuid.UUID `json:"id"`
	GroupID     uuid.UUID `json:"group_id"`
	AssignedTo  uuid.UUID `json:"assigned_to_id"`
	Title       string    `json:"title"`
	IsCompleted bool      `json:"is_completed"`
	DueDate     time.Time `json:"due_date"`
}
