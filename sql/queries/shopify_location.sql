-- name: CreateShopifyLocation :exec
INSERT INTO shopify_location(
    ID,
    shopify_warehouse_name,
    shopify_location_id,
    warehouse_name,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateShopifyLocation :exec
UPDATE shopify_location
SET
    shopify_warehouse_name = $1,
    shopify_location_id = $2,
    updated_at = $3
WHERE warehouse_name = $4;

-- name: GetShopifyLocationByWarehouse :one
SELECT
    shopify_warehouse_name,
    shopify_location_id,
    warehouse_name,
    created_at
FROM shopify_location
WHERE warehouse_name = $1;

-- name: RemoveShopifyLocationMap :exec
DELETE FROM shopify_location
WHERE id = $1;
