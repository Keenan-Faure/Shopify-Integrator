-- name: CreateCustomerOrder :exec
INSERT INTO customer_orders (
    id,
    customer_id,
    order_id,
    updated_at,
    created_at
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: GetCustomerByOrderID :one
SELECT
    customer_id
FROM customer_orders
WHERE order_id = $1;

-- name: GetOrdersByCustomerID :many
SELECT
    order_id
FROM customer_orders
WHERE customer_id = $1;

-- name: GetOrderIDByCustomerID :one
SELECT
    order_id
FROM customer_orders
WHERE customer_id = $1 AND order_id = $2;

-- name: RemoveCustomerOrdersByOrderID :exec
DELETE FROM customer_orders
WHERE order_id = (
    SELECT id
    FROM orders
    WHERE web_code = $1
);
