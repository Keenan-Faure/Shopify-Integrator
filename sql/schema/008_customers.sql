-- +goose Up
CREATE TABLE customers(
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL,
    first_name VARCHAR(32) NOT NULL,
    last_name VARCHAR(32) NOT NULL,
    email VARCHAR(32),
    phone VARCHAR(32),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_orders
        FOREIGN KEY (order_id)
            REFERENCES orders(id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE customers;