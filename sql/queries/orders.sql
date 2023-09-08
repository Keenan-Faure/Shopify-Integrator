-- name: CreateOrder :one
INSERT INTO orders(
    customer_id,
    notes,
    web_code,
    tax_total,
    order_total,
    shipping_total,
    discount_total,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: UpdateOrder :one
UPDATE orders
SET
    customer_id = $1,
    notes = $2,
    web_code = $3,
    tax_total = $4,
    order_total = $5,
    shipping_total = $6,
    discount_total = $7,
    updated_at = $8
WHERE id = $9
RETURNING *;

-- name: GetOrderByID :one
SELECT
    customer_id,
    notes,
    web_code,
    tax_total,
    order_total,
    shipping_total,
    discount_total,
    updated_at,
    created_at
FROM orders
WHERE id = $1;

-- name: GetOrderByCustomer :many
SELECT
    customer_id,
    notes,
    web_code,
    tax_total,
    order_total,
    shipping_total,
    discount_total,
    updated_at
FROM orders
WHERE customer_id = $1;

-- name: GetOrders :many
SELECT
    customer_id,
    notes,
    web_code,
    tax_total,
    order_total,
    shipping_total,
    discount_total,
    updated_at
FROM orders
LIMIT $1 OFFSET $2;

-- name: GetOrdersSearchWebCode :many
SELECT
    notes,
    web_code,
    tax_total,
    order_total,
    shipping_total,
    discount_total,
    updated_at
FROM orders
WHERE web_code SIMILAR TO $1
LIMIT 10;

-- name: GetOrdersSearchByCustomer :many
SELECT
    o.notes,
    o.web_code,
    o.tax_total,
    o.order_total,
    o.shipping_total,
    o.discount_total,
    o.updated_at
FROM orders o
INNER JOIN customers c
ON o.customer_id = c.id
WHERE CONCAT(c.first_name, ' ', c.last_name) SIMILAR TO $1
LIMIT 10;

