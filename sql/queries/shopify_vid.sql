-- name: CreateVID :exec
INSERT INTO shopify_vid(
    id,
    sku,
    shopify_variant_id,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: UpdateVID :exec
UPDATE shopify_vid
SET
    shopify_variant_id = $1,
    updated_at = $2
WHERE sku = $3;

-- name: GetVIDBySKU :one
SELECT
    sku,
    shopify_variant_id,
    updated_at
FROM shopify_vid
WHERE sku = $1
LIMIT 1;