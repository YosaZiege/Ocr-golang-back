-- name: CreateUser :one
INSERT INTO users (username, email, password_hash, provider)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2
WHERE username = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE username = $1;

