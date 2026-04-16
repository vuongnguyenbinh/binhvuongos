-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT * FROM users WHERE deleted_at IS NULL ORDER BY full_name LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users WHERE deleted_at IS NULL;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, full_name, role, phone, status)
VALUES ($1, $2, $3, $4, $5, 'active') RETURNING *;

-- name: UpdateUser :one
UPDATE users SET full_name = $2, role = $3, phone = $4, status = $5
WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: UpdateLastLogin :exec
UPDATE users SET last_login_at = NOW() WHERE id = $1;

-- name: SoftDeleteUser :exec
UPDATE users SET deleted_at = NOW() WHERE id = $1;
