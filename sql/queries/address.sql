-- name: CreateAddress :one
INSERT INTO address(
    id,
    customer_id,
    name,
    first_name,
    last_name,
    address1,
    address2,
    suburb,
    city,
    province,
    postal_code,
    company,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
)
RETURNING *;

-- name: UpdateAddress :exec
UPDATE address
SET
    customer_id = $1,
    first_name = $2,
    last_name = $3,
    address1 = $4,
    address2 = $5,
    suburb = $6,
    city = $7,
    province = $8,
    postal_code = $9,
    company = $10,
    updated_at = $11
WHERE id = $12;

-- name: UpdateAddressByNameAndCustomer :exec
UPDATE address
SET
    customer_id = $1,
    first_name = $2,
    last_name = $3,
    address1 = $4,
    address2 = $5,
    suburb = $6,
    city = $7,
    province = $8,
    postal_code = $9,
    company = $10,
    updated_at = $11
WHERE name = $12 AND
customer_id = $13;

-- name: GetAddressByCustomer :many
SELECT
    id,
    first_name,
    last_name,
    address1,
    address2,
    suburb,
    city,
    province,
    postal_code,
    company,
    updated_at
FROM address
WHERE customer_id = $1;

-- name: RemoveAddress :exec
DELETE FROM address
WHERE id = $1;
