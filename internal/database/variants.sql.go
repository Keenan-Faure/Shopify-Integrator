// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: variants.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createVariant = `-- name: CreateVariant :one
INSERT INTO variants(
    id,
    product_id,
    sku,
    option1,
    option2,
    option3,
    barcode,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING id, product_id, sku, option1, option2, option3, barcode, created_at, updated_at
`

type CreateVariantParams struct {
	ID        uuid.UUID      `json:"id"`
	ProductID uuid.UUID      `json:"product_id"`
	Sku       string         `json:"sku"`
	Option1   sql.NullString `json:"option1"`
	Option2   sql.NullString `json:"option2"`
	Option3   sql.NullString `json:"option3"`
	Barcode   sql.NullString `json:"barcode"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (q *Queries) CreateVariant(ctx context.Context, arg CreateVariantParams) (Variant, error) {
	row := q.db.QueryRowContext(ctx, createVariant,
		arg.ID,
		arg.ProductID,
		arg.Sku,
		arg.Option1,
		arg.Option2,
		arg.Option3,
		arg.Barcode,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Variant
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.Sku,
		&i.Option1,
		&i.Option2,
		&i.Option3,
		&i.Barcode,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProductVariants = `-- name: GetProductVariants :many
SELECT
    id,
    sku,
    option1,
    option2,
    option3,
    barcode,
    updated_at
FROM variants
WHERE product_id = $1
`

type GetProductVariantsRow struct {
	ID        uuid.UUID      `json:"id"`
	Sku       string         `json:"sku"`
	Option1   sql.NullString `json:"option1"`
	Option2   sql.NullString `json:"option2"`
	Option3   sql.NullString `json:"option3"`
	Barcode   sql.NullString `json:"barcode"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (q *Queries) GetProductVariants(ctx context.Context, productID uuid.UUID) ([]GetProductVariantsRow, error) {
	rows, err := q.db.QueryContext(ctx, getProductVariants, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductVariantsRow
	for rows.Next() {
		var i GetProductVariantsRow
		if err := rows.Scan(
			&i.ID,
			&i.Sku,
			&i.Option1,
			&i.Option2,
			&i.Option3,
			&i.Barcode,
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

const getVariantBySKU = `-- name: GetVariantBySKU :one
SELECT
    id,
    sku,
    option1,
    option2,
    option3,
    barcode
FROM variants
WHERE sku = $1
`

type GetVariantBySKURow struct {
	ID      uuid.UUID      `json:"id"`
	Sku     string         `json:"sku"`
	Option1 sql.NullString `json:"option1"`
	Option2 sql.NullString `json:"option2"`
	Option3 sql.NullString `json:"option3"`
	Barcode sql.NullString `json:"barcode"`
}

func (q *Queries) GetVariantBySKU(ctx context.Context, sku string) (GetVariantBySKURow, error) {
	row := q.db.QueryRowContext(ctx, getVariantBySKU, sku)
	var i GetVariantBySKURow
	err := row.Scan(
		&i.ID,
		&i.Sku,
		&i.Option1,
		&i.Option2,
		&i.Option3,
		&i.Barcode,
	)
	return i, err
}

const getVariantIDByCode = `-- name: GetVariantIDByCode :one
SELECT
    id
FROM variants
WHERE sku = $1
`

func (q *Queries) GetVariantIDByCode(ctx context.Context, sku string) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getVariantIDByCode, sku)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const removeVariant = `-- name: RemoveVariant :exec
DELETE FROM variants
WHERE id = $1
`

func (q *Queries) RemoveVariant(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, removeVariant, id)
	return err
}

const updateVariant = `-- name: UpdateVariant :exec
UPDATE variants
SET
    option1 = $1,
    option2 = $2,
    option3 = $3,
    barcode = $4,
    updated_at = $5
WHERE sku = $6
`

type UpdateVariantParams struct {
	Option1   sql.NullString `json:"option1"`
	Option2   sql.NullString `json:"option2"`
	Option3   sql.NullString `json:"option3"`
	Barcode   sql.NullString `json:"barcode"`
	UpdatedAt time.Time      `json:"updated_at"`
	Sku       string         `json:"sku"`
}

func (q *Queries) UpdateVariant(ctx context.Context, arg UpdateVariantParams) error {
	_, err := q.db.ExecContext(ctx, updateVariant,
		arg.Option1,
		arg.Option2,
		arg.Option3,
		arg.Barcode,
		arg.UpdatedAt,
		arg.Sku,
	)
	return err
}
