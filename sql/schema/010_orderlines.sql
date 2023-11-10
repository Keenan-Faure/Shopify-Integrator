-- +goose Up
CREATE TABLE order_lines(
    id UUID UNIQUE PRIMARY KEY,
    order_id UUID NOT NULL,
    line_type VARCHAR(16),
    sku VARCHAR(64) NOT NULL,
    price DECIMAL(9, 2) DEFAULT 0.00,
    barcode INTEGER DEFAULT 0,
    qty INTEGER DEFAULT 0,
    tax_total DECIMAL(9, 2) DEFAULT 0.00,
    tax_rate  DECIMAL(9, 2) DEFAULT 0.00,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_orders
        FOREIGN KEY (order_id)
            REFERENCES orders(id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE order_lines;

