// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
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
    variant_id,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)
`

type CreateVIDParams struct {
	ID               uuid.UUID `json:"id"`
	Sku              string    `json:"sku"`
	ShopifyVariantID string    `json:"shopify_variant_id"`
	VariantID        uuid.UUID `json:"variant_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (q *Queries) CreateVID(ctx context.Context, arg CreateVIDParams) error {
	_, err := q.db.ExecContext(ctx, createVID,
		arg.ID,
		arg.Sku,
		arg.ShopifyVariantID,
		arg.VariantID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
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

const updateVID = `-- name: UpdateVID :exec
UPDATE shopify_vid
SET
    shopify_variant_id = $1,
    updated_at = $2
WHERE sku = $3
`

type UpdateVIDParams struct {
	ShopifyVariantID string    `json:"shopify_variant_id"`
	UpdatedAt        time.Time `json:"updated_at"`
	Sku              string    `json:"sku"`
}

func (q *Queries) UpdateVID(ctx context.Context, arg UpdateVIDParams) error {
	_, err := q.db.ExecContext(ctx, updateVID, arg.ShopifyVariantID, arg.UpdatedAt, arg.Sku)
	return err
}
