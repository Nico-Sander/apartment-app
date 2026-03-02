-- +goose Up
-- +goose StatementBegin
-- 1. Users Table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 2. Groups (Apartments) Table
CREATE TABLE groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    invite_code TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 3. Group Memberships (Linking Users to Groups)
CREATE TABLE group_members (
    user_id UUID REFERENCES users (id) ON DELETE CASCADE,
    group_id UUID REFERENCES groups (id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id, group_id)
);

-- 4. Expenses Table
CREATE TABLE expenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID REFERENCES groups (id) ON DELETE CASCADE NOT NULL,
    paid_by_id UUID REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    description TEXT NOT NULL,
    date TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 5. Chores Table
CREATE TABLE chores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID REFERENCES groups (id) ON DELETE CASCADE NOT NULL,
    assigned_to_id UUID REFERENCES users (id) ON DELETE SET NULL,
    title TEXT NOT NULL,
    is_completed BOOLEAN DEFAULT FALSE,
    due_date TIMESTAMP WITH TIME ZONE
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chores;

DROP TABLE IF EXISTS expenses;

DROP TABLE IF EXISTS group_members;

DROP TABLE IF EXISTS groups;

DROP TABLE IF EXISTS users;

-- +goose StatementEnd
