-- +goose Up
CREATE TABLE products(
    id VARCHAR(32) PRIMARY KEY NOT NULL,
    active VARCHAR(1) NOT NULL,
    title VARCHAR(255),
    body_html TEXT,
    category VARCHAR(255) NOT NULL,
    product_type VARCHAR(255),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE (id)
);

-- +goose Down
DROP TABLE products;