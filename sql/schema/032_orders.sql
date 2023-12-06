-- +goose Up
ALTER TABLE orders ADD COLUMN status VARCHAR(32) NOT NULL DEFAULT 'not paid';

-- +goose Down
ALTER TABLE orders DROP COLUMN status;