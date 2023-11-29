// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: product_images.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createProductImage = `-- name: CreateProductImage :exec
INSERT INTO product_images(
    id,
    product_id,
    image_url,
    position,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)
`

type CreateProductImageParams struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	ImageUrl  string    `json:"image_url"`
	Position  int32     `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) CreateProductImage(ctx context.Context, arg CreateProductImageParams) error {
	_, err := q.db.ExecContext(ctx, createProductImage,
		arg.ID,
		arg.ProductID,
		arg.ImageUrl,
		arg.Position,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const getMaxImagePosition = `-- name: GetMaxImagePosition :one
SELECT 
    CAST(COALESCE(MAX("position"), 0) AS INTEGER)
FROM product_images
`

func (q *Queries) GetMaxImagePosition(ctx context.Context) (int32, error) {
	row := q.db.QueryRowContext(ctx, getMaxImagePosition)
	var column_1 int32
	err := row.Scan(&column_1)
	return column_1, err
}

const getProductImageByProductID = `-- name: GetProductImageByProductID :many
SELECT id, product_id, image_url, position, created_at, updated_at FROM product_images
WHERE product_id = $1
`

func (q *Queries) GetProductImageByProductID(ctx context.Context, productID uuid.UUID) ([]ProductImage, error) {
	rows, err := q.db.QueryContext(ctx, getProductImageByProductID, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ProductImage
	for rows.Next() {
		var i ProductImage
		if err := rows.Scan(
			&i.ID,
			&i.ProductID,
			&i.ImageUrl,
			&i.Position,
			&i.CreatedAt,
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

const updateProductImage = `-- name: UpdateProductImage :exec
UPDATE product_images
SET
    image_url = $1,
    updated_at = $2
WHERE product_id = $3
AND position = $4
`

type UpdateProductImageParams struct {
	ImageUrl  string    `json:"image_url"`
	UpdatedAt time.Time `json:"updated_at"`
	ProductID uuid.UUID `json:"product_id"`
	Position  int32     `json:"position"`
}

func (q *Queries) UpdateProductImage(ctx context.Context, arg UpdateProductImageParams) error {
	_, err := q.db.ExecContext(ctx, updateProductImage,
		arg.ImageUrl,
		arg.UpdatedAt,
		arg.ProductID,
		arg.Position,
	)
	return err
}
