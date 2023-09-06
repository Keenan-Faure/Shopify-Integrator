-- +goose Up
CREATE TABLE order_lines(
    id BINARY(16) PRIMARY KEY UNIQUE NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    order_id BINARY(16) NOT NULL,
    line_type VARCHAR(16),
    sku VARCHAR(64) NOT NULL,
    price DECIMAL(9, 2) DEFAULT 0.00,
    barcode INTEGER DEFAULT 0,
    qty INTEGER DEFAULT 0,
    tax_total DECIMAL(9, 2) DEFAULT 0.00,
    tax_rate  DECIMAL(9, 2) DEFAULT 0.00,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE order_lines;

