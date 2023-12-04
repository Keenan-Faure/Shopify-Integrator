-- name: CreateFetchStat :exec
INSERT INTO fetch_stats(
    id,
    amount_of_products,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4
);

-- name: GetFetchStat :one
SELECT
    id,
    amount_of_products,
    created_at,
    updated_at
FROM fetch_stats
WHERE id = $1;

-- name: GetFetchStats :many
SELECT 
    amount_of_products
FROM fetch_stats
WHERE 
    "created_at" BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()
ORDER BY "created_at" DESC;
