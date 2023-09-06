-- name: CreateOrderLine :execresult
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
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateOrderLine :execresult
UPDATE order_lines
SET
    order_id = ?,
    line_type = ?,
    sku = ?,
    price = ?,
    barcode = ?,
    qty = ?,
    tax_rate = ?,
    tax_total = ?,
    created_at = ?,
    updated_at = ?
WHERE id = ?;

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
WHERE order_id = ? AND line_type = 'shipping';

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
WHERE order_id = ? AND line_type = 'line';
