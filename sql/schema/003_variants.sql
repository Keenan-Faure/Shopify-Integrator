-- +goose Up
CREATE TABLE variants(
    id BINARY(16) PRIMARY KEY UNIQUE NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    product_id BINARY(16) NOT NULL,
    sku VARCHAR(64) NOT NULL,
    option1 VARCHAR(16),
    option2 VARCHAR(16),
    option3 VARCHAR(16),
    barcode VARCHAR(64),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE variants;