-- name: CreateProductOption :execresult
INSERT INTO product_options(
    product_id,
    name,
    value
) VALUES (
    ?, ?, ?
);

-- name: UpdateProductOption :execresult
UPDATE product_options
SET
    name = ?,
    value = ?
WHERE product_id = ?;

-- name: GetProductOptions :many
SELECT
    name,
    value
FROM product_options
WHERE id = ?;