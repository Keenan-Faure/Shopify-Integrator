-- +goose Up
CREATE TABLE variant_pricing(
    id VARCHAR(32) UNIQUE NOT NULL,
    variant_id VARCHAR(32) UNIQUE NOT NULL,
    name VARCHAR(16) NOT NULL,
    value VARCHAR(32) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE variant_pricing;