-- name: CreateVariantPricing :one
INSERT INTO variant_pricing(
    id,
    variant_id,
    name,
    value,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateVariantPricing :one
UPDATE variant_pricing
SET
    name = $1,
    value = $2
WHERE variant_id = $3
RETURNING *;

-- name: GetVariantPricing :many
SELECT 
    name,
    value
FROM variant_pricing
WHERE variant_id = $1;