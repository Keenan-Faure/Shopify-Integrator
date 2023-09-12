-- name: CreateOrderLine :one
INSERT INTO order_lines(
    order_id,
    line_type,
    sku,
    price,
    barcode,
    qty,
    tax_rate,
    tax_total,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: UpdateOrderLine :one
UPDATE order_lines
SET
    order_id = $1,
    line_type = $2,
    sku = $3,
    price = $4,
    barcode = $5,
    qty = $6,
    tax_rate = $7,
    tax_total = $8,
    created_at = $9,
    updated_at = $10
WHERE id = $11
RETURNING *;

-- name: GetShippingLinesByOrder :many
SELECT
    sku,
    line_type,
    price,
    barcode,
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
    barcode,
    qty,
    tax_rate,
    tax_total,
    updated_at
FROM order_lines
WHERE order_id = $1 AND line_type = 'line';
