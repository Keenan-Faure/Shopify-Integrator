-- +goose Up
CREATE TABLE variants(
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL,
    sku VARCHAR(64) NOT NULL,
    option1 VARCHAR(16),
    option2 VARCHAR(16),
    option3 VARCHAR(16),
    barcode VARCHAR(64),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_products
        FOREIGN KEY (product_id)
            REFERENCES products(id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE variants;