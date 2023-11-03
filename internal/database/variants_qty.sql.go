// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
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
    isdefault,
    value,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, variant_id, name, value, created_at, updated_at, isdefault
`

type CreateVariantQtyParams struct {
	ID        uuid.UUID     `json:"id"`
	VariantID uuid.UUID     `json:"variant_id"`
	Name      string        `json:"name"`
	Isdefault bool          `json:"isdefault"`
	Value     sql.NullInt32 `json:"value"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (q *Queries) CreateVariantQty(ctx context.Context, arg CreateVariantQtyParams) (VariantQty, error) {
	row := q.db.QueryRowContext(ctx, createVariantQty,
		arg.ID,
		arg.VariantID,
		arg.Name,
		arg.Isdefault,
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
		&i.Isdefault,
	)
	return i, err
}

const getVariantQty = `-- name: GetVariantQty :many
SELECT 
    name,
    value,
    isdefault
FROM variant_qty
WHERE variant_id = $1
`

type GetVariantQtyRow struct {
	Name      string        `json:"name"`
	Value     sql.NullInt32 `json:"value"`
	Isdefault bool          `json:"isdefault"`
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

const getVariantQtyBySKU = `-- name: GetVariantQtyBySKU :many
SELECT
    name,
    value,
    isdefault
FROM variant_qty
WHERE variant_id IN (
    SELECT id FROM variants
    WHERE sku = $1
)
`

type GetVariantQtyBySKURow struct {
	Name      string        `json:"name"`
	Value     sql.NullInt32 `json:"value"`
	Isdefault bool          `json:"isdefault"`
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
    name = $1,
    value = $2,
    isdefault = $3
WHERE variant_id IN (
    SELECT id FROM variants
    WHERE sku = $4
)
`

type UpdateVariantQtyParams struct {
	Name      string        `json:"name"`
	Value     sql.NullInt32 `json:"value"`
	Isdefault bool          `json:"isdefault"`
	Sku       string        `json:"sku"`
}

func (q *Queries) UpdateVariantQty(ctx context.Context, arg UpdateVariantQtyParams) error {
	_, err := q.db.ExecContext(ctx, updateVariantQty,
		arg.Name,
		arg.Value,
		arg.Isdefault,
		arg.Sku,
	)
	return err
}
