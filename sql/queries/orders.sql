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
WHERE web_code REGEXP ?
LIMIT 10;

