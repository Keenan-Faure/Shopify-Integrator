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
	SUM(amount_of_products) AS "amount",
	to_char(created_at, 'YYYY-MM-DD HH24:00') AS "hour"
FROM fetch_stats
WHERE created_at > current_date at time zone 'UTC'
GROUP BY "hour"
ORDER BY "hour" ASC;