-- name: CreateVariantQty :one
INSERT INTO variant_qty(
    id,
    variant_id,
    "name",
    "value",
    isdefault,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateVariantQty :exec
UPDATE variant_qty
SET
    "name" = COALESCE($1, "name"),
    "value" = COALESCE($2, "value"),
    isdefault = COALESCE($3, isdefault),
    updated_at = $4
WHERE variant_id IN (
    SELECT id FROM variants
    WHERE sku = $5
) AND "name" = $6;

-- name: GetVariantQty :many
SELECT 
    "name",
    "value",
    isdefault,
    updated_at
FROM variant_qty
WHERE variant_id = $1;

-- name: GetVariantQtyBySKU :many
SELECT
    "name",
    "value",
    isdefault,
    updated_at
FROM variant_qty
WHERE variant_id IN (
    SELECT id FROM variants
    WHERE sku = $1
) AND "name" = $2;

-- name: GetCountOfUniqueWarehouses :one
SELECT CAST(COALESCE(COUNT(DISTINCT "name"),0) AS INTEGER) FROM variant_qty;

-- name: GetUniqueWarehouses :many
SELECT DISTINCT "name" FROM variant_qty;

-- name: RemoveQty :exec
DELETE FROM variant_qty
WHERE id = $1;

-- name: RemoveQtyByWarehouseName :exec
DELETE FROM variant_qty
WHERE "name" = $1;

