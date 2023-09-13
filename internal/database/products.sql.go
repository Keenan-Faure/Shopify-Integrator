// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: products.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createProduct = `-- name: CreateProduct :one
INSERT INTO products(
    id,
    active,
    title,
    body_html,
    category,
    vendor,
    product_type,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING id, active, title, body_html, category, vendor, product_type, created_at, updated_at
`

type CreateProductParams struct {
	ID          uuid.UUID      `json:"id"`
	Active      string         `json:"active"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, createProduct,
		arg.ID,
		arg.Active,
		arg.Title,
		arg.BodyHtml,
		arg.Category,
		arg.Vendor,
		arg.ProductType,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Active,
		&i.Title,
		&i.BodyHtml,
		&i.Category,
		&i.Vendor,
		&i.ProductType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProductByID = `-- name: GetProductByID :one
SELECT
    active,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE id = $1
`

type GetProductByIDRow struct {
	Active      string         `json:"active"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (q *Queries) GetProductByID(ctx context.Context, id uuid.UUID) (GetProductByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getProductByID, id)
	var i GetProductByIDRow
	err := row.Scan(
		&i.Active,
		&i.Title,
		&i.BodyHtml,
		&i.Category,
		&i.Vendor,
		&i.ProductType,
		&i.UpdatedAt,
	)
	return i, err
}

const getProducts = `-- name: GetProducts :many
SELECT
    id,
    active,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
LIMIT $1 OFFSET $2
`

type GetProductsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetProductsRow struct {
	ID          uuid.UUID      `json:"id"`
	Active      string         `json:"active"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (q *Queries) GetProducts(ctx context.Context, arg GetProductsParams) ([]GetProductsRow, error) {
	rows, err := q.db.QueryContext(ctx, getProducts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductsRow
	for rows.Next() {
		var i GetProductsRow
		if err := rows.Scan(
			&i.ID,
			&i.Active,
			&i.Title,
			&i.BodyHtml,
			&i.Category,
			&i.Vendor,
			&i.ProductType,
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

const getProductsByCategory = `-- name: GetProductsByCategory :many
SELECT
    id,
    title,
    body_html,
    category,
    vendor,
    product_type
FROM products
WHERE category LIKE $1
LIMIT $2 OFFSET $3
`

type GetProductsByCategoryParams struct {
	Category sql.NullString `json:"category"`
	Limit    int32          `json:"limit"`
	Offset   int32          `json:"offset"`
}

type GetProductsByCategoryRow struct {
	ID          uuid.UUID      `json:"id"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
}

func (q *Queries) GetProductsByCategory(ctx context.Context, arg GetProductsByCategoryParams) ([]GetProductsByCategoryRow, error) {
	rows, err := q.db.QueryContext(ctx, getProductsByCategory, arg.Category, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductsByCategoryRow
	for rows.Next() {
		var i GetProductsByCategoryRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.BodyHtml,
			&i.Category,
			&i.Vendor,
			&i.ProductType,
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

const getProductsByType = `-- name: GetProductsByType :many
SELECT
    id,
    title,
    body_html,
    category,
    vendor,
    product_type
FROM products
WHERE product_type LIKE $1
LIMIT $2 OFFSET $3
`

type GetProductsByTypeParams struct {
	ProductType sql.NullString `json:"product_type"`
	Limit       int32          `json:"limit"`
	Offset      int32          `json:"offset"`
}

type GetProductsByTypeRow struct {
	ID          uuid.UUID      `json:"id"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
}

func (q *Queries) GetProductsByType(ctx context.Context, arg GetProductsByTypeParams) ([]GetProductsByTypeRow, error) {
	rows, err := q.db.QueryContext(ctx, getProductsByType, arg.ProductType, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductsByTypeRow
	for rows.Next() {
		var i GetProductsByTypeRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.BodyHtml,
			&i.Category,
			&i.Vendor,
			&i.ProductType,
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

const getProductsByVendor = `-- name: GetProductsByVendor :many
SELECT
    id,
    title,
    body_html,
    category,
    vendor,
    product_type
FROM products
WHERE vendor LIKE $1
LIMIT $2 OFFSET $3
`

type GetProductsByVendorParams struct {
	Vendor sql.NullString `json:"vendor"`
	Limit  int32          `json:"limit"`
	Offset int32          `json:"offset"`
}

type GetProductsByVendorRow struct {
	ID          uuid.UUID      `json:"id"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
}

func (q *Queries) GetProductsByVendor(ctx context.Context, arg GetProductsByVendorParams) ([]GetProductsByVendorRow, error) {
	rows, err := q.db.QueryContext(ctx, getProductsByVendor, arg.Vendor, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductsByVendorRow
	for rows.Next() {
		var i GetProductsByVendorRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.BodyHtml,
			&i.Category,
			&i.Vendor,
			&i.ProductType,
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

const getProductsSearchSKU = `-- name: GetProductsSearchSKU :many
SELECT
    p.id,
    p.title,
    p.category,
    p.vendor,
    p.product_type
FROM products p
INNER JOIN variants v
ON p.id = variants.product_id
WHERE v.sku LIKE $1
LIMIT 5
`

type GetProductsSearchSKURow struct {
	ID          uuid.UUID      `json:"id"`
	Title       sql.NullString `json:"title"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
}

func (q *Queries) GetProductsSearchSKU(ctx context.Context, sku string) ([]GetProductsSearchSKURow, error) {
	rows, err := q.db.QueryContext(ctx, getProductsSearchSKU, sku)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductsSearchSKURow
	for rows.Next() {
		var i GetProductsSearchSKURow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Category,
			&i.Vendor,
			&i.ProductType,
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

const getProductsSearchTitle = `-- name: GetProductsSearchTitle :many
SELECT
    id,
    title,
    category,
    vendor,
    product_type
FROM products
WHERE title LIKE $1
LIMIT 5
`

type GetProductsSearchTitleRow struct {
	ID          uuid.UUID      `json:"id"`
	Title       sql.NullString `json:"title"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
}

func (q *Queries) GetProductsSearchTitle(ctx context.Context, title sql.NullString) ([]GetProductsSearchTitleRow, error) {
	rows, err := q.db.QueryContext(ctx, getProductsSearchTitle, title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductsSearchTitleRow
	for rows.Next() {
		var i GetProductsSearchTitleRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Category,
			&i.Vendor,
			&i.ProductType,
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

const updateProduct = `-- name: UpdateProduct :one
UPDATE products
SET
    active = $1,
    title = $2,
    body_html = $3,
    category = $4,
    vendor = $5,
    product_type = $6,
    updated_at = $7
WHERE id = $8
RETURNING id, active, title, body_html, category, vendor, product_type, created_at, updated_at
`

type UpdateProductParams struct {
	Active      string         `json:"active"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	UpdatedAt   time.Time      `json:"updated_at"`
	ID          uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, updateProduct,
		arg.Active,
		arg.Title,
		arg.BodyHtml,
		arg.Category,
		arg.Vendor,
		arg.ProductType,
		arg.UpdatedAt,
		arg.ID,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Active,
		&i.Title,
		&i.BodyHtml,
		&i.Category,
		&i.Vendor,
		&i.ProductType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
