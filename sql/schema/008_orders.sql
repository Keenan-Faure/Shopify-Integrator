-- +goose Up
CREATE TABLE orders(
    id UUID PRIMARY KEY UNIQUE,
    customer_id UUID NOT NULL,
    notes VARCHAR(255) DEFAULT '',
    web_code VARCHAR(32) UNIQUE,
    tax_total DECIMAL(10, 2) DEFAULT 0.00,
    order_total DECIMAL(10, 2) DEFAULT 0.00,
    shipping_total DECIMAL(10, 2) DEFAULT 0.00,
    discount_total DECIMAL(10, 2) DEFAULT 0.00,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_customers
        FOREIGN KEY (customer_id)
            REFERENCES customers(id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE orders;