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
    active = COALESCE($1, active),
    title = COALESCE($2, title),
    body_html = COALESCE($3, body_html),
    category = COALESCE($4, category),
    vendor = COALESCE($5, vendor),
    product_type = COALESCE($6, product_type),
    updated_at = $7
WHERE product_code = $8;

-- name: UpdateProductByID :exec
UPDATE products
SET
    active = COALESCE($1, active),
    title = COALESCE($2, title),
    body_html = COALESCE($3, body_html),
    category = COALESCE($4, category),
    vendor = COALESCE($5, vendor),
    product_type = COALESCE($6, product_type),
    updated_at = $7
WHERE id = $8;

-- name: UpdateProductBySKU :exec
UPDATE products
SET
    title = COALESCE($1, title),
    body_html = COALESCE($2, body_html),
    category = COALESCE($3, category),
    vendor = COALESCE($4, vendor),
    product_type = COALESCE($5, product_type),
    updated_at = $6
WHERE id = (
    SELECT
        product_id
    FROM variants
    WHERE sku = $7
);

-- name: UpsertProduct :one
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
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT(product_code)
DO UPDATE 
SET
    active = COALESCE($3, products.active),
    title = COALESCE($4, products.title),
    body_html = COALESCE($5, products.body_html),
    category = COALESCE($6, products.category),
    vendor = COALESCE($7, products.vendor),
    product_type = COALESCE($8, products.product_type),
    updated_at = $9
RETURNING *, (xmax = 0) AS inserted;

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
SELECT DISTINCT
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
SELECT DISTINCT
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
WHERE category ILIKE $1
LIMIT $2 OFFSET $3;

-- name: GetProductsByVendor :many
SELECT DISTINCT
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
WHERE vendor ILIKE $1
LIMIT $2 OFFSET $3;

-- name: GetProductsByType :many
SELECT DISTINCT
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
WHERE product_type ILIKE $1
LIMIT $2 OFFSET $3;

-- name: GetProductByCategoryAndType :many
SELECT DISTINCT
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
WHERE category ILIKE $1
AND product_type ILIKE $2
LIMIT $3 OFFSET $4;

-- name: GetProductsByTypeAndVendor :many
SELECT DISTINCT
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
WHERE product_type ILIKE $1
AND vendor ILIKE $2
LIMIT $3 OFFSET $4;

-- name: GetProductsByVendorAndCategory :many
SELECT DISTINCT
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
WHERE vendor ILIKE $1
AND category ILIKE $2
LIMIT $3 OFFSET $4;

-- name: GetProductsByTypeAndCategory :many
SELECT DISTINCT
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
WHERE product_type ILIKE $1
AND category ILIKE $2
LIMIT $3 OFFSET $4;

-- name: GetProductsFilter :many
SELECT DISTINCT
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
WHERE category ILIKE $1
AND product_type ILIKE $2
AND vendor ILIKE $3
LIMIT $4 OFFSET $5;

-- name: GetProductsSearch :many
SELECT
    p.id,
    p.active,
    p.product_code,
    p.title,
    p.category,
    p.vendor,
    p.product_type,
    p.updated_at
FROM products p
INNER JOIN variants v
    ON p.id = v.product_id
WHERE v.sku ILIKE $1
UNION
SELECT
    id,
    active,
    product_code,
    title,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE title ILIKE $1;

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
ORDER BY updated_at DESC
LIMIT $1 OFFSET $2;

-- name: GetActiveProducts :many
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
WHERE active = '1'
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
