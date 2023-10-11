-- name: CreateVID :exec
INSERT INTO shopify_vid(
    id,
    sku,
    shopify_variant_id,
    shopify_inventory_id,
    variant_id,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: UpdateVID :exec
UPDATE shopify_vid
SET
    shopify_variant_id = $1,
    shopify_inventory_id = $2,
    updated_at = $3
WHERE sku = $4;

-- name: GetVIDBySKU :one
SELECT
    sku,
    shopify_variant_id,
    updated_at
FROM shopify_vid
WHERE sku = $1
LIMIT 1;

-- name: GetInventoryIDBySKU :one
select
    sku,
    shopify_inventory_id,
    updated_at
from shopify_vid
where sku = $1
LIMIT 1;