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
    id,
    product_id,
    name,
    position
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, product_id, name, position
`

type CreateProductOptionParams struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Name      string    `json:"name"`
	Position  int32     `json:"position"`
}

func (q *Queries) CreateProductOption(ctx context.Context, arg CreateProductOptionParams) (ProductOption, error) {
	row := q.db.QueryRowContext(ctx, createProductOption,
		arg.ID,
		arg.ProductID,
		arg.Name,
		arg.Position,
	)
	var i ProductOption
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.Name,
		&i.Position,
	)
	return i, err
}

const getProductOptions = `-- name: GetProductOptions :many
SELECT
    name,
    position
FROM product_options
WHERE product_id = $1
ORDER BY position ASC
`

type GetProductOptionsRow struct {
	Name     string `json:"name"`
	Position int32  `json:"position"`
}

func (q *Queries) GetProductOptions(ctx context.Context, productID uuid.UUID) ([]GetProductOptionsRow, error) {
	rows, err := q.db.QueryContext(ctx, getProductOptions, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductOptionsRow
	for rows.Next() {
		var i GetProductOptionsRow
		if err := rows.Scan(&i.Name, &i.Position); err != nil {
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

const getProductOptionsByCode = `-- name: GetProductOptionsByCode :many
SELECT
    name,
    position
FROM product_options
WHERE product_id IN (
    SELECT id
    FROM products
    WHERE product_code = $1
)
ORDER BY position ASC
`

type GetProductOptionsByCodeRow struct {
	Name     string `json:"name"`
	Position int32  `json:"position"`
}

func (q *Queries) GetProductOptionsByCode(ctx context.Context, productCode string) ([]GetProductOptionsByCodeRow, error) {
	rows, err := q.db.QueryContext(ctx, getProductOptionsByCode, productCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductOptionsByCodeRow
	for rows.Next() {
		var i GetProductOptionsByCodeRow
		if err := rows.Scan(&i.Name, &i.Position); err != nil {
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
    position = $2
WHERE product_id = $3
RETURNING id, product_id, name, position
`

type UpdateProductOptionParams struct {
	Name      string    `json:"name"`
	Position  int32     `json:"position"`
	ProductID uuid.UUID `json:"product_id"`
}

func (q *Queries) UpdateProductOption(ctx context.Context, arg UpdateProductOptionParams) (ProductOption, error) {
	row := q.db.QueryRowContext(ctx, updateProductOption, arg.Name, arg.Position, arg.ProductID)
	var i ProductOption
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.Name,
		&i.Position,
	)
	return i, err
}
