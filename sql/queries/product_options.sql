-- name: CreateProductOption :one
INSERT INTO product_options(
    id,
    product_id,
    "name",
    position
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateProductOption :one
UPDATE product_options
SET
    "name" = COALESCE($1, "name"),
    position = COALESCE($2, position)
WHERE product_id = $3
and position = $4
RETURNING *;

-- name: GetProductOptions :many
SELECT
    "name",
    position
FROM product_options
WHERE product_id = $1
ORDER BY position ASC;

-- name: GetProductOptionsByCode :many
SELECT
    "name",
    position
FROM product_options
WHERE product_id IN (
    SELECT id
    FROM products
    WHERE product_code = $1
)
ORDER BY position ASC;