-- name: CreateInventoryLocation :exec
INSERT INTO inventory_location(
    ID,
    shopify_location_id,
    inventory_item_id,
    warehouse_name,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: UpdateInventoryLocation :exec
UPDATE inventory_location
SET
    shopify_location_id = $1,
    warehouse_name = $2,
    updated_at = $3
WHERE inventory_item_id = $4;

-- name: GetInventoryLocationLink :one
SELECT
    shopify_location_id,
    warehouse_name,
    updated_at
FROM inventory_location
WHERE inventory_item_id = $1
AND  warehouse_name = $2;

-- name: RemoveLinkBySKU :exec
DELETE FROM inventory_location
WHERE inventory_item_id IN (
    SELECT shopify_inventory_id
    FROM shopify_vid
    WHERE sku = $1
);