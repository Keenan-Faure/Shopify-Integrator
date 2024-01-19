-- +goose Up
CREATE TABLE google_oauth(
    id UUID PRIMARY KEY UNIQUE NOT NULL,
    google_id VARCHAR(32) NOT NULL,
    email VARCHAR(32) NOT NULL,
    picture TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE google_oauth;