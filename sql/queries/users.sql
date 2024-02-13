-- name: CreateUser :one
INSERT INTO users (id, name, email, password, created_at, updated_at) 
VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING *;

-- name: GetUserByID :one
SELECT * 
FROM users 
WHERE id = ? AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT * 
FROM users 
WHERE email = ? AND deleted_at IS NULL;

-- name: UpdateUserPassword :exec
UPDATE users 
SET password = ?, updated_at = CURRENT_TIMESTAMP 
WHERE id = ?;

-- name: DeleteUser :exec
UPDATE users 
SET deleted_at = CURRENT_TIMESTAMP 
WHERE id = ?;

-- name: IsUserValidated :one
SELECT CASE 
         WHEN validated_at IS NOT NULL THEN 1 
         ELSE 0 
       END as is_validated
FROM users
WHERE id = ?;

