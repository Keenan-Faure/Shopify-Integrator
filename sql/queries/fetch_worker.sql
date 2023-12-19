-- name: CreateFetchWorker :exec
INSERT INTO fetch_worker(
    id,
    status,
    fetch_url,
    local_count,
    shopify_product_count,
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
    fetch_url = $2,
    local_count = $3,
    shopify_product_count = $4,
    updated_at = $5
WHERE id = $6;
