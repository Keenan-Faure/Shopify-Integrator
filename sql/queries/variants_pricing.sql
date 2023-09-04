-- name: CreateVariantPricing :execresult
INSERT INTO variant_pricing(
    variant_id,
    name,
    value,
    created_at,
    updated_at
) VALUES (
    ?, ?, ?, ?, ?
);

-- name: UpdateVariantPricing :execresult
UPDATE variant_pricing
SET
    name = ?,
    value = ?,
    updated_at = ?
WHERE variant_id = ?;

-- name: GetVariantPricing :many
SELECT 
    name,
    value,
    updated_at
FROM variant_pricing
WHERE variant_id = ?;