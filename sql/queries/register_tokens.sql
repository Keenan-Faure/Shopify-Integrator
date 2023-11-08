-- name: CreateToken :one
INSERT INTO register_tokens(
    id,
    name,
    email,
    token,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateToken :one
UPDATE register_tokens
SET
   token = $1
where email = $2
RETURNING *;

-- name: DeleteToken :exec
DELETE FROM register_tokens
WHERE
token = $1 AND
email = $2;

-- name: GetToken :one
SELECT
    name,
    email,
    token
FROM register_tokens
WHERE name = $1
AND email = $2;

-- name: GetTokenValidation :one
SELECT
    name,
    email,
    token
FROM register_tokens
WHERE name = $1
AND email = $2;
