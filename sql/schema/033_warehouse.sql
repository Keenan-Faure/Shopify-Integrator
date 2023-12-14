-- +goose Up
CREATE TABLE warehouses(
    id UUID PRIMARY KEY UNIQUE NOT NULL,
    name VARCHAR(64) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE warehouses;