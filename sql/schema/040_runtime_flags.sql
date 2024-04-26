-- +goose Up
CREATE TABLE runtime_flags(
    id UUID UNIQUE PRIMARY KEY,
    flag_name TEXT UNIQUE NOT NULL,
    flag_value BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE runtime_flags;