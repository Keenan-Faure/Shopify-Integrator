-- +goose Up
CREATE TABLE users(
    id UUID UNIQUE PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT (encode(sha256(random()::text::bytea), 'hex')),
    webhook_token VARCHAR(64) UNIQUE NOT NULL DEFAULT (encode(sha256(random()::text::bytea), 'hex')),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE users;