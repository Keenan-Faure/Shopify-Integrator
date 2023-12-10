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
    title = $2,
    body_html = $3,
    category = $4,
    vendor = $5,
    product_type = $6,
    updated_at = $7
WHERE product_code = $8;

-- name: UpdateProductByID :exec
UPDATE products
SET
    active = $1,
    title = $2,
    body_html = $3,
    category = $4,
    vendor = $5,
    product_type = $6,
    updated_at = $7
WHERE id = $8;

-- name: UpdateProductBySKU :exec
UPDATE products
SET
    active = $1,
    title = $2,
    body_html = $3,
    category = $4,
    vendor = $5,
    product_type = $6,
    updated_at = $7
WHERE id = (
    SELECT
        product_id
    FROM variants
    WHERE sku = $8
);


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
WHERE category LIKE $1
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
WHERE vendor LIKE $1
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
WHERE product_type LIKE $1
LIMIT $2 OFFSET $3;

-- name: GetProductByCategoryAndType :many
SELECT
    id,
    product_code,
    title,
    body_html,
    category,
    vendor,
    product_type
FROM products
WHERE category LIKE $1
AND product_type LIKE $2
LIMIT $3 OFFSET $4;

-- name: GetProductsByTypeAndVendor :many
SELECT
    id,
    product_code,
    title,
    body_html,
    category,
    vendor,
    product_type
FROM products
WHERE product_type LIKE $1
AND vendor LIKE $2
LIMIT $3 OFFSET $4;

-- name: GetProductsByVendorAndCategory :many
SELECT
    id,
    product_code,
    title,
    body_html,
    category,
    vendor,
    product_type
FROM products
WHERE vendor LIKE $1
AND category LIKE $2
LIMIT $3 OFFSET $4;

-- name: GetProductsFilter :many
SELECT
    id,
    product_code,
    title,
    body_html,
    category,
    vendor,
    product_type
FROM products
WHERE category LIKE $1
AND product_type LIKE $2
AND vendor LIKE $3
LIMIT $4 OFFSET $5;

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
    ON p.id = v.product_id
WHERE v.sku LIKE $1
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
WHERE title LIKE $1
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

-- name: RemoveProductByCode :exec
DELETE FROM products
WHERE product_code = $1;

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

-- name: GetProductIDByCode :one
SELECT
    id
FROM products
WHERE product_code = $1;

-- name: GetProductIDs :many
SELECT id FROM products;
