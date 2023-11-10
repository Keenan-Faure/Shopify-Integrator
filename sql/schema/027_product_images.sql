-- +goose Up
CREATE TABLE product_images(
    id UUID UNIQUE PRIMARY KEY,
    product_id UUID NOT NULL,
    image_url TEXT NOT NULL,
    position INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_product
        FOREIGN KEY (product_id)
            REFERENCES products(id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE product_images;