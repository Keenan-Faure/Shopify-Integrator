-- name: CreateAddress :one
INSERT INTO address(
    customer_id,
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
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING *;

-- name: UpdateAddress :one
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
WHERE id = $12
RETURNING *;

-- name: GetAddressByCustomer :many
SELECT
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