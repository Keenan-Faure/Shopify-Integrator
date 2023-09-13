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

const updateVariantQty = `-- name: UpdateVariantQty :one
UPDATE variant_qty
SET
    name = $1,
    value = $2
WHERE variant_id = $3
RETURNING id, variant_id, name, value, created_at, updated_at
`

type UpdateVariantQtyParams struct {
	Name      string        `json:"name"`
	Value     sql.NullInt32 `json:"value"`
	VariantID uuid.UUID     `json:"variant_id"`
}

func (q *Queries) UpdateVariantQty(ctx context.Context, arg UpdateVariantQtyParams) (VariantQty, error) {
	row := q.db.QueryRowContext(ctx, updateVariantQty, arg.Name, arg.Value, arg.VariantID)
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
