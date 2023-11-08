-- name: CreateVariant :one
INSERT INTO variants(
    id,
    product_id,
    sku,
    option1,
    option2,
    option3,
    barcode,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: UpdateVariant :exec
UPDATE variants
SET
    option1 = $1,
    option2 = $2,
    option3 = $3,
    barcode = $4,
    updated_at = $5
WHERE sku = $6;

-- name: GetProductVariants :many
SELECT
    id,
    sku,
    option1,
    option2,
    option3,
    barcode,
    updated_at
FROM variants
WHERE product_id = $1;

-- name: GetVariantBySKU :one
SELECT
    id,
    product_id,
    sku,
    option1,
    option2,
    option3,
    barcode
FROM variants
WHERE sku = $1;

-- name: GetVariantIDByCode :one
SELECT
    id
FROM variants
WHERE sku = $1;

-- name: GetVariantByVariantID :one
SELECT
    *
FROM variants
WHERE id = $1;

-- name: RemoveVariant :exec
DELETE FROM variants
WHERE id = $1;
