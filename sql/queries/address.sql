-- name: CreateAddress :execresult
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
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateAddress :execresult
UPDATE address
SET
    customer_id = ?,
    first_name = ?,
    last_name = ?,
    address1 = ?,
    address2 = ?,
    suburb = ?,
    city = ?,
    province = ?,
    postal_code = ?,
    company = ?,
    updated_at = ?
WHERE id = ?;

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
WHERE customer_id = ?;