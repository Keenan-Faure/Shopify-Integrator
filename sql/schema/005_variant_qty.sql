-- +goose Up
CREATE TABLE variant_qty(
    id UUID PRIMARY KEY UNIQUE,
    variant_id UUID NOT NULL,
    name VARCHAR(32) NOT NULL,
    value INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_variants
        FOREIGN KEY (variant_id)
            REFERENCES variants(id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE variant_qty;