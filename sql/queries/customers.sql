-- name: CreateCustomer :execresult
INSERT INTO customers(
    first_name,
    last_name,
    created_at,
    updated_at
) VALUES(
    ?, ?, ?, ?
);

-- name: UpdateCustomer :execresult
UPDATE customers
SET
    first_name = ?,
    last_name = ?,
    updated_at = ?
WHERE id = ?;

-- name: GetCustomerByID :one
SELECT
    first_name,
    last_name,
    updated_at
FROM customers
WHERE id = ?;

-- name: GetCustomersByName :many
SELECT
    first_name,
    last_name,
    updated_at
FROM customers
WHERE CONCAT(first_name, ' ', last_name) REGEXP ?
LIMIT 10;
