-- name: CreateProductOption :one
INSERT INTO product_options(
    product_id,
    name,
    value
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: UpdateProductOption :one
UPDATE product_options
SET
    name = $1,
    value = $2
WHERE product_id = $3
RETURNING *;

-- name: GetProductOptions :many
SELECT
    name,
    value
FROM product_options
WHERE id = $1;