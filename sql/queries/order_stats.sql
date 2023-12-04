-- name: CreateOrderStat :exec
INSERT INTO order_stats(
    id,
    order_total,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4
);

-- name: GetOrderStat :one
SELECT
    id,
    order_total,
    created_at,
    updated_at
FROM order_stats
WHERE id = $1;

-- name: GetOrderStats :many
SELECT 
    order_total
FROM order_stats
WHERE 
    "created_at" BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()
ORDER BY "created_at" DESC;
