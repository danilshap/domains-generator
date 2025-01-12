-- name: GetNotifications :many
SELECT * FROM notifications 
WHERE user_id = $1 
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetNotificationsCount :one
SELECT COUNT(*) FROM notifications 
WHERE user_id = $1;

-- name: GetUnreadNotificationsCount :one
SELECT COUNT(*) FROM notifications 
WHERE user_id = $1 AND read_at IS NULL;

-- name: MarkNotificationRead :exec
UPDATE notifications 
SET read_at = NOW() 
WHERE id = $1;

-- name: MarkAllNotificationsRead :exec
UPDATE notifications 
SET read_at = NOW() 
WHERE user_id = $1 AND read_at IS NULL;

-- name: CreateNotification :one
INSERT INTO notifications (
    user_id, title, message, type
) VALUES (
    $1, $2, $3, $4
) RETURNING *; 