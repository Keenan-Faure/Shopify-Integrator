// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: variants.sql

package database

import (
	"context"
	"database/sql"
)

const createVariant = `-- name: CreateVariant :execresult
INSERT INTO variants(
    product_id,
    sku,
    option1,
    option2,
    option3,
    barcode
) VALUES (
    ?, ?, ?, ?, ?, ?
)
`

type CreateVariantParams struct {
	ProductID []byte         `json:"product_id"`
	Sku       string         `json:"sku"`
	Option1   sql.NullString `json:"option1"`
	Option2   sql.NullString `json:"option2"`
	Option3   sql.NullString `json:"option3"`
	Barcode   sql.NullString `json:"barcode"`
}

func (q *Queries) CreateVariant(ctx context.Context, arg CreateVariantParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createVariant,
		arg.ProductID,
		arg.Sku,
		arg.Option1,
		arg.Option2,
		arg.Option3,
		arg.Barcode,
	)
}

const getProductVariants = `-- name: GetProductVariants :many
SELECT
    sku,
    option1,
    option2,
    option3,
    barcode
FROM variants
WHERE product_id = ?
`

type GetProductVariantsRow struct {
	Sku     string         `json:"sku"`
	Option1 sql.NullString `json:"option1"`
	Option2 sql.NullString `json:"option2"`
	Option3 sql.NullString `json:"option3"`
	Barcode sql.NullString `json:"barcode"`
}

func (q *Queries) GetProductVariants(ctx context.Context, productID []byte) ([]GetProductVariantsRow, error) {
	rows, err := q.db.QueryContext(ctx, getProductVariants, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductVariantsRow
	for rows.Next() {
		var i GetProductVariantsRow
		if err := rows.Scan(
			&i.Sku,
			&i.Option1,
			&i.Option2,
			&i.Option3,
			&i.Barcode,
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
    sku,
    option1,
    option2,
    option3,
    barcode
FROM variants
WHERE sku = ?
`

type GetVariantBySKURow struct {
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
		&i.Sku,
		&i.Option1,
		&i.Option2,
		&i.Option3,
		&i.Barcode,
	)
	return i, err
}

const updateVariant = `-- name: UpdateVariant :execresult
UPDATE variants
SET
    option1 = ?,
    option2 = ?,
    option3 = ?,
    barcode = ?
WHERE sku = ?
`

type UpdateVariantParams struct {
	Option1 sql.NullString `json:"option1"`
	Option2 sql.NullString `json:"option2"`
	Option3 sql.NullString `json:"option3"`
	Barcode sql.NullString `json:"barcode"`
	Sku     string         `json:"sku"`
}

func (q *Queries) UpdateVariant(ctx context.Context, arg UpdateVariantParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateVariant,
		arg.Option1,
		arg.Option2,
		arg.Option3,
		arg.Barcode,
		arg.Sku,
	)
}
