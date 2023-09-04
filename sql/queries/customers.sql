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
    c.first_name,
    c.last_name,
    c.updated_at,
    a.address1,
    a.address2,
    a.first_name,
    a.last_name,
    a.suburb,
    a.city,
    a.province,
    a.company,
    a.postal_code,
    a.updated_at
FROM customers c
INNER JOIN address a
ON c.id = a.customer_id
WHERE c.id = ?;

-- name: GetCustomersByName :many
SELECT
    c.first_name,
    c.last_name,
    c.updated_at,
    a.address1,
    a.address2,
    a.first_name,
    a.last_name,
    a.suburb,
    a.city,
    a.province,
    a.company,
    a.postal_code,
    a.updated_at
FROM customers c
INNER JOIN address a
ON c.id = a.customer_id
WHERE CONCAT(first_name, ' ', last_name) REGEXP ?
LIMIT ? OFFSET ?;

