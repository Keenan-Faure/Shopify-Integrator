-- name: CreateProduct :execresult
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
    ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateProduct :execresult
UPDATE products
SET
    active = ?,
    title = ?,
    body_html = ?,
    category = ?,
    vendor = ?,
    product_type = ?,
    updated_at = ?
WHERE id = ?;

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
WHERE id = ?;

-- name: GetProductsByCategory :many
SELECT
    active,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE category REGEXP ?
LIMIT ? OFFSET ?;

-- name: GetProductsByFilter :many
SELECT
    active,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE category IN (?)
AND product_type IN (?)
AND vendor IN (?)
LIMIT ? OFFSET ?;

-- name: GetProductsByType :many
SELECT
    active,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE product_type LIKE ?
LIMIT ? OFFSET ?;

-- name: GetProductsByVendor :many
SELECT
    active,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE vendor LIKE ?
LIMIT ? OFFSET ?;

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
WHERE v.sku LIKE ?
LIMIT 5;

-- name: GetProductsSearchTitle :many
SELECT
    id,
    title,
    category,
    vendor,
    product_type
FROM products
WHERE title LIKE ?
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
LIMIT ? OFFSET ?;
