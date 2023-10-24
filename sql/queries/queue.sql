-- name: CreateQueueItem :one
INSERT INTO queue_items(
    id,
    object,
    type,
    instruction,
    status,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id;

-- name: GetNextQueueItem :one
SELECT * FROM queue_items
ORDER BY created_at
LIMIT 1;

-- name: UpdateQueueItem :exec
UPDATE queue_items
SET
    status = $1,
    updated_at = $2
WHERE id = $3;

-- name: GetQueueItemByID :one
SELECT * FROM queue_items
WHERE ID = $1
LIMIT 1;

-- name: GetQueueItemsByDate :many
SELECT * FROM queue_items
WHERE status = $1
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;

-- name: GetQueueItemsByType :many
SELECT 
    id,
    object,
    type,
    instruction,
    status,
    updated_at
FROM queue_items
WHERE type = $1
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;

-- name: GetQueueItemsByStatus :many
SELECT 
    id,
    object,
    type,
    instruction,
    status,
    updated_at
FROM queue_items
WHERE status = $1
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;

-- name: GetQueueItemsByInstruction :many
SELECT 
    id,
    object,
    type,
    instruction,
    status,
    updated_at
FROM queue_items
WHERE instruction = $1
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;

-- name: GetQueueItemsByInstructionAndStatus :many
SELECT 
    id,
    object,
    type,
    instruction,
    status,
    updated_at
FROM queue_items
WHERE instruction = $1
AND status = $2
ORDER BY updated_at DESC
LIMIT $3 OFFSET $4;

-- name: GetQueueItemsByInstructionAndType :many
SELECT 
    id,
    object,
    type,
    instruction,
    status,
    updated_at
FROM queue_items
WHERE instruction = $1
AND type = $2
ORDER BY updated_at DESC
LIMIT $3 OFFSET $4;

-- name: GetQueueItemsByStatusAndType :many
SELECT 
    id,
    object,
    type,
    instruction,
    status,
    updated_at
FROM queue_items
WHERE status = $1
AND type = $2
ORDER BY updated_at DESC
LIMIT $3 OFFSET $4;

-- name: GetQueueItemsByFilter :many
SELECT 
    id,
    object,
    type,
    instruction,
    status,
    updated_at
FROM queue_items
WHERE status = $1
AND type = $2
AND instruction = $3
ORDER BY updated_at DESC
LIMIT $4 OFFSET $5;


-- name: RemoveQueueItemByID :exec
DELETE FROM queue_items
WHERE id = $1;

-- name: RemoveQueueItemsByStatus :exec
DELETE FROM queue_items
WHERE "status" IN ($1);

-- name: RemoveQueueItemsByInstruction :exec
DELETE FROM queue_items
WHERE instruction IN ($1);

