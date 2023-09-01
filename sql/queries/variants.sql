-- name: CreateVariant :execresult
INSERT INTO variants(
    id,
    product_id,
    sku,
    option1,
    option2,
    option3,
    barcode
) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: UpdateVariant :execresult
UPDATE variants
SET
    option1 = ?,
    option2 = ?,
    option3 = ?,
    barcode = ?
WHERE sku = ?;

-- name: GetProductVariants :many
SELECT
    sku,
    option1,
    option2,
    option3,
    barcode
FROM variants
WHERE product_id = ?;

-- name: GetVariantBySKU :one
SELECT
    sku,
    option1,
    option2,
    option3,
    barcode
FROM variants
WHERE sku = ?;
