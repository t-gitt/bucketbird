-- name: InsertBucket :one
INSERT INTO buckets (id, user_id, credential_id, name, region, description)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListBuckets :many
SELECT
    sqlc.embed(b),
    c.name as credential_name,
    c.provider as credential_provider
FROM buckets b
JOIN credentials c ON c.id = b.credential_id
WHERE b.user_id = $1
ORDER BY b.created_at DESC;

-- name: GetBucket :one
SELECT
    sqlc.embed(b),
    c.name as credential_name,
    c.provider as credential_provider
FROM buckets b
JOIN credentials c ON c.id = b.credential_id
WHERE b.id = $1 AND b.user_id = $2;

-- name: GetBucketByName :one
SELECT
    sqlc.embed(b),
    c.name as credential_name,
    c.provider as credential_provider
FROM buckets b
JOIN credentials c ON c.id = b.credential_id
WHERE b.user_id = $1 AND b.name = $2;

-- name: UpdateBucketSize :exec
UPDATE buckets
SET size_bytes = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateBucket :exec
UPDATE buckets
SET description = $3, updated_at = NOW()
WHERE id = $1 AND user_id = $2;

-- name: DeleteBucket :exec
DELETE FROM buckets WHERE id = $1 AND user_id = $2;
