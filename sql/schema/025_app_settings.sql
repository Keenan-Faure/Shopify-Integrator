-- +goose Up
CREATE TABLE app_settings(
    id UUID PRIMARY KEY UNIQUE,
    key VARCHAR(64) UNIQUE NOT NULL,
    description VARCHAR(255) NOT NULL DEFAULT '',
    value TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE app_settings;