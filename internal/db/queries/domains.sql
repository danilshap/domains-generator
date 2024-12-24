-- name: CreateDomain :one
INSERT INTO domains (name, provider, status, created_at, expires_at)
VALUES ($1, $2, $3, NOW(), $4)
RETURNING id;

-- name: GetDomainByName :one
SELECT id, name, provider, status, created_at, expires_at
FROM domains
WHERE name = $1;

-- name: GetAllDomains :many
SELECT id, name, provider, status, created_at, expires_at
FROM domains
ORDER BY created_at DESC;

-- name: SetDomainStatus :exec
UPDATE domains
SET status = $1
WHERE id = $2;

-- name: DeleteDomain :exec
DELETE FROM domains
WHERE id = $1;
