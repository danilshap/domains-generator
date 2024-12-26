-- name: GetMailboxesStats :one
SELECT 
    COUNT(*) FILTER (WHERE status = 1) as active_count,
    COUNT(*) FILTER (WHERE status = 2) as inactive_count,
    COUNT(*) FILTER (WHERE status = 3) as deleted_count,
    COUNT(*) as total_count
FROM mailboxes; 

-- name: GetMailboxesWithFilters :many
SELECT 
    m.id,
    m.address,
    m.status,
    m.domain_id,
    d.name as domain_name
FROM mailboxes m
LEFT JOIN domains d ON m.domain_id = d.id
WHERE 
    ($1::int[] IS NULL OR m.status = ANY($1)) AND
    ($2::int[] IS NULL OR m.domain_id = ANY($2)) AND
    ($3::text IS NULL OR m.address ILIKE '%' || $3 || '%')
ORDER BY m.id DESC
LIMIT $4 OFFSET $5; 

-- name: UpdateMailboxStatus :exec
UPDATE mailboxes
SET status = $2
WHERE id = $1; 