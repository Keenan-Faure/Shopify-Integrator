-- +goose Up
CREATE TABLE queue_items(
    id UUID PRIMARY KEY UNIQUE,
    type VARCHAR(32) NOT NULL,
    instruction VARCHAR(32) NOT NULL,
    status VARCHAR(16) NOT NULL,
    object JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE queue_items;