-- name: CreateQueueItem :one
INSERT INTO queue_items(
    id,
    queue_type,
    instruction,
    "status",
    "object",
    "description",
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id;

-- name: UpdateQueueItem :one
UPDATE queue_items
SET
    "status" = $1,
    updated_at = $2,
    "description" = $3
WHERE id = $4
RETURNING *;

-- name: GetNextQueueItem :one
SELECT * FROM queue_items
WHERE "status" NOT IN ('completed', 'failed')
ORDER BY instruction asc, created_at desc
LIMIT 1;

-- name: GetQueueSize :one
SELECT COUNT(*) FROM queue_items;

-- name: GetQueueItemByID :one
SELECT * FROM queue_items
WHERE ID = $1
LIMIT 1;

-- name: GetQueueItemsByDate :many
SELECT * FROM queue_items
WHERE "status" = $1
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;

-- name: GetQueueItemsByType :many
SELECT 
    id,
    "object",
    queue_type,
    instruction,
    "status",
    "description",
    updated_at
FROM queue_items
WHERE queue_type = $1
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;

-- name: GetQueueItemsByStatus :many
SELECT 
    id,
    "object",
    queue_type,
    instruction,
    "status",
    "description",
    updated_at
FROM queue_items
WHERE "status" = $1
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;

-- name: GetQueueItemsByInstruction :many
SELECT 
    id,
    "object",
    queue_type,
    instruction,
    "status",
    "description",
    updated_at
FROM queue_items
WHERE instruction = $1
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;

-- name: GetQueueItemsByInstructionAndStatus :many
SELECT 
    id,
    "object",
    queue_type,
    instruction,
    "status",
    "description",
    updated_at
FROM queue_items
WHERE instruction = $1
AND "status" = $2
ORDER BY updated_at DESC
LIMIT $3 OFFSET $4;

-- name: GetQueueItemsByInstructionAndType :many
SELECT 
    id,
    "object",
    queue_type,
    instruction,
    "status",
    "description",
    updated_at
FROM queue_items
WHERE instruction = $1
AND queue_type = $2
ORDER BY updated_at DESC
LIMIT $3 OFFSET $4;

-- name: GetQueueItemsByStatusAndType :many
SELECT 
    id,
    "object",
    queue_type,
    instruction,
    "status",
    "description",
    updated_at
FROM queue_items
WHERE "status" = $1
AND queue_type = $2
ORDER BY updated_at DESC
LIMIT $3 OFFSET $4;

-- name: GetQueueItemsByFilter :many
SELECT 
    id,
    "object",
    queue_type,
    instruction,
    "status",
    "description",
    updated_at
FROM queue_items
WHERE "status" = $1
AND queue_type = $2
AND instruction = $3
ORDER BY updated_at DESC
LIMIT $4 OFFSET $5;

-- name: GetQueueItemsCount :one
SELECT
    COUNT(*)
FROM queue_items
WHERE instruction = $1 AND
status NOT IN ('completed', 'failed');

-- name: RemoveQueueItemByID :exec
DELETE FROM queue_items
WHERE id = $1;

-- name: RemoveQueueItemsByStatus :exec
DELETE FROM queue_items WHERE
"status" = $1;

-- name: RemoveQueueItemsByInstruction :exec
DELETE FROM queue_items WHERE
instruction = $1;

-- name: RemoveQueueItemsByType :exec
DELETE FROM queue_items WHERE
queue_type = $1;

-- name: RemoveQueueItemsByStatusAndInstruction :exec
DELETE FROM queue_items WHERE
"status" = $1 AND
instruction = $2;

-- name: RemoveQueueItemsByTypeAndInstruction :exec
DELETE FROM queue_items WHERE
queue_type = $1 AND
instruction = $2;

-- name: RemoveQueueItemsByStatusAndType :exec
DELETE FROM queue_items WHERE
"status" = $1 AND
queue_type = $2;

-- name: RemoveQueueItemsFilter :exec
DELETE FROM queue_items WHERE
"status" = $1 AND
queue_type = $2 AND
instruction = $3;
