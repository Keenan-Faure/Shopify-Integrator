-- name: CreateVariantPricing :one
INSERT INTO variant_pricing(
    id,
    variant_id,
    "name",
    "value",
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
    "name" = COALESCE($1, "name"),
    "value" = COALESCE($2, "value"),
    isdefault = COALESCE($3, isdefault),
    updated_at = $4
WHERE variant_id IN (
    SELECT id FROM variants
    WHERE sku = $5
) AND "name" = $6;

-- name: GetPriceTierBySKU :one
SELECT 
    "name",
    "value",
    isdefault
FROM variant_pricing
WHERE variant_id IN (
    SELECT id FROM variants
    WHERE sku = $1
) AND "name" = $2;

-- name: GetVariantPricing :many
SELECT 
    "name",
    "value",
    isdefault
FROM variant_pricing
WHERE variant_id = $1;

-- name: GetVariantPricingBySKU :many
SELECT
    "name",
    "value",
    isdefault
FROM variant_pricing
WHERE variant_id IN (
    SELECT id FROM variants
    WHERE sku = $1
);

-- name: GetCountOfUniquePrices :one
SELECT CAST(COALESCE(COUNT(DISTINCT "name"),0) AS INTEGER) FROM variant_pricing;

-- name: GetUniquePriceTiers :many
SELECT DISTINCT "name" FROM variant_pricing;

-- name: RemovePricing :exec
DELETE FROM variant_pricing
WHERE id = $1;