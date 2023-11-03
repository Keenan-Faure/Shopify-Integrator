-- +goose Up
CREATE TABLE shopify_collections(
    id UUID PRIMARY KEY UNIQUE,
    product_collection VARCHAR(64) UNIQUE,
    shopify_collection_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLe shopify_collections;