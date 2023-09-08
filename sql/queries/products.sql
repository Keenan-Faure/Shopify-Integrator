-- name: CreateProduct :one
INSERT INTO products(
    active,
    title,
    body_html,
    category,
    vendor,
    product_type,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET
    active = $1,
    title = $2,
    body_html = $3,
    category = $4,
    vendor = $5,
    product_type = $6,
    updated_at = $7
WHERE id = $8
RETURNING *;

-- name: GetProductByID :one
SELECT
    active,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE id = $1;

-- name: GetProductsByCategory :many
SELECT
    id,
    title,
    body_html,
    category,
    vendor,
    product_type
FROM products
WHERE category LIKE $1
LIMIT $2 OFFSET $3;

-- name: GetProductsByType :many
SELECT
    id,
    title,
    body_html,
    category,
    vendor,
    product_type
FROM products
WHERE product_type LIKE $1
LIMIT $2 OFFSET $3;

-- name: GetProductsByVendor :many
SELECT
    id,
    title,
    body_html,
    category,
    vendor,
    product_type
FROM products
WHERE vendor LIKE $1
LIMIT $2 OFFSET $3;

-- name: GetProductsSearchSKU :many
SELECT
    p.id,
    p.title,
    p.category,
    p.vendor,
    p.product_type
FROM products p
INNER JOIN variants v
ON p.id = variants.product_id
WHERE v.sku LIKE $1
LIMIT 5;

-- name: GetProductsSearchTitle :many
SELECT
    id,
    title,
    category,
    vendor,
    product_type
FROM products
WHERE title LIKE $1
LIMIT 5;

-- name: GetProducts :many
SELECT
    active,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
LIMIT $1 OFFSET $2;
