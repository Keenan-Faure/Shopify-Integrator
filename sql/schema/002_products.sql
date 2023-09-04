-- +goose Up
CREATE TABLE products(
    id BINARY(16) PRIMARY KEY UNIQUE NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    active VARCHAR(1) NOT NULL,
    title VARCHAR(255),
    body_html TEXT,
    category VARCHAR(64),
    vendor VARCHAR(64),
    product_type VARCHAR(64),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE (id)
);

-- +goose Down
DROP TABLE products;