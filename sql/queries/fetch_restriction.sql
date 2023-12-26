-- name: CreateFetchRestriction :exec
INSERT INTO fetch_restriction(
    id,
    field,
    flag,
    updated_at,
    created_at
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: UpdateFetchRestriction :exec
UPDATE fetch_restriction
SET
    flag = $1,
    updated_at = $2
WHERE field = $3;

-- name: GetFetchRestriction :many
SELECT * FROM fetch_restriction;
