-- +goose Up
CREATE TABLE customer_orders(
    id UUID UNIQUE PRIMARY KEY,
    customer_id UUID NOT NULL,
    order_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_orders
        FOREIGN KEY (order_id)
            REFERENCES orders(id)
            ON DELETE CASCADE,
    CONSTRAINT fk_customers
        FOREIGN KEY (customer_id)
            REFERENCES customers(id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE customer_orders;