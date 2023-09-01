-- +goose Up
CREATE TABLE users(
    id VARCHAR(32) PRIMARY KEY NOT NULL,
    webhook_token VARCHAR(32) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    api_key VARCHAR(64) UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;