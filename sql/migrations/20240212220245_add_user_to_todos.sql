-- +goose Up
-- +goose StatementBegin

-- Seed a user
INSERT INTO users (id, name, email, password) VALUES ('seeded_user_id', 'Seeded User', 'seededuser@example.com', 'hashed_password');

-- Create a new table with the additional column and foreign key
CREATE TABLE new_todo_items (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    description TEXT NOT NULL,
    complete BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Copy data from the old table to the new table, setting user_id to the seeded user's ID
INSERT INTO new_todo_items (id, user_id, description, complete, created_at, completed_at)
    SELECT id, 'seeded_user_id', description, complete, created_at, completed_at FROM todo_items;

-- Drop the old table
DROP TABLE todo_items;

-- Rename the new table to the original name
ALTER TABLE new_todo_items RENAME TO todo_items;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- To reverse the migration, follow the reverse process
CREATE TABLE old_todo_items (
    id TEXT PRIMARY KEY,
    description TEXT NOT NULL,
    complete BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL
);

INSERT INTO old_todo_items (id, description, complete, created_at, completed_at)
    SELECT id, description, complete, created_at, completed_at FROM todo_items;

DROP TABLE todo_items;

ALTER TABLE old_todo_items RENAME TO todo_items;

-- Optionally, remove the seeded user
-- DELETE FROM users WHERE id = 'seeded_user_id';

-- +goose StatementEnd

