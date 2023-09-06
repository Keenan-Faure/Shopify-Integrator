-- +goose Up
CREATE TABLE address(
    id BINARY(16) PRIMARY KEY UNIQUE NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    customer_id BINARY(16) NOT NULL,
    first_name VARCHAR(32) NOT NULL,
    last_name VARCHAR(32) NOT NULL,
    address1 VARCHAR(64) DEFAULT '',
    address2 VARCHAR(64) DEFAULT '',
    suburb VARCHAR(64) DEFAULT '',
    city VARCHAR(64) DEFAULT '',
    province VARCHAR(64) DEFAULT '',
    postal_code VARCHAR(64) DEFAULT '',
    company VARCHAR(64) DEFAULT '',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE address;