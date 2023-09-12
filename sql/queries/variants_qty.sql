-- name: CreateVariantQty :one
INSERT INTO variant_qty(
    variant_id,
    name,
    value,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5
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