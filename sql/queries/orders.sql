-- name: CreateOrder :one
INSERT INTO orders(
    id,
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
    notes = $1,
    web_code = $2,
    tax_total = $3,
    order_total = $4,
    shipping_total = $5,
    discount_total = $6,
    updated_at = $7
WHERE id = $8
RETURNING *;

-- name: GetOrderByID :one
SELECT
    id,
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
    o.id,
    o.notes,
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
    web_code,
    tax_total,
    order_total,
    shipping_total,
    discount_total,
    updated_at
FROM orders
LIMIT $1 OFFSET $2;

-- name: GetOrderByWebCode :one
SELECT
    id,
    notes,
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
    WHERE CONCAT(LOWER(c.first_name), ' ', LOWER(c.last_name)) SIMILAR TO LOWER($1)
    AND LOWER(c.first_name) LIKE CONCAT('%',LOWER($1),'%')
    AND LOWER(c.last_name) LIKE CONCAT('%',LOWER($1),'%')
);

-- name: RemoveOrder :exec
DELETE FROM orders
WHERE id = $1;