-- name: CreateDomain :one
INSERT INTO domains (name, provider, status, created_at, expires_at)
VALUES ($1, $2, $3, NOW(), $4)
RETURNING *;

-- name: GetDomainByID :one
SELECT * FROM domains
WHERE id = $1 LIMIT 1;

-- name: GetAllDomains :many
SELECT * FROM domains
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetDomainsCount :one
SELECT COUNT(*) FROM domains;

-- name: GetDomainByName :one
SELECT id, name, provider, status, created_at, expires_at, is_deleted
FROM domains
WHERE name = $1;

-- name: SetDomainStatus :exec
UPDATE domains
SET status = $1
WHERE id = $2;

-- name: DeleteDomain :exec
DELETE FROM domains
WHERE id = $1;
