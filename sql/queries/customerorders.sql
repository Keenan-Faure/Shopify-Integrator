-- name: CreateCustomerOrder :exec
INSERT INTO customerorders (
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
FROM customerorders
WHERE order_id = $1;

-- name: GetOrdersByCustomerID :many
SELECT
    order_id
FROM customerorders
WHERE customer_id = $1;
