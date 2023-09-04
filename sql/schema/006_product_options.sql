-- +goose Up
CREATE TABLE product_options(
    id BINARY(16) PRIMARY KEY UNIQUE NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    product_id BINARY(16) UNIQUE NOT NULL,
    name VARCHAR(16) NOT NULL,
    value VARCHAR(32) NOT NULL
);

-- +goose Down
DROP TABLE product_options;