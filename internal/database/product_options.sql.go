// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: product_options.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createProductOption = `-- name: CreateProductOption :one
INSERT INTO product_options(
    product_id,
    name,
    value
) VALUES (
    $1, $2, $3
)
RETURNING id, product_id, name, value
`

type CreateProductOptionParams struct {
	ProductID uuid.UUID `json:"product_id"`
	Name      string    `json:"name"`
	Value     string    `json:"value"`
}

func (q *Queries) CreateProductOption(ctx context.Context, arg CreateProductOptionParams) (ProductOption, error) {
	row := q.db.QueryRowContext(ctx, createProductOption, arg.ProductID, arg.Name, arg.Value)
	var i ProductOption
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.Name,
		&i.Value,
	)
	return i, err
}

const getProductOptions = `-- name: GetProductOptions :many
SELECT
    name,
    value
FROM product_options
WHERE id = $1
`

type GetProductOptionsRow struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (q *Queries) GetProductOptions(ctx context.Context, id uuid.UUID) ([]GetProductOptionsRow, error) {
	rows, err := q.db.QueryContext(ctx, getProductOptions, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductOptionsRow
	for rows.Next() {
		var i GetProductOptionsRow
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

const updateProductOption = `-- name: UpdateProductOption :one
UPDATE product_options
SET
    name = $1,
    value = $2
WHERE product_id = $3
RETURNING id, product_id, name, value
`

type UpdateProductOptionParams struct {
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	ProductID uuid.UUID `json:"product_id"`
}

func (q *Queries) UpdateProductOption(ctx context.Context, arg UpdateProductOptionParams) (ProductOption, error) {
	row := q.db.QueryRowContext(ctx, updateProductOption, arg.Name, arg.Value, arg.ProductID)
	var i ProductOption
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.Name,
		&i.Value,
	)
	return i, err
}
