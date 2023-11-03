-- +goose Up
ALTER TABLE shopify_vid ADD COLUMN shopify_inventory_id VARCHAR(16) NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE shopify_vid DROP COLUMN shopify_inventory_id;