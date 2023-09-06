-- +goose Up
CREATE TABLE variant_pricing(
    id BINARY(16) PRIMARY KEY UNIQUE NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    variant_id BINARY(16) UNIQUE NOT NULL,
    name VARCHAR(16) NOT NULL,
    value DECIMAL(9, 2) DEFAULT 0.00,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE variant_pricing;