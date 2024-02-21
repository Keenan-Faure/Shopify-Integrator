-- name: CreateCustomerAddress :exec
INSERT INTO customer_address (
    id,
    customer_id,
    address_type,
    updated_at,
    created_at
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: GetAddressByCustomerAndType :one
SELECT
    address_id
FROM customer_address
WHERE customer_id = $1 AND address_type = $2;

-- name: GetAddressByCustomerID :many
SELECT
    address_id
FROM customer_address
WHERE customer_id = $1;
