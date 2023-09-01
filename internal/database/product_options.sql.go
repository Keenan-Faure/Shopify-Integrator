// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: product_options.sql

package database

import (
	"context"
	"database/sql"
)

const createProductOption = `-- name: CreateProductOption :execresult
INSERT INTO product_options(
    id,
    product_id,
    name,
    value
) VALUES (?, ?, ?, ?)
`

type CreateProductOptionParams struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Value     string `json:"value"`
}

func (q *Queries) CreateProductOption(ctx context.Context, arg CreateProductOptionParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createProductOption,
		arg.ID,
		arg.ProductID,
		arg.Name,
		arg.Value,
	)
}

const updateProductOption = `-- name: UpdateProductOption :execresult
UPDATE product_options
SET
    name = ?,
    value = ?
WHERE product_id = ?
`

type UpdateProductOptionParams struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	ProductID string `json:"product_id"`
}

func (q *Queries) UpdateProductOption(ctx context.Context, arg UpdateProductOptionParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateProductOption, arg.Name, arg.Value, arg.ProductID)
}
