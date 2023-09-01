-- name: CreateProductOption :execresult
INSERT INTO product_options(
    id,
    product_id,
    name,
    value
) VALUES (?, ?, ?, ?);

-- name: UpdateProductOption :execresult
UPDATE product_options
SET
    name = ?,
    value = ?
WHERE product_id = ?;

SELECT
    name,
    value
FROM product_options
WHERE product_id = ?;