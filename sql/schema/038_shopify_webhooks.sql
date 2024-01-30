-- +goose Up
CREATE TABLE shopify_webhooks(
    id UUID PRIMARY KEY NOT NULL UNIQUE,
    shopify_webhook_id VARCHAR(64) NOT NULL UNIQUE,
    webhook_url TEXT NOT NULL,
    topic VARCHAR(32) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

INSERT INTO shopify_webhooks(
    id,
    shopify_webhook_id,
    webhook_url,
    topic,
    created_at,
    updated_at
) VALUES (
    uuid_generate_v4(),
    '',
    '',
    '',
    NOW(),
    NOW()
);

-- +goose Down
DROP TABLE shopify_webhooks;