-- name: CreateUser :execresult
INSERT INTO users (
    id,
    name,
    created_at,
    updated_at,
    api_key
) VALUES (
    ?, ?, ?, ?, ?
);

-- name: GetUserByApiKey :one
SELECT * FROM users
WHERE api_key = ?
LIMIT 1;

-- name: GetUserByName :one
SELECT * FROM users
WHERE name = ?
LIMIT 1;