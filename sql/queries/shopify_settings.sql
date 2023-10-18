-- name: AddShopifySetting :exec
INSERT INTO shopify_settings(
    id,
    key,
    value,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: UpdateShopifySetting :exec
UPDATE shopify_settings
SET
    value = $1
WHERE key = $2;

-- name: RemoveShopifySetting :exec
DELETE FROM shopify_settings
WHERE key = $1;

-- name: GetShopifySettingByKey :one
SELECT
    value,
    updated_at
FROM shopify_settings
WHERE key = $1;