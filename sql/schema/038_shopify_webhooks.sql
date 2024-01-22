-- +goose Up
CREATE TABLE shopify_webhooks(
    id UUID PRIMARY KEY NOT NULL UNIQUE,
    shopify_webhook_id VARCHAR(64) NOT NULL UNIQUE,
    webhook_url TEXT NOT NULL,
    topic VARCHAR(32) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE shopify_webhooks;