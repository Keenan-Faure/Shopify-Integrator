-- name: CreateVariant :execresult
INSERT INTO variants(
    product_id,
    sku,
    option1,
    option2,
    option3,
    barcode,
    created_at,
    updated_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateVariant :execresult
UPDATE variants
SET
    option1 = ?,
    option2 = ?,
    option3 = ?,
    barcode = ?,
    updated_at = ?
WHERE sku = ?;

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
WHERE product_id = ?;

-- name: GetVariantBySKU :one
SELECT
    id,
    sku,
    option1,
    option2,
    option3,
    barcode
FROM variants
WHERE sku = ?;
