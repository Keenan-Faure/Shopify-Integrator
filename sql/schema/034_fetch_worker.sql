-- +goose Up
create table fetch_worker(
    id UUID PRIMARY KEY NOT NULL,
    status VARCHAR(1) NOT NULL,
    fetch_url VARCHAR(255) NOT NULL,
    local_count INTEGER NOT NULL,
    shopify_product_count INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

INSERT INTO fetch_worker(
    id,
    status,
    fetch_url,
    local_count,
    shopify_product_count,
    created_at,
    updated_at
) VALUES (
    uuid_generate_v4(),
    '0',
    '',
    0,
    0,
    NOW(),
    NOW()
);

-- +goose Down
drop table fetch_worker;