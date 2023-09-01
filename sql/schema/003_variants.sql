-- +goose Up
CREATE TABLE variants(
    id VARCHAR(32) PRIMARY KEY NOT NULL,
    product_id VARCHAR(32) NOT NULL,
    sku VARCHAR(255) NOT NULL,
    option1 VARCHAR(16) NOT NULL,
    option2 VARCHAR(16) NOT NULL,
    option3 VARCHAR(16) NOT NULL,
    barcode VARCHAR(255)
);

-- +goose Down
DROP TABLE variants;