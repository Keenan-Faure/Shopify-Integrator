-- name: CreateCustomer :one
INSERT INTO customers(
    first_name,
    last_name,
    email,
    phone,
    created_at,
    updated_at
) VALUES(
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateCustomer :one
UPDATE customers
SET
    first_name = $1,
    last_name = $2,
    email = $3,
    phone = $4,
    updated_at = $5
WHERE id = $6
RETURNING *;

-- name: GetCustomers :many
SELECT
    first_name,
    last_name,
    email,
    phone,
    updated_at
FROM customers
LIMIT $1 OFFSET $2;

-- name: GetCustomerByID :one
SELECT
    first_name,
    last_name,
    email,
    phone,
    updated_at
FROM customers
WHERE id = $1;

-- name: GetCustomersByName :many
SELECT
    first_name,
    last_name,
    email,
    phone,
    updated_at
FROM customers
WHERE CONCAT(first_name, ' ', last_name) SIMILAR TO $1
LIMIT 10;
