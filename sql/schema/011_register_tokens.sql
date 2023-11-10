-- +goose Up
CREATE TABLE register_tokens(
    id UUID UNIQUE PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    email VARCHAR(32) UNIQUE NOT NULL,
    token UUID UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE register_tokens;