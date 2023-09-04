-- name: CreateVariantQty :execresult
INSERT INTO variant_qty(
    variant_id,
    name,
    value,
    created_at,
    updated_at
) VALUES (
    ?, ?, ?, ?, ?
);

-- name: UpdateVariantQty :execresult
UPDATE variant_qty
SET
    name = ?,
    value = ?,
    updated_at = ?
WHERE variant_id = ?;

-- name: GetVariantQty :many
SELECT 
    name,
    value,
    updated_at
FROM variant_qty
WHERE variant_id = ?;