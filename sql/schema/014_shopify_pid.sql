-- +goose Up
CREATE TABLE shopify_pid(
    id UUID PRIMARY KEY,
    product_code TEXT UNIQUE NOT NULL,
    product_id UUID NOT NULL,
    shopify_product_id VARCHAR(16) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_products
        FOREIGN KEY (product_id)
            REFERENCES products(id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE shopify_pid;