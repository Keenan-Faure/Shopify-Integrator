-- +goose Up
CREATE TABLE shopify_settings(
    id UUID PRIMARY KEY UNIQUE,
    key VARCHAR(64) UNIQUE NOT NULL,
    value TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE shopify_settings;