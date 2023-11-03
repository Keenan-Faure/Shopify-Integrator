-- name: CreateShopifyCollection :exec
INSERT INTO shopify_collections(
    ID,
    product_collection,
    shopify_collection_id,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5);

-- name: GetShopifyCollection :one
SELECT
    product_collection,
    shopify_collection_id,
    updated_at
FROM shopify_collections
WHERE product_collection = $1
LIMIT 1;

-- name: UpdateShopifyCollection :exec
UPDATE shopify_collections
SET
    shopify_collection_id = $1
WHERE product_collection = $2;
