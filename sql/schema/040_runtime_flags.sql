-- +goose Up
CREATE TABLE runtime_flags(
    id UUID UNIQUE PRIMARY KEY,
    flag_name TEXT UNIQUE NOT NULL,
    flag_value BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

INSERT INTO runtime_flags(
    id,
    flag_name,
    flag_value,
    created_at,
    updated_at
) VALUES (
    uuid_generate_v4(),
    'workers',
    FALSE,
    NOW(),
    NOW()
),(
    uuid_generate_v4(),
    'localhost',
    TRUE,
    NOW(),
    NOW()
);

-- +goose Down
DROP TABLE runtime_flags;