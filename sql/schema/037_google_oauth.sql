-- +goose Up
CREATE TABLE google_oauth(
    id UUID PRIMARY KEY UNIQUE NOT NULL,
    user_id UUID UNIQUE NOT NULL,
    cookie_secret VARCHAR(64) UNIQUE NOT NULL DEFAULT (encode(sha256(random()::text::bytea), 'hex')),
    cookie_token VARCHAR(64) NOT NULL,
    google_id VARCHAR(32) NOT NULL UNIQUE,
    email VARCHAR(32) NOT NULL,
    picture TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_users
        FOREIGN KEY (user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE google_oauth;