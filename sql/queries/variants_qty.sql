-- name: CreateVariantQty :one
INSERT INTO variant_qty(
    id,
    variant_id,
    name,
    isdefault,
    value,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateVariantQty :exec
UPDATE variant_qty
SET
    name = $1,
    value = $2,
    isdefault = $3
WHERE variant_id IN (
    SELECT id FROM variants
    WHERE sku = $4
) AND name = $5;

-- name: GetVariantQty :many
SELECT 
    name,
    value,
    isdefault
FROM variant_qty
WHERE variant_id = $1;

-- name: GetVariantQtyBySKU :many
SELECT
    name,
    value,
    isdefault
FROM variant_qty
WHERE variant_id IN (
    SELECT id FROM variants
    WHERE sku = $1
) AND name = $2;

-- name: GetCountOfUniqueWarehouses :one
SELECT COUNT(DISTINCT "name") FROM variant_qty;

-- name: GetUniqueWarehouses :many
SELECT DISTINCT "name" FROM variant_qty;

-- name: RemoveQty :exec
DELETE FROM variant_qty
WHERE id = $1;

