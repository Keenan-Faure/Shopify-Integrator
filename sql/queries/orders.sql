-- name: CreateOrder :one
INSERT INTO orders(
    id,
    status,
    notes,
    web_code,
    tax_total,
    order_total,
    shipping_total,
    discount_total,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: UpdateOrder :one
UPDATE orders
SET
    notes = $1,
    status = $2,
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
    id,
    notes,
    status,
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
    o.id,
    o.notes,
    o.status,
    o.web_code,
    o.tax_total,
    o.order_total,
    o.shipping_total,
    o.discount_total,
    o.updated_at
FROM orders o 
WHERE orders.id in (
    SELECT order_id FROM customerorders
    WHERE customer_id = $1
);

-- name: GetOrders :many
SELECT
    id,
    notes,
    status,
    web_code,
    tax_total,
    order_total,
    shipping_total,
    discount_total,
    updated_at
FROM orders
ORDER BY updated_at DESC
LIMIT $1 OFFSET $2;

-- name: GetOrderByWebCode :one
SELECT
    id,
    notes,
    status,
    web_code,
    tax_total,
    order_total,
    shipping_total,
    discount_total,
    updated_at
FROM orders
WHERE web_code = $1;

-- name: GetOrdersSearchWebCode :many
SELECT
    id,
    notes,
    status,
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
    o.id,
    o.notes,
    o.status,
    o.web_code,
    o.tax_total,
    o.order_total,
    o.shipping_total,
    o.discount_total,
    o.updated_at
FROM orders o
WHERE o.id in (
    SELECT order_id FROM customerorders co
    INNER JOIN customers c
    ON co.customer_id = c.id
    WHERE CONCAT(c.first_name, ' ', c.last_name) SIMILAR TO $1
    OR c.first_name LIKE $1
    OR c.last_name LIKE $1
);

-- name: FetchOrderStatsPaid :many
SELECT
	COUNT(id) AS "count",
	to_char(created_at, 'YYYY-MM-DD') AS "day"
FROM orders
WHERE
    created_at > current_date at time zone 'UTC' - interval '7 day' AND
    "status" = 'paid'
GROUP BY "day"
ORDER BY "day" DESC;

-- name: FetchOrderStatsNotPaid :many
SELECT
	COUNT(id) AS "count",
	to_char(created_at, 'YYYY-MM-DD') AS "day"
FROM orders
WHERE
    created_at > current_date at time zone 'UTC' - interval '7 day' AND
    "status" != 'paid'
GROUP BY "day"
ORDER BY "day" DESC;

-- name: RemoveOrder :exec
DELETE FROM orders
WHERE id = $1;