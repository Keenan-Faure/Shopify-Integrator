-- +goose Up
CREATE TABLE variants(
    id VARCHAR(32) PRIMARY KEY,
    product_id VARCHAR(32) NOT NULL,
    sku VARCHAR(255) NOT NULL,
    barcode VARCHAR(255),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

-- +goose Down
DROP TABLE variants;