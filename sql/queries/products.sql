-- name: CreateProduct :execresult
INSERT INTO products(
    id,
    active,
    title,
    body_html,
    category,
    product_type,
    created_at,
    updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateProduct :execresult
UPDATE products SET
active = ?
title = ?
body_html = ?
category = ?
product_type = ?
updated_at = ?;

