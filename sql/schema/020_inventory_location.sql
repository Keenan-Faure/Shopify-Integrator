-- +goose Up
CREATE TABLE inventory_location(
    id UUID UNIQUE PRIMARY KEY,
    shopify_location_id VARCHAR(16) NOT NULL DEFAULT 0,
    inventory_item_id VARCHAR(16) NOT NULL DEFAULT 0,
    warehouse_name VARCHAR(32) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE inventory_location;