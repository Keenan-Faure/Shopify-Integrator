// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: shopify_vid.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createVID = `-- name: CreateVID :exec
INSERT INTO shopify_vid(
    id,
    sku,
    shopify_variant_id,
    shopify_inventory_id,
    variant_id,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
`

type CreateVIDParams struct {
	ID                 uuid.UUID `json:"id"`
	Sku                string    `json:"sku"`
	ShopifyVariantID   string    `json:"shopify_variant_id"`
	ShopifyInventoryID string    `json:"shopify_inventory_id"`
	VariantID          uuid.UUID `json:"variant_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (q *Queries) CreateVID(ctx context.Context, arg CreateVIDParams) error {
	_, err := q.db.ExecContext(ctx, createVID,
		arg.ID,
		arg.Sku,
		arg.ShopifyVariantID,
		arg.ShopifyInventoryID,
		arg.VariantID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const getInventoryIDBySKU = `-- name: GetInventoryIDBySKU :one
select
    sku,
    shopify_inventory_id,
    updated_at
from shopify_vid
where sku = $1
LIMIT 1
`

type GetInventoryIDBySKURow struct {
	Sku                string    `json:"sku"`
	ShopifyInventoryID string    `json:"shopify_inventory_id"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (q *Queries) GetInventoryIDBySKU(ctx context.Context, sku string) (GetInventoryIDBySKURow, error) {
	row := q.db.QueryRowContext(ctx, getInventoryIDBySKU, sku)
	var i GetInventoryIDBySKURow
	err := row.Scan(&i.Sku, &i.ShopifyInventoryID, &i.UpdatedAt)
	return i, err
}

const getVIDBySKU = `-- name: GetVIDBySKU :one
SELECT
    sku,
    shopify_variant_id,
    updated_at
FROM shopify_vid
WHERE sku = $1
LIMIT 1
`

type GetVIDBySKURow struct {
	Sku              string    `json:"sku"`
	ShopifyVariantID string    `json:"shopify_variant_id"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (q *Queries) GetVIDBySKU(ctx context.Context, sku string) (GetVIDBySKURow, error) {
	row := q.db.QueryRowContext(ctx, getVIDBySKU, sku)
	var i GetVIDBySKURow
	err := row.Scan(&i.Sku, &i.ShopifyVariantID, &i.UpdatedAt)
	return i, err
}

const removeShopifyVIDBySKU = `-- name: RemoveShopifyVIDBySKU :exec
DELETE FROM shopify_vid
WHERE sku = $1
`

func (q *Queries) RemoveShopifyVIDBySKU(ctx context.Context, sku string) error {
	_, err := q.db.ExecContext(ctx, removeShopifyVIDBySKU, sku)
	return err
}

const updateVID = `-- name: UpdateVID :exec
UPDATE shopify_vid
SET
    shopify_variant_id = $1,
    shopify_inventory_id = $2,
    updated_at = $3
WHERE sku = $4
`

type UpdateVIDParams struct {
	ShopifyVariantID   string    `json:"shopify_variant_id"`
	ShopifyInventoryID string    `json:"shopify_inventory_id"`
	UpdatedAt          time.Time `json:"updated_at"`
	Sku                string    `json:"sku"`
}

func (q *Queries) UpdateVID(ctx context.Context, arg UpdateVIDParams) error {
	_, err := q.db.ExecContext(ctx, updateVID,
		arg.ShopifyVariantID,
		arg.ShopifyInventoryID,
		arg.UpdatedAt,
		arg.Sku,
	)
	return err
}

const upsertVID = `-- name: UpsertVID :exec
INSERT INTO shopify_vid(
    id,
    sku,
    shopify_variant_id,
    shopify_inventory_id,
    variant_id,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT(sku)
DO UPDATE
SET
    shopify_variant_id = COALESCE($4, shopify_vid.shopify_variant_id),
    updated_at = $7
`

type UpsertVIDParams struct {
	ID                 uuid.UUID `json:"id"`
	Sku                string    `json:"sku"`
	ShopifyVariantID   string    `json:"shopify_variant_id"`
	ShopifyInventoryID string    `json:"shopify_inventory_id"`
	VariantID          uuid.UUID `json:"variant_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (q *Queries) UpsertVID(ctx context.Context, arg UpsertVIDParams) error {
	_, err := q.db.ExecContext(ctx, upsertVID,
		arg.ID,
		arg.Sku,
		arg.ShopifyVariantID,
		arg.ShopifyInventoryID,
		arg.VariantID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}
