// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: shopify_inventory.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createShopifyInventoryRecord = `-- name: CreateShopifyInventoryRecord :exec
INSERT INTO shopify_inventory(
    id,
    shopify_location_id,
    inventory_item_id,
    available,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)
`

type CreateShopifyInventoryRecordParams struct {
	ID                uuid.UUID `json:"id"`
	ShopifyLocationID string    `json:"shopify_location_id"`
	InventoryItemID   string    `json:"inventory_item_id"`
	Available         int32     `json:"available"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (q *Queries) CreateShopifyInventoryRecord(ctx context.Context, arg CreateShopifyInventoryRecordParams) error {
	_, err := q.db.ExecContext(ctx, createShopifyInventoryRecord,
		arg.ID,
		arg.ShopifyLocationID,
		arg.InventoryItemID,
		arg.Available,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const getShopifyInventory = `-- name: GetShopifyInventory :one
SELECT
    available,
    created_at
FROM shopify_inventory
WHERE inventory_item_id = $1 AND
shopify_location_id = $2
`

type GetShopifyInventoryParams struct {
	InventoryItemID   string `json:"inventory_item_id"`
	ShopifyLocationID string `json:"shopify_location_id"`
}

type GetShopifyInventoryRow struct {
	Available int32     `json:"available"`
	CreatedAt time.Time `json:"created_at"`
}

func (q *Queries) GetShopifyInventory(ctx context.Context, arg GetShopifyInventoryParams) (GetShopifyInventoryRow, error) {
	row := q.db.QueryRowContext(ctx, getShopifyInventory, arg.InventoryItemID, arg.ShopifyLocationID)
	var i GetShopifyInventoryRow
	err := row.Scan(&i.Available, &i.CreatedAt)
	return i, err
}

const removeShopifyInventoryRecord = `-- name: RemoveShopifyInventoryRecord :exec
DELETE FROM shopify_inventory
WHERE shopify_location_id = $1
AND inventory_item_id = $2
`

type RemoveShopifyInventoryRecordParams struct {
	ShopifyLocationID string `json:"shopify_location_id"`
	InventoryItemID   string `json:"inventory_item_id"`
}

func (q *Queries) RemoveShopifyInventoryRecord(ctx context.Context, arg RemoveShopifyInventoryRecordParams) error {
	_, err := q.db.ExecContext(ctx, removeShopifyInventoryRecord, arg.ShopifyLocationID, arg.InventoryItemID)
	return err
}

const updateShopifyInventoryRecord = `-- name: UpdateShopifyInventoryRecord :exec
UPDATE shopify_inventory
SET
    available = $1,
    updated_at = $2
WHERE shopify_location_id = $3
AND inventory_item_id = $4
`

type UpdateShopifyInventoryRecordParams struct {
	Available         int32     `json:"available"`
	UpdatedAt         time.Time `json:"updated_at"`
	ShopifyLocationID string    `json:"shopify_location_id"`
	InventoryItemID   string    `json:"inventory_item_id"`
}

func (q *Queries) UpdateShopifyInventoryRecord(ctx context.Context, arg UpdateShopifyInventoryRecordParams) error {
	_, err := q.db.ExecContext(ctx, updateShopifyInventoryRecord,
		arg.Available,
		arg.UpdatedAt,
		arg.ShopifyLocationID,
		arg.InventoryItemID,
	)
	return err
}