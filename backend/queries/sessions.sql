-- name: CreateSession :one
INSERT INTO sessions (id, user_id, refresh_token_hash, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetSessionByHash :one
SELECT * FROM sessions WHERE refresh_token_hash = $1;

-- name: UpdateSessionToken :exec
UPDATE sessions
SET refresh_token_hash = $2, expires_at = $3, updated_at = NOW()
WHERE id = $1;

-- name: DeleteSessionByHash :exec
DELETE FROM sessions WHERE refresh_token_hash = $1;

-- name: DeleteSessionsForUser :exec
DELETE FROM sessions WHERE user_id = $1;
