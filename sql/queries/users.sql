-- name: CreateUser :execresult
INSERT INTO users (
    id,
    name,
    webhook_token,
    created_at,
    updated_at,
    api_key
) VALUES (
    ?, ?, ?, ?, ?, ?
);

-- name: UpdateUser :execresult
UPDATE users 
SET
    name = ?,
    updated_at = ?
WHERE id = ?;

-- name: GetUserByApiKey :one
SELECT * FROM users
WHERE api_key = ?
LIMIT 1;

-- name: GetUserByName :one
SELECT * FROM users
WHERE name = ?
LIMIT 1;

-- name: ValidateWebhookByUser :one
SELECT
    name
FROM users
WHERE 
webhook_token = ? AND name = ?;