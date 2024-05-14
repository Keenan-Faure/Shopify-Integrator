-- name: AddRunTimeFlag :exec
INSERT INTO runtime_flags(
    id,
    flag_name,
    flag_value,
    updated_at,
    created_at
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: UpsertRunTimeFlag :exec
INSERT INTO runtime_flags(
    id,
    flag_name,
    flag_value,
    updated_at,
    created_at
) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT(flag_name)
DO UPDATE
SET
    flag_value = $3,
    updated_at = $4
;

-- name: GetRuntimeFlag :one
SELECT
    flag_name,
    flag_value,
    updated_at
FROM runtime_flags
WHERE flag_name = $1;

-- name: RemoveRuntimeFlag :exec
DELETE FROM runtime_flags WHERE flag_name = $1;