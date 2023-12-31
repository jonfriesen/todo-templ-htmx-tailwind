-- +goose Up
-- +goose StatementBegin
INSERT INTO todo_items (id, description, complete) 
VALUES 
    ('seed-id1', 'Walk the dog', 0),
    ('seed-id2', 'Pay taxes', 1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM todoitems
WHERE id IN ('seed-id1', 'seed-id2');
-- +goose StatementEnd
