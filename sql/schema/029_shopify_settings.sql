-- +goose Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

INSERT INTO shopify_settings(
    id,
    key,
    description,
    field_name,
    value,
    created_at,
    updated_at
) VALUES (
    uuid_generate_v4(),
    'shopify_enable_dynamic_sku_search',
    'Enables the dynamic searching of SKUs on Shopify when adding new products. If disabled, only first product SKU will be considered.',
    'Shopify Enable Dynamic SKU Search',
    'true',
    NOW(),
    NOW()
);

-- +goose Down
DELETE FROM shopify_settings;