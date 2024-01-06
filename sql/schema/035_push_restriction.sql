-- +goose Up
CREATE TABLE push_restriction(
    id UUID PRIMARY KEY UNIQUE NOT NULL,
    field VARCHAR(64) NOT NULL UNIQUE,
    flag VARCHAR(32) NOT NULL, 
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

INSERT INTO push_restriction(
    id,
    field,
    flag,
    created_at,
    updated_at
) VALUES
(
    uuid_generate_v4(),
    'title',
    'app',
    NOW(),
    NOW()
),
(
    uuid_generate_v4(),
    'body_html',
    'app',
    NOW(),
    NOW()
),
(
    uuid_generate_v4(),
    'category',
    'app',
    NOW(),
    NOW()
),
(
    uuid_generate_v4(),
    'vendor',
    'app',
    NOW(),
    NOW()
),
(
    uuid_generate_v4(),
    'product_type',
    'app',
    NOW(),
    NOW()
),
(
    uuid_generate_v4(),
    'barcode',
    'app',
    NOW(),
    NOW()
),
(
    uuid_generate_v4(),
    'options',
    'app',
    NOW(),
    NOW()
),
(
    uuid_generate_v4(),
    'pricing',
    'app',
    NOW(),
    NOW()
),
(
    uuid_generate_v4(),
    'warehousing',
    'app',
    NOW(),
    NOW()
);

-- +goose Down
DROP TABLE push_restriction;