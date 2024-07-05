-- name: CreateFeedFollow :one
INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id)
VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteFeedFollowByID :exec
DELETE FROM feed_follows
WHERE id = $1;

-- name: GetFeedFollowByID :one
SELECT * FROM feed_follows
WHERE id = $1
LIMIT 1;

-- name: GetFeedFollowsByUserID :many
SELECT * FROM feed_follows
WHERE user_id = $1;