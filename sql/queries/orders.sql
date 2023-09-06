-- name: CreateOrder :execresult
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
    ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateOrder :execresult
UPDATE orders
SET
    customer_id = ?,
    notes = ?,
    web_code = ?,
    tax_total = ?,
    order_total = ?,
    shipping_total = ?,
    discount_total = ?,
    updated_at = ?
WHERE id = ?;

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
WHERE id = ?;

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
WHERE customer_id = ?;

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
LIMIT ? OFFSET ?;

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
WHERE web_code LIKE ?
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
WHERE CONCAT(c.first_name, ' ', c.last_name) LIKE ?
LIMIT 10;

