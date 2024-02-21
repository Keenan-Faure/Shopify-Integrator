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
