-- +goose Up
CREATE TABLE users(
    id UUID PRIMARY KEY UNIQUE,
    webhook_token VARCHAR(32) UNIQUE NOT NULL DEFAULT (encode(sha256(random()::text::bytea), 'hex')),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    api_key VARCHAR(32) UNIQUE NOT NULL DEFAULT (encode(sha256(random()::text::bytea), 'hex'))
);

-- +goose Down
DROP TABLE users;