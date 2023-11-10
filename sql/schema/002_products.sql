-- +goose Up
CREATE TABLE products(
    id UUID UNIQUE PRIMARY KEY,
    active VARCHAR(1) NOT NULL,
    product_code TEXT UNIQUE NOT NULL,
    title VARCHAR(255),
    body_html TEXT,
    category VARCHAR(64),
    vendor VARCHAR(64),
    product_type VARCHAR(64),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE products;