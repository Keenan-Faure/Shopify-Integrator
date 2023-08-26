-- name: CreateUser :exec
INSERT INTO users (
    id,
    name,
    created_at,
    updated_at,
    api_key
) VALUES (
    ?, ?, ?, ?, ?
);
--