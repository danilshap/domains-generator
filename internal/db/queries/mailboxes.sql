-- name: CreateMailbox :one
INSERT INTO mailboxes (address, password, domain_id, created_at, status)
VALUES ($1, $2, $3, NOW(), $4)
RETURNING *;

-- name: GetMailboxesByDomain :many
SELECT id, address, password, domain_id, status, created_at
FROM mailboxes
WHERE address = $1 AND is_deleted = false;

-- name: GetAllMailboxes :many
SELECT m.*, d.name as domain_name 
FROM mailboxes m
LEFT JOIN domains d ON m.domain_id = d.id
WHERE m.is_deleted = false
ORDER BY m.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetMailboxesByDomainName :many
SELECT id, address, status, created_at
FROM mailboxes
WHERE domain_id = $1 AND is_deleted = false
ORDER BY created_at DESC;

-- name: SetMailboxStatus :exec
UPDATE mailboxes
SET status = $1
WHERE id = $2;

-- name: DeleteMailbox :exec
UPDATE mailboxes 
SET is_deleted = true 
WHERE id = $1;

-- name: GetMailboxesByDomainID :many
SELECT * FROM mailboxes 
WHERE domain_id = $1 AND is_deleted = false
LIMIT $2 OFFSET $3;

-- name: GetMailboxesCountByDomainID :one
SELECT COUNT(*) FROM mailboxes 
WHERE domain_id = $1;

-- name: GetMailboxCountByDomainID :one
SELECT COUNT(*) FROM mailboxes 
WHERE domain_id = $1 AND is_deleted = false;

-- name: GetMailboxesCount :one
SELECT COUNT(*) FROM mailboxes
WHERE is_deleted = false;

-- name: GetMailboxByID :one
SELECT * FROM mailboxes
WHERE id = $1 AND is_deleted = false
LIMIT 1;

-- name: UpdateMailbox :exec
UPDATE mailboxes
SET address = $2, domain_id = $3
WHERE id = $1;

-- name: GetMailboxesWithFilters :many
SELECT m.*, d.name as domain_name 
FROM mailboxes m
LEFT JOIN domains d ON m.domain_id = d.id
WHERE m.is_deleted = false
  AND CASE WHEN array_length($1::int[] /* status_filter */, 1) > 0 THEN m.status = ANY($1) ELSE true END
  AND CASE WHEN array_length($2::int[] /* domain_filter */, 1) > 0 THEN m.domain_id = ANY($2) ELSE true END
  AND ($3 /* search */ = '' OR m.address ILIKE '%' || $3 || '%')
ORDER BY m.created_at DESC
LIMIT $4 /* limit */ OFFSET $5 /* offset */;

-- name: GetMailboxesStats :one
SELECT 
    COUNT(*) as total_count,
    COUNT(*) FILTER (WHERE status = 1) as active_count,
    COUNT(*) FILTER (WHERE status = 2) as inactive_count,
    COUNT(DISTINCT domain_id) as domains_count
FROM mailboxes
WHERE is_deleted = false;

-- name: UpdateMailboxesStatusByDomainID :exec
UPDATE mailboxes 
SET status = $1
WHERE domain_id = $2 AND is_deleted = false;
