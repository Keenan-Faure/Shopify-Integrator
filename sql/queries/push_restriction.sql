-- name: CreatePushReestriction :exec
INSERT INTO push_restriction(
    id,
    field,
    flag,
    updated_at,
    created_at
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: UpdatePushRestriction :exec
UPDATE push_restriction
SET
    flag = $1,
    updated_at = $2
WHERE field = $3;

-- name: GetPushRestriction :many
SELECT * FROM push_restriction;