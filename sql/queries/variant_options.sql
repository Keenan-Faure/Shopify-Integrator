-- name: CreateVariantOption :execresult
INSERT INTO variant_options(
    id,
    product_id,
    name,
    value
) VALUES (?, ?, ?, ?);

-- name: UpdateVariantOption
UPDATE variant_options SET
name = ?
value = ?;