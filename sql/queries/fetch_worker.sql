-- name: CreateFetchWorker :exec
INSERT INTO fetch_worker(
    id,
    status,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4
);

-- name: GetFetchWorker :one
SELECT * FROM fetch_worker
LIMIT 1;

-- name: UpdateFetchWorker :exec
UPDATE fetch_worker
SET
    status = $1,
    updated_at = $2
WHERE id = $3;
