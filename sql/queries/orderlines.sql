-- name: CreateOrderLine :execresult
INSERT INTO order_lines(
    order_id,
    sku,
    price,
    barcode,
    qty,
    tax_rate,
    tax_total,
    created_at,
    updated_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateOrderLine :execresult
UPDATE order_lines
SET
    order_id = ?,
    sku = ?,
    price = ?,
    barcode = ?,
    qty = ?,
    tax_rate = ?,
    tax_total = ?,
    created_at = ?,
    updated_at = ?
WHERE id = ?;

-- name: GetOrderLinesByOrder :many
SELECT
    sku,
    price,
    barcode,
    qty,
    tax_rate,
    tax_total,
    updated_at
FROM order_lines
WHERE order_id = ?;
