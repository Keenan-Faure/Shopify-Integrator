-- +goose Up
CREATE TABLE fetch_stats(
    id UUID PRIMARY KEY NOT NULL UNIQUE,
    amount_of_products INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE fetch_stats;