-- name: CreateProduct :execresult
INSERT INTO products(
    id,
    active,
    title,
    body_html,
    category,
    product_type,
    created_at,
    updated_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateProduct :execresult
UPDATE products
SET
    active = ?,
    title = ?,
    body_html = ?,
    category = ?,
    product_type = ?,
    updated_at = ?
WHERE id = ?;

-- name: GetProductByID :one
SELECT
    active,
    title,
    body_html,
    category,
    product_type,
    updated_at
FROM products
WHERE id = ?;

-- name: GetProductByActiveStatus :many
SELECT
    active,
    title,
    body_html,
    category,
    product_type,
    updated_at
FROM products
WHERE active = ?;
 

