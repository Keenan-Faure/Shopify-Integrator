-- +goose Up
CREATE TABLE register_tokens(
    id UUID PRIMARY KEY UNIQUE,
    name VARCHAR(64) NOT NULL,
    email VARCHAR(32) UNIQUE NOT NULL,
    token VARCHAR(6) UNIQUE NOT NULL DEFAULT (encode(sha256(random()::text::bytea), 'hex')),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE register_tokens;

