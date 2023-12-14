-- +goose Up
CREATE TABLE shopify_settings(
    id UUID UNIQUE PRIMARY KEY,
    key VARCHAR(64) UNIQUE NOT NULL,
    value TEXT NOT NULL DEFAULT '',
    field_name VARCHAR(64) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE shopify_settings;