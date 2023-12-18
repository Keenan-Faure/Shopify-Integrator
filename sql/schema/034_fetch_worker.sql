-- +goose Up
create table fetch_worker(
    id UUID PRIMARY KEY NOT NULL,
    status VARCHAR(1) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

INSERT INTO fetch_worker(
    id,
    status,
    created_at,
    updated_at
) VALUES (
    uuid_generate_v4(),
    '0',
    NOW(),
    NOW()
);

-- +goose Down
drop table fetch_worker;