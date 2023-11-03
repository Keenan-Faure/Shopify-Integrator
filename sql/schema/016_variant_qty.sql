-- +goose Up
ALTER TABLE variant_qty ADD COLUMN isdefault BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE variant_qty DROP COLUMN isdefault;