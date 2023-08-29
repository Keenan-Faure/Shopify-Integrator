-- +goose Up
CREATE TABLE variants(
    id VARCHAR(32) PRIMARY KEY NOT NULL,
    product_id VARCHAR(32) NOT NULL,
    sku VARCHAR(255) NOT NULL,
    barcode VARCHAR(255)
);

-- +goose Down
DROP TABLE variants;