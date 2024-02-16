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

-- name: UpsertVariant :one
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
ON CONFLICT(sku)
DO UPDATE 
SET
    option1 = COALESCE($4, variants.option1),
    option2 = COALESCE($5, variants.option2),
    option3 = COALESCE($6, variants.option3),
    barcode = COALESCE($7, variants.barcode),
    updated_at = $9
RETURNING *, (xmax = 0) AS inserted;


-- name: UpdateVariant :exec
UPDATE variants
SET
    option1 = COALESCE($1, option1),
    option2 = COALESCE($2, option2),
    option3 = COALESCE($3, option3),
    barcode = COALESCE($4, barcode),
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

-- name: GetVariants :many
SELECT id FROM variants;

-- name: GetUnindexedVariants :many
SELECT
    id
FROM variants
WHERE id NOT IN (
    SELECT
        variant_id
    FROM variant_qty
);

-- name: RemoveVariant :exec
DELETE FROM variants
WHERE id = $1 AND product_id = $2;
