// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: warehouses.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createWarehouse = `-- name: CreateWarehouse :exec
INSERT INTO warehouses(
    id,
    name,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4
)
`

type CreateWarehouseParams struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) CreateWarehouse(ctx context.Context, arg CreateWarehouseParams) error {
	_, err := q.db.ExecContext(ctx, createWarehouse,
		arg.ID,
		arg.Name,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const getWarehouseByID = `-- name: GetWarehouseByID :one
SELECT id, name, created_at, updated_at FROM warehouses
WHERE id = $1
`

func (q *Queries) GetWarehouseByID(ctx context.Context, id uuid.UUID) (Warehouse, error) {
	row := q.db.QueryRowContext(ctx, getWarehouseByID, id)
	var i Warehouse
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getWarehouseByName = `-- name: GetWarehouseByName :one
SELECT id, name, created_at, updated_at FROM warehouses
WHERE name = $1
`

func (q *Queries) GetWarehouseByName(ctx context.Context, name string) (Warehouse, error) {
	row := q.db.QueryRowContext(ctx, getWarehouseByName, name)
	var i Warehouse
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getWarehouses = `-- name: GetWarehouses :many
SELECT
    id,
    name,
    updated_at
FROM warehouses
`

type GetWarehousesRow struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) GetWarehouses(ctx context.Context) ([]GetWarehousesRow, error) {
	rows, err := q.db.QueryContext(ctx, getWarehouses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetWarehousesRow
	for rows.Next() {
		var i GetWarehousesRow
		if err := rows.Scan(&i.ID, &i.Name, &i.UpdatedAt); err != nil {
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

const removeWarehouse = `-- name: RemoveWarehouse :exec
DELETE FROM warehouses
WHERE id = $1
`

func (q *Queries) RemoveWarehouse(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, removeWarehouse, id)
	return err
}
