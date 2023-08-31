-- name: CreateVariantPricing :execresult
INSERT INTO variant_pricing(
    id,
    variant_id,
    name,
    value
    created_at,
    updated_at
) VALUES (?, ?, ?, ?);

-- name: UpdateVariantPricing
UPDATE variant_pricing SET
name = ?
value = ?
updated_at = ?;