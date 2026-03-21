-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS chores;

-- ADD THIS LINE TO CLEAR THE GHOST TABLE
CREATE TABLE chores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL REFERENCES groups (id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    -- "AssignedTo" can be null in case a chore is created but not assigned
    assigned_to UUID REFERENCES users (id) ON DELETE SET NULL,
    -- Recurring Logic
    is_recurring BOOLEAN NOT NULL DEFAULT FALSE,
    interval_unit TEXT, -- e.g., 'day', 'week', 'month', 'year'
    deadline_weekday INT, -- 0 = Sunday, 1 = Monday, ...
    -- Concrete date to sort chores by "upcoming"
    due_date DATE,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE chores
-- +goose StatementEnd
