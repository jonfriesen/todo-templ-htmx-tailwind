-- name: ListTodos :many
SELECT * FROM todo_items WHERE user_id = ?;

-- name: InsertTodo :exec
INSERT INTO todo_items (id, user_id, description, complete)
VALUES (?, ?, ?, ?);

-- name: GetTodo :one
SELECT * FROM todo_items WHERE id = ?;

-- name: CompleteTodo :exec
UPDATE todo_items
SET
    complete = CASE WHEN complete = 0 THEN 1 ELSE 0 END,
    completed_at = CASE WHEN complete = 0 THEN CURRENT_TIMESTAMP ELSE NULL END
WHERE id = ?;
