-- name: CreateUser :one
INSERT INTO users (
    id,
    "name",
    email,
    "password",
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUsers :one
SELECT * FROM users LIMIT 1;

-- name: GetUserByName :one
SELECT
    "name"
FROM users
WHERE "name" = $1
LIMIT 1;

-- name: GetUserCredentials :one
SELECT
    "name",
    "password",
    api_key
FROM users
WHERE "name" = $1
AND "password" = $2
LIMIT 1;

-- name: UpdateUser :execresult
UPDATE users 
SET
    "name" = $1,
    email = $2,
    updated_at = $3
WHERE id = $4;

-- name: GetUserByApiKey :one
SELECT * FROM users
WHERE api_key = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: ValidateWebhookByUser :one
SELECT
    "name"
FROM users
WHERE 
webhook_token = $1 AND api_key = $2;

-- name: RemoveUser :exec
DELETE FROM users
WHERE api_key = $1;