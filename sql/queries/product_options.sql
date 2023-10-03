-- name: CreateProductOption :one
INSERT INTO product_options(
    id,
    product_id,
    name
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: UpdateProductOption :one
UPDATE product_options
SET
    name = $1
WHERE product_id = $2
RETURNING *;

-- name: GetProductOptions :many
SELECT
    name
FROM product_options
WHERE product_id = $1;

-- name: GetProductOptionsByCode :many
SELECT
    name
FROM product_options
WHERE product_id IN (
    SELECT id
    FROM products
    WHERE product_code = $1
);