-- +goose Up
CREATE TABLE customers(
    id UUID UNIQUE PRIMARY KEY,
    web_customer_code VARCHAR(32) UNIQUE NOT NULL,
    first_name VARCHAR(32) NOT NULL,
    last_name VARCHAR(32) NOT NULL,
    email VARCHAR(32),
    phone VARCHAR(32),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE customers;