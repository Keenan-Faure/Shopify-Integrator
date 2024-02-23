-- name: CreateOrderLine :one
INSERT INTO order_lines(
    id,
    order_id,
    line_type,
    sku,
    price,
    qty,
    tax_rate,
    tax_total,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: UpdateOrderLine :exec
UPDATE order_lines
SET
    order_id = $1,
    line_type = $2,
    sku = $3,
    price = $4,
    qty = $5,
    tax_rate = $6,
    tax_total = $7,
    updated_at = $8
WHERE id = $9;

-- name: UpdateOrderLineByOrderAndSKU :exec
UPDATE order_lines
SET
    line_type = $1,
    sku = $2,
    price = $3,
    qty = $4,
    tax_rate = $5,
    tax_total = $6,
    updated_at = $7
WHERE order_id = $8
AND sku = $9;

-- name: GetShippingLinesByOrder :many
SELECT
    sku,
    line_type,
    price,
    qty,
    tax_rate,
    tax_total,
    updated_at
FROM order_lines
WHERE order_id = $1 AND line_type = 'shipping';

-- name: GetOrderLinesByOrder :many
SELECT
    sku,
    line_type,
    price,
    qty,
    tax_rate,
    tax_total,
    updated_at
FROM order_lines
WHERE order_id = $1 AND line_type = 'product';
