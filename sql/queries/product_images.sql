-- name: CreateProductImage :exec
INSERT INTO product_images(
    id,
    product_id,
    image_url,
    position,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: UpdateProductImage :exec
UPDATE product_images
SET
    image_url = $1,
    updated_at = $2
WHERE product_id = $3
AND position = $4;

-- name: GetProductImageByProductID :many
SELECT * FROM product_images
WHERE product_id = $1;

-- name: GetMaxImagePosition :one
SELECT 
    CAST(MAX("position") AS INTEGER)
FROM product_images;
