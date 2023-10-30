-- +goose Up
ALTER TABLE shopify_settings ADD COLUMN description VARCHAR(255) NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE shopify_settings DROP COLUMN description;