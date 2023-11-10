-- +goose Up
CREATE TABLE product_options(
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL,
    name VARCHAR(16) NOT NULL,
    position INT NOT NULL,
    CONSTRAINT fk_products
        FOREIGN KEY (product_id)
            REFERENCES products(id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE product_options;