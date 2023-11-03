-- name: CreateVariantPricing :one
INSERT INTO variant_pricing(
    id,
    variant_id,
    name,
    value,
    isdefault,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateVariantPricing :exec
UPDATE variant_pricing
SET
    name = $1,
    value = $2,
    isdefault = $3
WHERE variant_id IN (
    SELECT id FROM variants
    WHERE sku = $4
);

-- name: GetVariantPricing :many
SELECT 
    name,
    value,
    isdefault
FROM variant_pricing
WHERE variant_id = $1;

-- name: GetVariantPricingBySKU :many
SELECT
    name,
    value,
    isdefault
FROM variant_pricing
WHERE variant_id IN (
    SELECT id FROM variants
    WHERE sku = $1
);

-- name: RemovePricing :exec
DELETE FROM variant_pricing
WHERE id = $1;