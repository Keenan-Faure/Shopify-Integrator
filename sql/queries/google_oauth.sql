-- name: CreateOAuthRecord :one
INSERT INTO google_oauth(
    id,
    user_id,
    google_id,
    email,
    picture,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetUserByGoogleID :one
SELECT * FROM google_oauth
WHERE google_id = $1;

