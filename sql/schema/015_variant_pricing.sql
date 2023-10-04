-- +goose Up
ALTER TABLE variant_pricing ADD COLUMN isdefault BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE variant_pricing DROP COLUMN isdefault;