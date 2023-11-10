-- +goose Up
CREATE TABLE variants(
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL,
    sku TEXT UNIQUE NOT NULL,
    option1 VARCHAR(32),
    option2 VARCHAR(32),
    option3 VARCHAR(32),
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