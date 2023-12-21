-- name: CreateFetchWorker :exec
INSERT INTO fetch_worker(
    id,
    status,
    local_count,
    shopify_product_count,
    fetch_url,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7 
);

-- name: GetFetchWorker :one
SELECT * FROM fetch_worker
LIMIT 1;

-- name: UpdateFetchWorker :exec
UPDATE fetch_worker
SET
    status = $1,
    local_count = $2,
    shopify_product_count = $3,
    fetch_url = $4,
    updated_at = $5
WHERE id = $6;
