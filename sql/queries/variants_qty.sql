-- name: CreateVariantQty :one
INSERT INTO variant_qty(
    id,
    variant_id,
    name,
    value,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateVariantQty :one
UPDATE variant_qty
SET
    name = $1,
    value = $2
WHERE variant_id = $3
RETURNING *;

-- name: GetVariantQty :many
SELECT 
    name,
    value
FROM variant_qty
WHERE variant_id = $1;

-- name: RemoveQty :exec
DELETE FROM variant_qty
WHERE id = $1;