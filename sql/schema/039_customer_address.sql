-- +goose Up
CREATE TABLE customer_address(
    id UUID UNIQUE PRIMARY KEY,
    customer_id UUID NOT NULL,
    address_type VARCHAR(32) NOT NULL,
    address_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_address
        FOREIGN KEY (address_id)
            REFERENCES address(id)
            ON DELETE CASCADE,
    CONSTRAINT fk_customers
        FOREIGN KEY (customer_id)
            REFERENCES customers(id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE customer_address;