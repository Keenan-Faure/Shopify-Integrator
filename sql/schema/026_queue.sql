-- +goose Up
ALTER TABLE queue_items ADD COLUMN description TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE DROP COLUMN description;