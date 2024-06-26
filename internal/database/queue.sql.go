// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queue.sql

package database

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const createQueueItem = `-- name: CreateQueueItem :one
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
RETURNING id
`

type CreateQueueItemParams struct {
	ID          uuid.UUID       `json:"id"`
	QueueType   string          `json:"queue_type"`
	Instruction string          `json:"instruction"`
	Status      string          `json:"status"`
	Object      json.RawMessage `json:"object"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (q *Queries) CreateQueueItem(ctx context.Context, arg CreateQueueItemParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createQueueItem,
		arg.ID,
		arg.QueueType,
		arg.Instruction,
		arg.Status,
		arg.Object,
		arg.Description,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getNextQueueItem = `-- name: GetNextQueueItem :one
SELECT id, queue_type, instruction, status, object, created_at, updated_at, description FROM queue_items
WHERE "status" NOT IN ('completed', 'failed')
ORDER BY instruction asc, created_at desc
LIMIT 1
`

func (q *Queries) GetNextQueueItem(ctx context.Context) (QueueItem, error) {
	row := q.db.QueryRowContext(ctx, getNextQueueItem)
	var i QueueItem
	err := row.Scan(
		&i.ID,
		&i.QueueType,
		&i.Instruction,
		&i.Status,
		&i.Object,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Description,
	)
	return i, err
}

const getQueueItemByID = `-- name: GetQueueItemByID :one
SELECT id, queue_type, instruction, status, object, created_at, updated_at, description FROM queue_items
WHERE ID = $1
LIMIT 1
`

func (q *Queries) GetQueueItemByID(ctx context.Context, id uuid.UUID) (QueueItem, error) {
	row := q.db.QueryRowContext(ctx, getQueueItemByID, id)
	var i QueueItem
	err := row.Scan(
		&i.ID,
		&i.QueueType,
		&i.Instruction,
		&i.Status,
		&i.Object,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Description,
	)
	return i, err
}

const getQueueItemsByDate = `-- name: GetQueueItemsByDate :many
SELECT id, queue_type, instruction, status, object, created_at, updated_at, description FROM queue_items
WHERE "status" = $1
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3
`

type GetQueueItemsByDateParams struct {
	Status string `json:"status"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) GetQueueItemsByDate(ctx context.Context, arg GetQueueItemsByDateParams) ([]QueueItem, error) {
	rows, err := q.db.QueryContext(ctx, getQueueItemsByDate, arg.Status, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []QueueItem
	for rows.Next() {
		var i QueueItem
		if err := rows.Scan(
			&i.ID,
			&i.QueueType,
			&i.Instruction,
			&i.Status,
			&i.Object,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getQueueItemsByFilter = `-- name: GetQueueItemsByFilter :many
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
LIMIT $4 OFFSET $5
`

type GetQueueItemsByFilterParams struct {
	Status      string `json:"status"`
	QueueType   string `json:"queue_type"`
	Instruction string `json:"instruction"`
	Limit       int32  `json:"limit"`
	Offset      int32  `json:"offset"`
}

type GetQueueItemsByFilterRow struct {
	ID          uuid.UUID       `json:"id"`
	Object      json.RawMessage `json:"object"`
	QueueType   string          `json:"queue_type"`
	Instruction string          `json:"instruction"`
	Status      string          `json:"status"`
	Description string          `json:"description"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (q *Queries) GetQueueItemsByFilter(ctx context.Context, arg GetQueueItemsByFilterParams) ([]GetQueueItemsByFilterRow, error) {
	rows, err := q.db.QueryContext(ctx, getQueueItemsByFilter,
		arg.Status,
		arg.QueueType,
		arg.Instruction,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetQueueItemsByFilterRow
	for rows.Next() {
		var i GetQueueItemsByFilterRow
		if err := rows.Scan(
			&i.ID,
			&i.Object,
			&i.QueueType,
			&i.Instruction,
			&i.Status,
			&i.Description,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getQueueItemsByInstruction = `-- name: GetQueueItemsByInstruction :many
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
LIMIT $2 OFFSET $3
`

type GetQueueItemsByInstructionParams struct {
	Instruction string `json:"instruction"`
	Limit       int32  `json:"limit"`
	Offset      int32  `json:"offset"`
}

type GetQueueItemsByInstructionRow struct {
	ID          uuid.UUID       `json:"id"`
	Object      json.RawMessage `json:"object"`
	QueueType   string          `json:"queue_type"`
	Instruction string          `json:"instruction"`
	Status      string          `json:"status"`
	Description string          `json:"description"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (q *Queries) GetQueueItemsByInstruction(ctx context.Context, arg GetQueueItemsByInstructionParams) ([]GetQueueItemsByInstructionRow, error) {
	rows, err := q.db.QueryContext(ctx, getQueueItemsByInstruction, arg.Instruction, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetQueueItemsByInstructionRow
	for rows.Next() {
		var i GetQueueItemsByInstructionRow
		if err := rows.Scan(
			&i.ID,
			&i.Object,
			&i.QueueType,
			&i.Instruction,
			&i.Status,
			&i.Description,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getQueueItemsByInstructionAndStatus = `-- name: GetQueueItemsByInstructionAndStatus :many
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
LIMIT $3 OFFSET $4
`

type GetQueueItemsByInstructionAndStatusParams struct {
	Instruction string `json:"instruction"`
	Status      string `json:"status"`
	Limit       int32  `json:"limit"`
	Offset      int32  `json:"offset"`
}

type GetQueueItemsByInstructionAndStatusRow struct {
	ID          uuid.UUID       `json:"id"`
	Object      json.RawMessage `json:"object"`
	QueueType   string          `json:"queue_type"`
	Instruction string          `json:"instruction"`
	Status      string          `json:"status"`
	Description string          `json:"description"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (q *Queries) GetQueueItemsByInstructionAndStatus(ctx context.Context, arg GetQueueItemsByInstructionAndStatusParams) ([]GetQueueItemsByInstructionAndStatusRow, error) {
	rows, err := q.db.QueryContext(ctx, getQueueItemsByInstructionAndStatus,
		arg.Instruction,
		arg.Status,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetQueueItemsByInstructionAndStatusRow
	for rows.Next() {
		var i GetQueueItemsByInstructionAndStatusRow
		if err := rows.Scan(
			&i.ID,
			&i.Object,
			&i.QueueType,
			&i.Instruction,
			&i.Status,
			&i.Description,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getQueueItemsByInstructionAndType = `-- name: GetQueueItemsByInstructionAndType :many
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
LIMIT $3 OFFSET $4
`

type GetQueueItemsByInstructionAndTypeParams struct {
	Instruction string `json:"instruction"`
	QueueType   string `json:"queue_type"`
	Limit       int32  `json:"limit"`
	Offset      int32  `json:"offset"`
}

type GetQueueItemsByInstructionAndTypeRow struct {
	ID          uuid.UUID       `json:"id"`
	Object      json.RawMessage `json:"object"`
	QueueType   string          `json:"queue_type"`
	Instruction string          `json:"instruction"`
	Status      string          `json:"status"`
	Description string          `json:"description"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (q *Queries) GetQueueItemsByInstructionAndType(ctx context.Context, arg GetQueueItemsByInstructionAndTypeParams) ([]GetQueueItemsByInstructionAndTypeRow, error) {
	rows, err := q.db.QueryContext(ctx, getQueueItemsByInstructionAndType,
		arg.Instruction,
		arg.QueueType,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetQueueItemsByInstructionAndTypeRow
	for rows.Next() {
		var i GetQueueItemsByInstructionAndTypeRow
		if err := rows.Scan(
			&i.ID,
			&i.Object,
			&i.QueueType,
			&i.Instruction,
			&i.Status,
			&i.Description,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getQueueItemsByStatus = `-- name: GetQueueItemsByStatus :many
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
LIMIT $2 OFFSET $3
`

type GetQueueItemsByStatusParams struct {
	Status string `json:"status"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

type GetQueueItemsByStatusRow struct {
	ID          uuid.UUID       `json:"id"`
	Object      json.RawMessage `json:"object"`
	QueueType   string          `json:"queue_type"`
	Instruction string          `json:"instruction"`
	Status      string          `json:"status"`
	Description string          `json:"description"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (q *Queries) GetQueueItemsByStatus(ctx context.Context, arg GetQueueItemsByStatusParams) ([]GetQueueItemsByStatusRow, error) {
	rows, err := q.db.QueryContext(ctx, getQueueItemsByStatus, arg.Status, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetQueueItemsByStatusRow
	for rows.Next() {
		var i GetQueueItemsByStatusRow
		if err := rows.Scan(
			&i.ID,
			&i.Object,
			&i.QueueType,
			&i.Instruction,
			&i.Status,
			&i.Description,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getQueueItemsByStatusAndType = `-- name: GetQueueItemsByStatusAndType :many
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
LIMIT $3 OFFSET $4
`

type GetQueueItemsByStatusAndTypeParams struct {
	Status    string `json:"status"`
	QueueType string `json:"queue_type"`
	Limit     int32  `json:"limit"`
	Offset    int32  `json:"offset"`
}

type GetQueueItemsByStatusAndTypeRow struct {
	ID          uuid.UUID       `json:"id"`
	Object      json.RawMessage `json:"object"`
	QueueType   string          `json:"queue_type"`
	Instruction string          `json:"instruction"`
	Status      string          `json:"status"`
	Description string          `json:"description"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (q *Queries) GetQueueItemsByStatusAndType(ctx context.Context, arg GetQueueItemsByStatusAndTypeParams) ([]GetQueueItemsByStatusAndTypeRow, error) {
	rows, err := q.db.QueryContext(ctx, getQueueItemsByStatusAndType,
		arg.Status,
		arg.QueueType,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetQueueItemsByStatusAndTypeRow
	for rows.Next() {
		var i GetQueueItemsByStatusAndTypeRow
		if err := rows.Scan(
			&i.ID,
			&i.Object,
			&i.QueueType,
			&i.Instruction,
			&i.Status,
			&i.Description,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getQueueItemsByType = `-- name: GetQueueItemsByType :many
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
LIMIT $2 OFFSET $3
`

type GetQueueItemsByTypeParams struct {
	QueueType string `json:"queue_type"`
	Limit     int32  `json:"limit"`
	Offset    int32  `json:"offset"`
}

type GetQueueItemsByTypeRow struct {
	ID          uuid.UUID       `json:"id"`
	Object      json.RawMessage `json:"object"`
	QueueType   string          `json:"queue_type"`
	Instruction string          `json:"instruction"`
	Status      string          `json:"status"`
	Description string          `json:"description"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (q *Queries) GetQueueItemsByType(ctx context.Context, arg GetQueueItemsByTypeParams) ([]GetQueueItemsByTypeRow, error) {
	rows, err := q.db.QueryContext(ctx, getQueueItemsByType, arg.QueueType, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetQueueItemsByTypeRow
	for rows.Next() {
		var i GetQueueItemsByTypeRow
		if err := rows.Scan(
			&i.ID,
			&i.Object,
			&i.QueueType,
			&i.Instruction,
			&i.Status,
			&i.Description,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getQueueItemsCount = `-- name: GetQueueItemsCount :one
SELECT
    COUNT(*)
FROM queue_items
WHERE instruction = $1 AND
status NOT IN ('completed', 'failed')
`

func (q *Queries) GetQueueItemsCount(ctx context.Context, instruction string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getQueueItemsCount, instruction)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getQueueSize = `-- name: GetQueueSize :one
SELECT COUNT(*) FROM queue_items
WHERE "status" IN ('in-queue', 'processing')
`

func (q *Queries) GetQueueSize(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, getQueueSize)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const removeQueueItemByID = `-- name: RemoveQueueItemByID :exec
DELETE FROM queue_items
WHERE id = $1
`

func (q *Queries) RemoveQueueItemByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, removeQueueItemByID, id)
	return err
}

const removeQueueItemsByInstruction = `-- name: RemoveQueueItemsByInstruction :exec
DELETE FROM queue_items WHERE
instruction = $1
`

func (q *Queries) RemoveQueueItemsByInstruction(ctx context.Context, instruction string) error {
	_, err := q.db.ExecContext(ctx, removeQueueItemsByInstruction, instruction)
	return err
}

const removeQueueItemsByStatus = `-- name: RemoveQueueItemsByStatus :exec
DELETE FROM queue_items WHERE
"status" = $1
`

func (q *Queries) RemoveQueueItemsByStatus(ctx context.Context, status string) error {
	_, err := q.db.ExecContext(ctx, removeQueueItemsByStatus, status)
	return err
}

const removeQueueItemsByStatusAndInstruction = `-- name: RemoveQueueItemsByStatusAndInstruction :exec
DELETE FROM queue_items WHERE
"status" = $1 AND
instruction = $2
`

type RemoveQueueItemsByStatusAndInstructionParams struct {
	Status      string `json:"status"`
	Instruction string `json:"instruction"`
}

func (q *Queries) RemoveQueueItemsByStatusAndInstruction(ctx context.Context, arg RemoveQueueItemsByStatusAndInstructionParams) error {
	_, err := q.db.ExecContext(ctx, removeQueueItemsByStatusAndInstruction, arg.Status, arg.Instruction)
	return err
}

const removeQueueItemsByStatusAndType = `-- name: RemoveQueueItemsByStatusAndType :exec
DELETE FROM queue_items WHERE
"status" = $1 AND
queue_type = $2
`

type RemoveQueueItemsByStatusAndTypeParams struct {
	Status    string `json:"status"`
	QueueType string `json:"queue_type"`
}

func (q *Queries) RemoveQueueItemsByStatusAndType(ctx context.Context, arg RemoveQueueItemsByStatusAndTypeParams) error {
	_, err := q.db.ExecContext(ctx, removeQueueItemsByStatusAndType, arg.Status, arg.QueueType)
	return err
}

const removeQueueItemsByType = `-- name: RemoveQueueItemsByType :exec
DELETE FROM queue_items WHERE
queue_type = $1
`

func (q *Queries) RemoveQueueItemsByType(ctx context.Context, queueType string) error {
	_, err := q.db.ExecContext(ctx, removeQueueItemsByType, queueType)
	return err
}

const removeQueueItemsByTypeAndInstruction = `-- name: RemoveQueueItemsByTypeAndInstruction :exec
DELETE FROM queue_items WHERE
queue_type = $1 AND
instruction = $2
`

type RemoveQueueItemsByTypeAndInstructionParams struct {
	QueueType   string `json:"queue_type"`
	Instruction string `json:"instruction"`
}

func (q *Queries) RemoveQueueItemsByTypeAndInstruction(ctx context.Context, arg RemoveQueueItemsByTypeAndInstructionParams) error {
	_, err := q.db.ExecContext(ctx, removeQueueItemsByTypeAndInstruction, arg.QueueType, arg.Instruction)
	return err
}

const removeQueueItemsFilter = `-- name: RemoveQueueItemsFilter :exec
DELETE FROM queue_items WHERE
"status" = $1 AND
queue_type = $2 AND
instruction = $3
`

type RemoveQueueItemsFilterParams struct {
	Status      string `json:"status"`
	QueueType   string `json:"queue_type"`
	Instruction string `json:"instruction"`
}

func (q *Queries) RemoveQueueItemsFilter(ctx context.Context, arg RemoveQueueItemsFilterParams) error {
	_, err := q.db.ExecContext(ctx, removeQueueItemsFilter, arg.Status, arg.QueueType, arg.Instruction)
	return err
}

const updateQueueItem = `-- name: UpdateQueueItem :one
UPDATE queue_items
SET
    "status" = $1,
    updated_at = $2,
    "description" = $3
WHERE id = $4
RETURNING id, queue_type, instruction, status, object, created_at, updated_at, description
`

type UpdateQueueItemParams struct {
	Status      string    `json:"status"`
	UpdatedAt   time.Time `json:"updated_at"`
	Description string    `json:"description"`
	ID          uuid.UUID `json:"id"`
}

func (q *Queries) UpdateQueueItem(ctx context.Context, arg UpdateQueueItemParams) (QueueItem, error) {
	row := q.db.QueryRowContext(ctx, updateQueueItem,
		arg.Status,
		arg.UpdatedAt,
		arg.Description,
		arg.ID,
	)
	var i QueueItem
	err := row.Scan(
		&i.ID,
		&i.QueueType,
		&i.Instruction,
		&i.Status,
		&i.Object,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Description,
	)
	return i, err
}
