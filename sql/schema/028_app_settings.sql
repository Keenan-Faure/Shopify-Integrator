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
    'Enable Shopify Fetch',
    'false',
    NOW(),
    NOW()
), (
    uuid_generate_v4(),
    'app_enable_queue_worker',
    'Enables the queue worker to process queue items.',
    'Enable Queue Worker',
    'true',
    NOW(),
    NOW()
), (
    uuid_generate_v4(),
    'app_shopify_fetch_time',
    'Duration between each product fetch from Shopify.',
    'Fetch Interval',
    'false',
    NOW(),
    NOW()
), (
    uuid_generate_v4(),
    'app_enable_shopify_push',
    'Enables products to be pushed to Shopify.',
    'Enable Shopify Push',
    'false',
    NOW(),
    NOW()
), (
    uuid_generate_v4(),
    'app_queue_size',
    'Maximum amount of queue items that can exist in the queue at any time.',
    'Queue Size',
    '500',
    NOW(),
    NOW()
), (
    uuid_generate_v4(),
    'app_queue_process_limit',
    'Maximum amount of queue items that can be processed each iteration.',
    'Queue Process Limit',
    '20',
    NOW(),
    NOW()
), (
    uuid_generate_v4(),
    'app_queue_cron_time',
    'Interval between each run of the queue worker.',
    'Queue Cron Time',
    '5',
    NOW(),
    NOW()
), (
    uuid_generate_v4(),
    'app_fetch_add_products',
    'Enables the creation of products that does not exist locally when fetching data from Shopify.',
    'Add Fetched Products',
    'false',
    NOW(),
    NOW()
), (
    uuid_generate_v4(),
    'app_fetch_overwrite_products',
    'Enables local data to be overwritten by Shopify data if the product exists locally. Restrictions still apply',
    'Fetch Overwrite Products',
    'false',
    NOW(),
    NOW()
),(
    uuid_generate_v4(),
    'app_fetch_sync_images',
    'Enabled products to be pulled from Shopify when fetching data.',
    'Add Shopify Images',
    'false',
    NOW(),
    NOW()
);

-- +goose Down
DELETE FROM app_settings;