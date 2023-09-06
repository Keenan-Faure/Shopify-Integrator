-- +goose Up
CREATE TABLE variant_qty(
    id BINARY(16) PRIMARY KEY UNIQUE NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    variant_id BINARY(16) UNIQUE NOT NULL,
    name VARCHAR(16) NOT NULL,
    value INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE variant_qty;