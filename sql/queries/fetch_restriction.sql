-- name: CreateFetchRestriction :exec
INSERT INTO push_restriction(
    id,
    field,
    flag,
    updated_at,
    created_at
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: UpdateFetchRestriction :exec
UPDATE push_restriction
SET
    flag = $1,
    updated_at = $2
WHERE field = $3;
