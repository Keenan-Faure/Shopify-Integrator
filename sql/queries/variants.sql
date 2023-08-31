-- name: CreateVariant :execresult
INSERT INTO variants(
    id,
    product_id,
    sku,
    barcode
) VALUES (?, ?, ?, ?);

-- name: UpdateVariant
UPDATE variants SET
sku = ?
price = ?
compare_at_price = ?
option1 = ?
option2 = ?
option3 = ?
barcode = ?;