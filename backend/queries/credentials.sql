-- name: CreateCredential :one
INSERT INTO credentials (
    id, user_id, name, provider, region, endpoint,
    encrypted_access_key, encrypted_secret_key,
    use_ssl, status, logo
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: ListCredentials :many
SELECT * FROM credentials
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetCredential :one
SELECT * FROM credentials
WHERE id = $1 AND user_id = $2;

-- name: UpdateCredential :exec
UPDATE credentials
SET name = $3, provider = $4, region = $5, endpoint = $6,
    encrypted_access_key = $7, encrypted_secret_key = $8,
    use_ssl = $9, status = $10, logo = $11, updated_at = NOW()
WHERE id = $1 AND user_id = $2;

-- name: DeleteCredential :exec
DELETE FROM credentials WHERE id = $1 AND user_id = $2;
