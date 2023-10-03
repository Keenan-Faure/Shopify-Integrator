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

-- name: UpdateVariantPricing :exec
UPDATE variant_pricing
SET
    value = $1
WHERE variant_id IN (
    SELECT id FROM variants
    WHERE sku = $2
) AND name = $3;

-- name: GetVariantPricing :many
SELECT 
    name,
    value
FROM variant_pricing
WHERE variant_id = $1;

-- name: GetVariantPricingBySKU :many
SELECT
    name,
    value
FROM variant_pricing
WHERE variant_id IN (
    SELECT id FROM variants
    WHERE sku = $1
);

-- name: RemovePricing :exec
DELETE FROM variant_pricing
WHERE id = $1;