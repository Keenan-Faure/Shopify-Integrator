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
WHERE product_type REGEXP ?
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
WHERE vendor REGEXP ?
LIMIT ? OFFSET ?;

-- name: GetProductsSearchSKU :many
SELECT
    active,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE sku REGEXP ?
LIMIT 5;

-- name: GetProductsSearchTitle :many
SELECT
    active,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE title REGEXP ?
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
