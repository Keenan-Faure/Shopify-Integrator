-- +goose Up
CREATE TABLE shopify_location(
    id UUID UNIQUE PRIMARY KEY,
    shopify_warehouse_name VARCHAR(32) NOT NULL,
    shopify_location_id VARCHAR(16) NOT NULL,
    warehouse_name VARCHAR(32) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLe shopify_location;