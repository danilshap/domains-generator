-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    hashed_password,
    full_name,
    role
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND is_active = true LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND is_active = true LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 AND is_active = true LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
WHERE is_active = true
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET
    username = COALESCE(sqlc.narg(username), username),
    email = COALESCE(sqlc.narg(email), email),
    full_name = COALESCE(sqlc.narg(full_name), full_name),
    role = COALESCE(sqlc.narg(role), role),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET
    hashed_password = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: DeactivateUser :exec
UPDATE users
SET
    is_active = false,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetUserStats :one
SELECT
    COUNT(DISTINCT d.id) as domains_count,
    COUNT(DISTINCT m.id) as mailboxes_count,
    COUNT(DISTINCT CASE WHEN m.status = 1 THEN m.id END) as active_mailboxes_count,
    COUNT(DISTINCT CASE WHEN m.status = 2 THEN m.id END) as inactive_mailboxes_count
FROM users u
LEFT JOIN domains d ON d.user_id = u.id
LEFT JOIN mailboxes m ON m.user_id = u.id
WHERE u.id = $1;

-- name: VerifyUserCredentials :one
SELECT * FROM users
WHERE email = $1 AND is_active = true
LIMIT 1;