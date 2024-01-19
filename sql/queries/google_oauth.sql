-- name: CreateOAuthRecord :exec
INSERT INTO google_oauth(
    id,
    google_id,
    email,
    picture,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
);