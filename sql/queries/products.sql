-- name: CreateProduct :one
INSERT INTO products(
    id,
    product_code,
    active,
    title,
    body_html,
    category,
    vendor,
    product_type,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: UpdateProduct :exec
UPDATE products
SET
    active = $1,
    product_code = $2,
    title = $3,
    body_html = $4,
    category = $5,
    vendor = $6,
    product_type = $7,
    updated_at = $8
WHERE product_code = $9;

-- name: GetProductByID :one
SELECT
    active,
    product_code,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE id = $1;

-- name: GetProductByProductCode :one
SELECT
    active,
    product_code,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE product_code = $1;

-- name: GetProductsByCategory :many
SELECT
    id,
    product_code,
    title,
    body_html,
    category,
    vendor,
    product_type
FROM products
WHERE LOWER(category) LIKE CONCAT('%',LOWER($1),'%')
LIMIT $2 OFFSET $3;

-- name: GetProductsByType :many
SELECT
    id,
    product_code,
    title,
    body_html,
    category,
    vendor,
    product_type
FROM products
WHERE LOWER(product_type) LIKE CONCAT('%',LOWER($1),'%')
LIMIT $2 OFFSET $3;

-- name: GetProductsByVendor :many
SELECT
    id,
    product_code,
    title,
    body_html,
    category,
    vendor,
    product_type
FROM products
WHERE LOWER(vendor) LIKE CONCAT('%',LOWER($1),'%')
LIMIT $2 OFFSET $3;

-- name: GetProductsSearchSKU :many
SELECT
    p.id,
    p.product_code,
    p.title,
    p.category,
    p.vendor,
    p.product_type
FROM products p
INNER JOIN variants v
ON p.id = variants.product_id
WHERE LOWER(v.sku) LIKE CONCAT('%',LOWER($1),'%')
LIMIT 5;

-- name: GetProductsSearchTitle :many
SELECT
    id,
    product_code,
    title,
    category,
    vendor,
    product_type
FROM products
WHERE LOWER(title) LIKE CONCAT('%',LOWER($1),'%')
LIMIT 5;

-- name: GetProducts :many
SELECT
    id,
    active,
    product_code,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
LIMIT $1 OFFSET $2;

-- name: RemoveProduct :exec
DELETE FROM products
WHERE id = $1;

-- name: GetVariantOptionsByProductCode :many
SELECT
    v.sku,
    v.option1,
    v.option2,
    v.option3
FROM variants v
WHERE v.product_id IN (
    SELECT product_id
    FROM products
    WHERE product_code = $1
);
