-- name: CreateDomain :one
INSERT INTO domains (name, provider, status, created_at, expires_at)
VALUES ($1, $2, $3, NOW(), $4)
RETURNING *;

-- name: GetDomainByID :one
SELECT id, name, provider, status, created_at, expires_at, is_deleted, settings FROM domains
WHERE id = $1 AND is_deleted = false LIMIT 1;

-- name: GetAllDomains :many
SELECT d.*, COUNT(m.id) as mailboxes_count 
FROM domains d
LEFT JOIN mailboxes m ON d.id = m.domain_id AND m.is_deleted = false
WHERE d.is_deleted = false AND d.user_id = $3
GROUP BY d.id
ORDER BY d.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetDomainsCount :one
SELECT COUNT(*) FROM domains
WHERE is_deleted = false AND user_id = $1;

-- name: GetDomainByName :one
SELECT id, name, provider, status, created_at, expires_at, is_deleted
FROM domains
WHERE name = $1 AND is_deleted = false;

-- name: UpdateDomain :exec
UPDATE domains 
SET name = $2, provider = $3
WHERE id = $1;

-- name: SetDomainStatus :exec
UPDATE domains 
SET status = $2
WHERE id = $1;

-- name: DeleteDomain :exec
UPDATE domains 
SET is_deleted = true 
WHERE id = $1;

-- name: GetDomainsByUserID :many
SELECT * FROM domains
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: GetUserByDomainID :one
SELECT u.* FROM users u
JOIN domains d ON d.user_id = u.id
WHERE d.id = $1 AND u.is_active = true
LIMIT 1;
