// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: variants_qty.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createVariantQty = `-- name: CreateVariantQty :one
INSERT INTO variant_qty(
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

type CreateVariantQtyParams struct {
	ID        uuid.UUID     `json:"id"`
	VariantID uuid.UUID     `json:"variant_id"`
	Name      string        `json:"name"`
	Value     sql.NullInt32 `json:"value"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (q *Queries) CreateVariantQty(ctx context.Context, arg CreateVariantQtyParams) (VariantQty, error) {
	row := q.db.QueryRowContext(ctx, createVariantQty,
		arg.ID,
		arg.VariantID,
		arg.Name,
		arg.Value,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i VariantQty
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

const getVariantQty = `-- name: GetVariantQty :many
SELECT 
    name,
    value
FROM variant_qty
WHERE variant_id = $1
`

type GetVariantQtyRow struct {
	Name  string        `json:"name"`
	Value sql.NullInt32 `json:"value"`
}

func (q *Queries) GetVariantQty(ctx context.Context, variantID uuid.UUID) ([]GetVariantQtyRow, error) {
	rows, err := q.db.QueryContext(ctx, getVariantQty, variantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetVariantQtyRow
	for rows.Next() {
		var i GetVariantQtyRow
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

const getVariantQtyBySKU = `-- name: GetVariantQtyBySKU :many
SELECT
    name,
    value
FROM variant_qty
WHERE variant_id IN (
    SELECT id FROM variants
    WHERE sku = $1
)
`

type GetVariantQtyBySKURow struct {
	Name  string        `json:"name"`
	Value sql.NullInt32 `json:"value"`
}

func (q *Queries) GetVariantQtyBySKU(ctx context.Context, sku string) ([]GetVariantQtyBySKURow, error) {
	rows, err := q.db.QueryContext(ctx, getVariantQtyBySKU, sku)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetVariantQtyBySKURow
	for rows.Next() {
		var i GetVariantQtyBySKURow
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

const removeQty = `-- name: RemoveQty :exec
DELETE FROM variant_qty
WHERE id = $1
`

func (q *Queries) RemoveQty(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, removeQty, id)
	return err
}

const updateVariantQty = `-- name: UpdateVariantQty :exec
UPDATE variant_qty
SET
    value = $1
WHERE variant_id IN (
    SELECT id FROM variants
    WHERE sku = $2
) AND name = $3
`

type UpdateVariantQtyParams struct {
	Value sql.NullInt32 `json:"value"`
	Sku   string        `json:"sku"`
	Name  string        `json:"name"`
}

func (q *Queries) UpdateVariantQty(ctx context.Context, arg UpdateVariantQtyParams) error {
	_, err := q.db.ExecContext(ctx, updateVariantQty, arg.Value, arg.Sku, arg.Name)
	return err
}
