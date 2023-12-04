-- name: CreateOrderLine :one
INSERT INTO order_lines(
    id,
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
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING *;

-- name: UpdateOrderLine :exec
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
    updated_at = $9
WHERE id = $10;

-- name: UpdateOrderLineByOrderAndSKU :exec
UPDATE order_lines
SET
    line_type = $1,
    sku = $2,
    price = $3,
    barcode = $4,
    qty = $5,
    tax_rate = $6,
    tax_total = $7,
    updated_at = $8
WHERE order_id = $9
AND sku = $10;

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
WHERE order_id = $1 AND line_type = 'product';
