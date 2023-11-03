-- name: CreateShopifyInventoryRecord :exec
INSERT INTO shopify_inventory(
    id,
    shopify_location_id,
    inventory_item_id,
    available,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: GetShopifyInventory :one
SELECT
    available,
    created_at
FROM shopify_inventory
WHERE inventory_item_id = $1 AND
shopify_location_id = $2;

-- name: UpdateShopifyInventoryRecord :exec
UPDATE shopify_inventory
SET
    available = $1,
    updated_at = $2
WHERE shopify_location_id = $3
AND inventory_item_id = $4;

-- name: RemoveShopifyInventoryRecord :exec
DELETE FROM shopify_inventory
WHERE shopify_location_id = $1
AND inventory_item_id = $2;