-- name: CreateUser :execresult
INSERT INTO users (
    name,
    email,
    created_at,
    updated_at
) VALUES (
    ?, ?, ?, ?
);

-- name: GetUsers :one
SELECT * FROM users LIMIT 1;

-- name: UpdateUser :execresult
UPDATE users 
SET
    name = ?,
    email = ?,
    updated_at = ?
WHERE id = ?;

-- name: GetUserByApiKey :one
SELECT * FROM users
WHERE api_key = ?
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?;

-- name: ValidateWebhookByUser :one
SELECT
    name
FROM users
WHERE 
webhook_token = ? AND name = ?;