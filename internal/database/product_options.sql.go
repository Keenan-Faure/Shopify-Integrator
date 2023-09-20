// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
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
    name
) VALUES (
    $1, $2, $3
)
RETURNING id, product_id, name
`

type CreateProductOptionParams struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Name      string    `json:"name"`
}

func (q *Queries) CreateProductOption(ctx context.Context, arg CreateProductOptionParams) (ProductOption, error) {
	row := q.db.QueryRowContext(ctx, createProductOption, arg.ID, arg.ProductID, arg.Name)
	var i ProductOption
	err := row.Scan(&i.ID, &i.ProductID, &i.Name)
	return i, err
}

const getProductOptions = `-- name: GetProductOptions :many
SELECT
    name
FROM product_options
WHERE product_id = $1
`

func (q *Queries) GetProductOptions(ctx context.Context, productID uuid.UUID) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getProductOptions, productID)
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

const getProductOptionsByCode = `-- name: GetProductOptionsByCode :many
SELECT
    name
FROM product_options
WHERE product_id IN (
    SELECT id
    FROM products
    WHERE product_code = $1
)
`

func (q *Queries) GetProductOptionsByCode(ctx context.Context, productCode string) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getProductOptionsByCode, productCode)
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

const updateProductOption = `-- name: UpdateProductOption :one
UPDATE product_options
SET
    name = $1
WHERE product_id = $2
RETURNING id, product_id, name
`

type UpdateProductOptionParams struct {
	Name      string    `json:"name"`
	ProductID uuid.UUID `json:"product_id"`
}

func (q *Queries) UpdateProductOption(ctx context.Context, arg UpdateProductOptionParams) (ProductOption, error) {
	row := q.db.QueryRowContext(ctx, updateProductOption, arg.Name, arg.ProductID)
	var i ProductOption
	err := row.Scan(&i.ID, &i.ProductID, &i.Name)
	return i, err
}
