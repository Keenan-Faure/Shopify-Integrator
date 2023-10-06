-- name: CreatePID :exec
INSERT INTO shopify_pid(
    id,
    product_code,
    product_id,
    shopify_product_id,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: UpdatePID :exec
UPDATE shopify_pid
SET
    shopify_product_id = $1,
    updated_at = $2
WHERE product_code = $3;

-- name: GetPIDByProductCode :one
SELECT
    product_code,
    shopify_product_id,
    updated_at
FROM shopify_pid
WHERE product_code = $1
LIMIT 1;