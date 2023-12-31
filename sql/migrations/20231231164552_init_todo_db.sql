-- +goose Up
-- +goose StatementBegin
CREATE TABLE todo_items (
    id TEXT PRIMARY KEY,
    description TEXT NOT NULL,
    complete BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABOLE todo_items;
-- +goose StatementEnd
