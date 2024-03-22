// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: shopify_location.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createShopifyLocation = `-- name: CreateShopifyLocation :one
INSERT INTO shopify_location(
    ID,
    shopify_warehouse_name,
    shopify_location_id,
    warehouse_name,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, shopify_warehouse_name, shopify_location_id, warehouse_name, created_at, updated_at
`

type CreateShopifyLocationParams struct {
	ID                   uuid.UUID `json:"id"`
	ShopifyWarehouseName string    `json:"shopify_warehouse_name"`
	ShopifyLocationID    string    `json:"shopify_location_id"`
	WarehouseName        string    `json:"warehouse_name"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func (q *Queries) CreateShopifyLocation(ctx context.Context, arg CreateShopifyLocationParams) (ShopifyLocation, error) {
	row := q.db.QueryRowContext(ctx, createShopifyLocation,
		arg.ID,
		arg.ShopifyWarehouseName,
		arg.ShopifyLocationID,
		arg.WarehouseName,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i ShopifyLocation
	err := row.Scan(
		&i.ID,
		&i.ShopifyWarehouseName,
		&i.ShopifyLocationID,
		&i.WarehouseName,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getShopifyLocationByLocationID = `-- name: GetShopifyLocationByLocationID :one
SELECT
    id,
    shopify_warehouse_name,
    shopify_location_id,
    warehouse_name,
    created_at
FROM shopify_location
WHERE shopify_location_id = $1
`

type GetShopifyLocationByLocationIDRow struct {
	ID                   uuid.UUID `json:"id"`
	ShopifyWarehouseName string    `json:"shopify_warehouse_name"`
	ShopifyLocationID    string    `json:"shopify_location_id"`
	WarehouseName        string    `json:"warehouse_name"`
	CreatedAt            time.Time `json:"created_at"`
}

func (q *Queries) GetShopifyLocationByLocationID(ctx context.Context, shopifyLocationID string) (GetShopifyLocationByLocationIDRow, error) {
	row := q.db.QueryRowContext(ctx, getShopifyLocationByLocationID, shopifyLocationID)
	var i GetShopifyLocationByLocationIDRow
	err := row.Scan(
		&i.ID,
		&i.ShopifyWarehouseName,
		&i.ShopifyLocationID,
		&i.WarehouseName,
		&i.CreatedAt,
	)
	return i, err
}

const getShopifyLocationByWarehouse = `-- name: GetShopifyLocationByWarehouse :one
SELECT
    id,
    shopify_warehouse_name,
    shopify_location_id,
    warehouse_name,
    created_at
FROM shopify_location
WHERE warehouse_name = $1
`

type GetShopifyLocationByWarehouseRow struct {
	ID                   uuid.UUID `json:"id"`
	ShopifyWarehouseName string    `json:"shopify_warehouse_name"`
	ShopifyLocationID    string    `json:"shopify_location_id"`
	WarehouseName        string    `json:"warehouse_name"`
	CreatedAt            time.Time `json:"created_at"`
}

func (q *Queries) GetShopifyLocationByWarehouse(ctx context.Context, warehouseName string) (GetShopifyLocationByWarehouseRow, error) {
	row := q.db.QueryRowContext(ctx, getShopifyLocationByWarehouse, warehouseName)
	var i GetShopifyLocationByWarehouseRow
	err := row.Scan(
		&i.ID,
		&i.ShopifyWarehouseName,
		&i.ShopifyLocationID,
		&i.WarehouseName,
		&i.CreatedAt,
	)
	return i, err
}

const getShopifyLocations = `-- name: GetShopifyLocations :many
SELECT id, shopify_warehouse_name, shopify_location_id, warehouse_name, created_at, updated_at FROM shopify_location
LIMIT $1 OFFSET $2
`

type GetShopifyLocationsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetShopifyLocations(ctx context.Context, arg GetShopifyLocationsParams) ([]ShopifyLocation, error) {
	rows, err := q.db.QueryContext(ctx, getShopifyLocations, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ShopifyLocation
	for rows.Next() {
		var i ShopifyLocation
		if err := rows.Scan(
			&i.ID,
			&i.ShopifyWarehouseName,
			&i.ShopifyLocationID,
			&i.WarehouseName,
			&i.CreatedAt,
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

const removeShopifyLocationMap = `-- name: RemoveShopifyLocationMap :exec
DELETE FROM shopify_location
WHERE id = $1
`

func (q *Queries) RemoveShopifyLocationMap(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, removeShopifyLocationMap, id)
	return err
}

const removeShopifyLocationMapByLocationID = `-- name: RemoveShopifyLocationMapByLocationID :exec
DELETE FROM shopify_location
WHERE shopify_location_id = $1
`

func (q *Queries) RemoveShopifyLocationMapByLocationID(ctx context.Context, shopifyLocationID string) error {
	_, err := q.db.ExecContext(ctx, removeShopifyLocationMapByLocationID, shopifyLocationID)
	return err
}

const removeShopifyLocationMapByWarehouse = `-- name: RemoveShopifyLocationMapByWarehouse :exec
DELETE FROM shopify_location
WHERE warehouse_name = $1
`

func (q *Queries) RemoveShopifyLocationMapByWarehouse(ctx context.Context, warehouseName string) error {
	_, err := q.db.ExecContext(ctx, removeShopifyLocationMapByWarehouse, warehouseName)
	return err
}

const updateShopifyLocation = `-- name: UpdateShopifyLocation :exec
UPDATE shopify_location
SET
    shopify_warehouse_name = $1,
    shopify_location_id = $2,
    updated_at = $3
WHERE warehouse_name = $4
`

type UpdateShopifyLocationParams struct {
	ShopifyWarehouseName string    `json:"shopify_warehouse_name"`
	ShopifyLocationID    string    `json:"shopify_location_id"`
	UpdatedAt            time.Time `json:"updated_at"`
	WarehouseName        string    `json:"warehouse_name"`
}

func (q *Queries) UpdateShopifyLocation(ctx context.Context, arg UpdateShopifyLocationParams) error {
	_, err := q.db.ExecContext(ctx, updateShopifyLocation,
		arg.ShopifyWarehouseName,
		arg.ShopifyLocationID,
		arg.UpdatedAt,
		arg.WarehouseName,
	)
	return err
}
