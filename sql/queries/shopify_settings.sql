-- name: AddShopifySetting :exec
INSERT INTO shopify_settings(
    id,
    key,
    description,
    field_name,
    value,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: UpdateShopifySetting :exec
UPDATE shopify_settings
SET
    value = $1,
    updated_at = $2
WHERE key = $3;

-- name: GetShopifySettingByKey :one
SELECT
    id,
    key,
    description,
    field_name,
    value,
    updated_at
FROM shopify_settings
WHERE key = $1;

-- name: GetShopifySettings :many
SELECT
    id,
    key,
    description,
    field_name,
    value,
    updated_at
FROM shopify_settings;

-- name: RemoveShopifySetting :exec
DELETE FROM shopify_settings
WHERE key = $1;