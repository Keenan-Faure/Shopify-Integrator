-- +goose Up
CREATE TABLE product_options(
    id VARCHAR(32) UNIQUE NOT NULL,
    product_id VARCHAR(32) UNIQUE NOT NULL,
    name VARCHAR(16) NOT NULL,
    value VARCHAR(32) NOT NULL
);

-- +goose Down
DROP TABLE product_options;