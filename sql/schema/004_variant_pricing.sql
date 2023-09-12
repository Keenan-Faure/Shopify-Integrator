-- +goose Up
CREATE TABLE variant_pricing(
    id UUID PRIMARY KEY,
    variant_id UUID NOT NULL,
    name VARCHAR(16) NOT NULL,
    value DECIMAL(9, 2) DEFAULT 0.00,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_variants
        FOREIGN KEY (variant_id)
            REFERENCES variants(id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE variant_pricing;