-- +goose Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

INSERT INTO app_settings(
    id,
    key,
    description,
    field_name,
    value,
    created_at,
    updated_at
) VALUES (
    uuid_generate_v4(),
    'app_enable_shopify_fetch',
    'Enables the automatic pulling of products from Shopify.',
    'App Enable Shopify Fetch',
    'false',
    NOW(),
    NOW()
), (
    uuid_generate_v4(),
    'app_enable_queue_worker',
    'Enables the queue worker to process queue items.',
    'App Enable Queue Worker',
    'false',
    NOW(),
    NOW()
), (
    uuid_generate_v4(),
    'app_shopify_fetch_time',
    'Duration between each product fetch from Shopify.',
    'App Shopify Fetch Time',
    'false',
    NOW(),
    NOW()
), (
    uuid_generate_v4(),
    'app_enable_shopify_push',
    'Enables products to be pushed to Shopify.',
    'App Enable Shopify Push',
    'false',
    NOW(),
    NOW()
), (
    uuid_generate_v4(),
    'app_queue_size',
    'Maximum amount of queue items that can exist in the queue at any time.',
    'App Queue Size',
    '100',
    NOW(),
    NOW()
), (
    uuid_generate_v4(),
    'app_queue_process_limit',
    'Maximum amount of queue items that can be processed each iteration.',
    'App Queue Process Limit',
    '10',
    NOW(),
    NOW()
), (
    uuid_generate_v4(),
    'app_queue_cron_time',
    'Interval between each run of the queue worker.',
    'App Queue Cron Time',
    '7',
    NOW(),
    NOW()
), (
    uuid_generate_v4(),
    'app_fetch_add_products',
    'Enables the creation of products that does not exist locally when fetching data from Shopify.',
    'App Fetch Add Products',
    'false',
    NOW(),
    NOW()
), (
    uuid_generate_v4(),
    'app_fetch_overwrite_products',
    'Enables local data to be overwritten by Shopify data if the product exists locally.',
    'App Fetch Overwrite Products',
    'false',
    NOW(),
    NOW()
), (
    uuid_generate_v4(),
    'app_fetch_create_price_tier_enabled',
    'Enables price tiers to be created when fetching data from Shopify.',
    'App Fetch Create Price Tier',
    'false',
    NOW(),
    NOW()
), (
    uuid_generate_v4(),
    'app_fetch_sync_images',
    'Enabled products to be pulled from Shopify when fetching data.',
    'App Fetch Sync Images',
    'false',
    NOW(),
    NOW()
);

-- +goose Down
DELETE FROM app_settings;