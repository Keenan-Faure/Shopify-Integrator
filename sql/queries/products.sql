-- name: CreateProduct :execresult
INSERT INTO products(
    id,
    active,
    title,
    body_html,
    category,
    vendor,
    product_type,
    created_at,
    updated_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?
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
WHERE active = ?
AND category IN (?);

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
WHERE active = ?
AND category IN (?)
AND product_type IN (?)
AND vendor IN (?);

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
WHERE active = ?
AND product_type in (?);

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
WHERE active = ?
AND vendor IN (?);
 

