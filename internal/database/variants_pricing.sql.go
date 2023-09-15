// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: variants_pricing.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createVariantPricing = `-- name: CreateVariantPricing :one
INSERT INTO variant_pricing(
    id,
    variant_id,
    name,
    value,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id, variant_id, name, value, created_at, updated_at
`

type CreateVariantPricingParams struct {
	ID        uuid.UUID      `json:"id"`
	VariantID uuid.UUID      `json:"variant_id"`
	Name      string         `json:"name"`
	Value     sql.NullString `json:"value"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (q *Queries) CreateVariantPricing(ctx context.Context, arg CreateVariantPricingParams) (VariantPricing, error) {
	row := q.db.QueryRowContext(ctx, createVariantPricing,
		arg.ID,
		arg.VariantID,
		arg.Name,
		arg.Value,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i VariantPricing
	err := row.Scan(
		&i.ID,
		&i.VariantID,
		&i.Name,
		&i.Value,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getVariantPricing = `-- name: GetVariantPricing :many
SELECT 
    name,
    value
FROM variant_pricing
WHERE variant_id = $1
`

type GetVariantPricingRow struct {
	Name  string         `json:"name"`
	Value sql.NullString `json:"value"`
}

func (q *Queries) GetVariantPricing(ctx context.Context, variantID uuid.UUID) ([]GetVariantPricingRow, error) {
	rows, err := q.db.QueryContext(ctx, getVariantPricing, variantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetVariantPricingRow
	for rows.Next() {
		var i GetVariantPricingRow
		if err := rows.Scan(&i.Name, &i.Value); err != nil {
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

const removePricing = `-- name: RemovePricing :exec
DELETE FROM variant_pricing
WHERE id = $1
`

func (q *Queries) RemovePricing(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, removePricing, id)
	return err
}

const updateVariantPricing = `-- name: UpdateVariantPricing :one
UPDATE variant_pricing
SET
    name = $1,
    value = $2
WHERE variant_id = $3
RETURNING id, variant_id, name, value, created_at, updated_at
`

type UpdateVariantPricingParams struct {
	Name      string         `json:"name"`
	Value     sql.NullString `json:"value"`
	VariantID uuid.UUID      `json:"variant_id"`
}

func (q *Queries) UpdateVariantPricing(ctx context.Context, arg UpdateVariantPricingParams) (VariantPricing, error) {
	row := q.db.QueryRowContext(ctx, updateVariantPricing, arg.Name, arg.Value, arg.VariantID)
	var i VariantPricing
	err := row.Scan(
		&i.ID,
		&i.VariantID,
		&i.Name,
		&i.Value,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
