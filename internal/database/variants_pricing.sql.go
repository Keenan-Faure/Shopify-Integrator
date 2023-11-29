// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
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
    isdefault,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, variant_id, name, value, created_at, updated_at, isdefault
`

type CreateVariantPricingParams struct {
	ID        uuid.UUID      `json:"id"`
	VariantID uuid.UUID      `json:"variant_id"`
	Name      string         `json:"name"`
	Value     sql.NullString `json:"value"`
	Isdefault bool           `json:"isdefault"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (q *Queries) CreateVariantPricing(ctx context.Context, arg CreateVariantPricingParams) (VariantPricing, error) {
	row := q.db.QueryRowContext(ctx, createVariantPricing,
		arg.ID,
		arg.VariantID,
		arg.Name,
		arg.Value,
		arg.Isdefault,
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
		&i.Isdefault,
	)
	return i, err
}

const getCountOfUniquePrices = `-- name: GetCountOfUniquePrices :one
SELECT CAST(COALESCE(COUNT(DISTINCT "name"),0) AS INTEGER) FROM variant_pricing
`

func (q *Queries) GetCountOfUniquePrices(ctx context.Context) (int32, error) {
	row := q.db.QueryRowContext(ctx, getCountOfUniquePrices)
	var column_1 int32
	err := row.Scan(&column_1)
	return column_1, err
}

const getPriceTierBySKU = `-- name: GetPriceTierBySKU :one
SELECT 
    name,
    value,
    isdefault
FROM variant_pricing
WHERE variant_id IN (
    SELECT id FROM variants
    WHERE sku = $1
) AND name = $2
`

type GetPriceTierBySKUParams struct {
	Sku  string `json:"sku"`
	Name string `json:"name"`
}

type GetPriceTierBySKURow struct {
	Name      string         `json:"name"`
	Value     sql.NullString `json:"value"`
	Isdefault bool           `json:"isdefault"`
}

func (q *Queries) GetPriceTierBySKU(ctx context.Context, arg GetPriceTierBySKUParams) (GetPriceTierBySKURow, error) {
	row := q.db.QueryRowContext(ctx, getPriceTierBySKU, arg.Sku, arg.Name)
	var i GetPriceTierBySKURow
	err := row.Scan(&i.Name, &i.Value, &i.Isdefault)
	return i, err
}

const getUniquePriceTiers = `-- name: GetUniquePriceTiers :many
SELECT DISTINCT "name" FROM variant_pricing
`

func (q *Queries) GetUniquePriceTiers(ctx context.Context) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getUniquePriceTiers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getVariantPricing = `-- name: GetVariantPricing :many
SELECT 
    name,
    value,
    isdefault
FROM variant_pricing
WHERE variant_id = $1
`

type GetVariantPricingRow struct {
	Name      string         `json:"name"`
	Value     sql.NullString `json:"value"`
	Isdefault bool           `json:"isdefault"`
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
		if err := rows.Scan(&i.Name, &i.Value, &i.Isdefault); err != nil {
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

const getVariantPricingBySKU = `-- name: GetVariantPricingBySKU :many
SELECT
    name,
    value,
    isdefault
FROM variant_pricing
WHERE variant_id IN (
    SELECT id FROM variants
    WHERE sku = $1
)
`

type GetVariantPricingBySKURow struct {
	Name      string         `json:"name"`
	Value     sql.NullString `json:"value"`
	Isdefault bool           `json:"isdefault"`
}

func (q *Queries) GetVariantPricingBySKU(ctx context.Context, sku string) ([]GetVariantPricingBySKURow, error) {
	rows, err := q.db.QueryContext(ctx, getVariantPricingBySKU, sku)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetVariantPricingBySKURow
	for rows.Next() {
		var i GetVariantPricingBySKURow
		if err := rows.Scan(&i.Name, &i.Value, &i.Isdefault); err != nil {
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

const updateVariantPricing = `-- name: UpdateVariantPricing :exec
UPDATE variant_pricing
SET
    name = $1,
    value = $2,
    isdefault = $3
WHERE variant_id IN (
    SELECT id FROM variants
    WHERE sku = $4
) AND name = $5
`

type UpdateVariantPricingParams struct {
	Name      string         `json:"name"`
	Value     sql.NullString `json:"value"`
	Isdefault bool           `json:"isdefault"`
	Sku       string         `json:"sku"`
	Name_2    string         `json:"name_2"`
}

func (q *Queries) UpdateVariantPricing(ctx context.Context, arg UpdateVariantPricingParams) error {
	_, err := q.db.ExecContext(ctx, updateVariantPricing,
		arg.Name,
		arg.Value,
		arg.Isdefault,
		arg.Sku,
		arg.Name_2,
	)
	return err
}
