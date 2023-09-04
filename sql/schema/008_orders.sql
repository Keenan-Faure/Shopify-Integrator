-- +goose Up
CREATE TABLE orders(
    id BINARY(16) PRIMARY KEY UNIQUE NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    customer_id BINARY(16) NOT NULL,
    notes VARCHAR(255) DEFAULT '',
    web_code VARCHAR(32) UNIQUE,
    tax_total DECIMAL(10, 2) DEFAULT 0.00,
    order_total DECIMAL(10, 2) DEFAULT 0.00,
    shipping_total DECIMAL(10, 2) DEFAULT 0.00,
    discount_total DECIMAL(10, 2) DEFAULT 0.00,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE orders;