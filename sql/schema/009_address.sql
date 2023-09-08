-- +goose Up
CREATE TABLE address(
    id UUID PRIMARY KEY UNIQUE,
    customer_id UUID NOT NULL,
    name VARCHAR(16),
    first_name VARCHAR(32) NOT NULL,
    last_name VARCHAR(32) NOT NULL,
    address1 VARCHAR(64) DEFAULT '',
    address2 VARCHAR(64) DEFAULT '',
    suburb VARCHAR(64) DEFAULT '',
    city VARCHAR(64) DEFAULT '',
    province VARCHAR(64) DEFAULT '',
    postal_code VARCHAR(64) DEFAULT '',
    company VARCHAR(64) DEFAULT '',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_customers
        FOREIGN KEY (customer_id)
            REFERENCES customers(id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE address;