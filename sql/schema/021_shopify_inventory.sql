-- +goose Up
CREATE TABLE shopify_inventory(
    id UUID PRIMARY KEY UNIQUE,
    shopify_location_id VARCHAR(16) NOT NULL DEFAULT 0,
    inventory_item_id VARCHAR(16) NOT NULL DEFAULT 0,
    available INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE shopify_inventory;