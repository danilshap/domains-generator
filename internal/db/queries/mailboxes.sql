-- name: CreateMailbox :one
INSERT INTO mailboxes (
    address, 
    password, 
    domain_id, 
    user_id,
    created_at, 
    status
) VALUES (
    $1, $2, $3, 
    $4,
    NOW(), 
    $5
) RETURNING *;

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
WHERE is_deleted = false AND user_id = $1;

-- name: GetMailboxByID :one
SELECT * FROM mailboxes
WHERE id = $1 AND is_deleted = false
LIMIT 1;

-- name: UpdateMailbox :exec
UPDATE mailboxes
SET address = $2, domain_id = $3, user_id = $4
WHERE id = $1;

-- name: GetMailboxesWithFilters :many
SELECT m.*, d.name as domain_name 
FROM mailboxes m
LEFT JOIN domains d ON m.domain_id = d.id
WHERE m.is_deleted = false
  AND CASE WHEN array_length(@status_filter::int[], 1) > 0 THEN m.status = ANY(@status_filter) ELSE true END
  AND CASE WHEN array_length(@domain_filter::int[], 1) > 0 THEN m.domain_id = ANY(@domain_filter) ELSE true END
  AND (@search_query::text = '' OR m.address ILIKE '%' || @search_query || '%')
  AND (@user_id::int IS NULL OR m.user_id = @user_id)
ORDER BY m.created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: GetMailboxesStats :one
SELECT 
    COUNT(*) as total_count,
    COUNT(*) FILTER (WHERE status = 1) as active_count,
    COUNT(*) FILTER (WHERE status = 2) as inactive_count,
    COUNT(DISTINCT domain_id) as domains_count
FROM mailboxes
WHERE is_deleted = false AND user_id = $1;

-- name: UpdateMailboxesStatusByDomainID :exec
UPDATE mailboxes 
SET status = $1
WHERE domain_id = $2 AND is_deleted = false;

-- name: UpdateMailboxesStatusByID :exec
UPDATE mailboxes 
SET status = $1
WHERE id = $2 AND is_deleted = false;

-- name: GetMailboxesByUserID :many
SELECT * FROM mailboxes
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: GetUserByMailboxID :one
SELECT u.* FROM users u
JOIN mailboxes m ON m.user_id = u.id
WHERE m.id = $1 AND u.is_active = true
LIMIT 1;