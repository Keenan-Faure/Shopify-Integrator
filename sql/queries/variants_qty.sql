-- name: CreateVariantQty :execresult
INSERT INTO variant_qty(
    id,
    variant_id,
    name,
    value
    created_at,
    updated_at
) VALUES (?, ?, ?, ?);

-- name: UpdateVariantQty
UPDATE variant_qty SET
name = ?
value = ?
updated_at = ?;