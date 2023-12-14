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
    'shopify_default_price_tier',
    'Price tier to use when pushing pricing to Shopify.',
    'Shopify Default Price Tier',
    'default',
    NOW(),
    NOW()
), (
    uuid_generate_v4(),
    'shopify_default_compare_at_price_tier',
    'Price tier to use as the compare at tier when pushing pricing to Shopify.',
    'Shopify Default Compare At Price Tier',
    'default',
    NOW(),
    NOW()
), (
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