-- +goose Up
CREATE TABLE order_stats(
    id UUID PRIMARY KEY NOT NULL UNIQUE,
    order_total INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE order_stats;