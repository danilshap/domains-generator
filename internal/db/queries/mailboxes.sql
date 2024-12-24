-- name: CreateMailbox :one
INSERT INTO mailboxes (address, password, domain_id, created_at, status)
VALUES ($1, $2, $3, NOW(), $4)
RETURNING id;

-- name: GetMailboxesByDomain :many
SELECT id, address, password, domain_id, status, created_at
FROM mailboxes
WHERE address = $1;

-- name: GetAllMailboxes :many
SELECT id, address, domain_id, status, created_at
FROM mailboxes
ORDER BY created_at DESC;

-- name: GetmailboxesByDomain :many
SELECT id, address, status, created_at
FROM mailboxes
WHERE domain_id = $1
ORDER BY created_at DESC;

-- name: SetMailboxStatus :exec
UPDATE mailboxes
SET status = $1
WHERE id = $2;

-- name: DeleteMailbox :exec
DELETE FROM mailboxes
WHERE id = $1;
