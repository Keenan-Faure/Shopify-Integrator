-- +goose Up
CREATE TABLE shopify_vid(
    id UUID PRIMARY KEY,
    sku TEXT UNIQUE NOT NULL,
    variant_id UUID NOT NULL,
    shopify_variant_id VARCHAR(16) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_variants
        FOREIGN KEY (variant_id)
            REFERENCES variants(id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE shopify_vid;