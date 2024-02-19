// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
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
    product_code,
    active,
    title,
    body_html,
    category,
    vendor,
    product_type,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING id, active, product_code, title, body_html, category, vendor, product_type, created_at, updated_at
`

type CreateProductParams struct {
	ID          uuid.UUID      `json:"id"`
	ProductCode string         `json:"product_code"`
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
		arg.ProductCode,
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
		&i.ProductCode,
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

const getActiveProducts = `-- name: GetActiveProducts :many
SELECT
    id,
    active,
    product_code,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE active = '1'
LIMIT $1 OFFSET $2
`

type GetActiveProductsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetActiveProductsRow struct {
	ID          uuid.UUID      `json:"id"`
	Active      string         `json:"active"`
	ProductCode string         `json:"product_code"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (q *Queries) GetActiveProducts(ctx context.Context, arg GetActiveProductsParams) ([]GetActiveProductsRow, error) {
	rows, err := q.db.QueryContext(ctx, getActiveProducts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetActiveProductsRow
	for rows.Next() {
		var i GetActiveProductsRow
		if err := rows.Scan(
			&i.ID,
			&i.Active,
			&i.ProductCode,
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

const getProductByCategoryAndType = `-- name: GetProductByCategoryAndType :many
SELECT DISTINCT
    id,
    active,
    product_code,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE category ILIKE $1
AND product_type ILIKE $2
LIMIT $3 OFFSET $4
`

type GetProductByCategoryAndTypeParams struct {
	Category    sql.NullString `json:"category"`
	ProductType sql.NullString `json:"product_type"`
	Limit       int32          `json:"limit"`
	Offset      int32          `json:"offset"`
}

type GetProductByCategoryAndTypeRow struct {
	ID          uuid.UUID      `json:"id"`
	Active      string         `json:"active"`
	ProductCode string         `json:"product_code"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (q *Queries) GetProductByCategoryAndType(ctx context.Context, arg GetProductByCategoryAndTypeParams) ([]GetProductByCategoryAndTypeRow, error) {
	rows, err := q.db.QueryContext(ctx, getProductByCategoryAndType,
		arg.Category,
		arg.ProductType,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductByCategoryAndTypeRow
	for rows.Next() {
		var i GetProductByCategoryAndTypeRow
		if err := rows.Scan(
			&i.ID,
			&i.Active,
			&i.ProductCode,
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

const getProductByID = `-- name: GetProductByID :one
SELECT
    active,
    product_code,
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
	ProductCode string         `json:"product_code"`
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
		&i.ProductCode,
		&i.Title,
		&i.BodyHtml,
		&i.Category,
		&i.Vendor,
		&i.ProductType,
		&i.UpdatedAt,
	)
	return i, err
}

const getProductByProductCode = `-- name: GetProductByProductCode :one
SELECT DISTINCT
    active,
    product_code,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE product_code = $1
`

type GetProductByProductCodeRow struct {
	Active      string         `json:"active"`
	ProductCode string         `json:"product_code"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (q *Queries) GetProductByProductCode(ctx context.Context, productCode string) (GetProductByProductCodeRow, error) {
	row := q.db.QueryRowContext(ctx, getProductByProductCode, productCode)
	var i GetProductByProductCodeRow
	err := row.Scan(
		&i.Active,
		&i.ProductCode,
		&i.Title,
		&i.BodyHtml,
		&i.Category,
		&i.Vendor,
		&i.ProductType,
		&i.UpdatedAt,
	)
	return i, err
}

const getProductIDByCode = `-- name: GetProductIDByCode :one
SELECT
    id
FROM products
WHERE product_code = $1
`

func (q *Queries) GetProductIDByCode(ctx context.Context, productCode string) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getProductIDByCode, productCode)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getProductIDs = `-- name: GetProductIDs :many
SELECT id FROM products
`

func (q *Queries) GetProductIDs(ctx context.Context) ([]uuid.UUID, error) {
	rows, err := q.db.QueryContext(ctx, getProductIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProducts = `-- name: GetProducts :many
SELECT
    id,
    active,
    product_code,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
ORDER BY updated_at DESC
LIMIT $1 OFFSET $2
`

type GetProductsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetProductsRow struct {
	ID          uuid.UUID      `json:"id"`
	Active      string         `json:"active"`
	ProductCode string         `json:"product_code"`
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
			&i.ProductCode,
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
SELECT DISTINCT
    id,
    active,
    product_code,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE category ILIKE $1
LIMIT $2 OFFSET $3
`

type GetProductsByCategoryParams struct {
	Category sql.NullString `json:"category"`
	Limit    int32          `json:"limit"`
	Offset   int32          `json:"offset"`
}

type GetProductsByCategoryRow struct {
	ID          uuid.UUID      `json:"id"`
	Active      string         `json:"active"`
	ProductCode string         `json:"product_code"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	UpdatedAt   time.Time      `json:"updated_at"`
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
			&i.Active,
			&i.ProductCode,
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

const getProductsByType = `-- name: GetProductsByType :many
SELECT DISTINCT
    id,
    active,
    product_code,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE product_type ILIKE $1
LIMIT $2 OFFSET $3
`

type GetProductsByTypeParams struct {
	ProductType sql.NullString `json:"product_type"`
	Limit       int32          `json:"limit"`
	Offset      int32          `json:"offset"`
}

type GetProductsByTypeRow struct {
	ID          uuid.UUID      `json:"id"`
	Active      string         `json:"active"`
	ProductCode string         `json:"product_code"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	UpdatedAt   time.Time      `json:"updated_at"`
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
			&i.Active,
			&i.ProductCode,
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

const getProductsByTypeAndCategory = `-- name: GetProductsByTypeAndCategory :many
SELECT DISTINCT
    id,
    active,
    product_code,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE product_type ILIKE $1
AND category ILIKE $2
LIMIT $3 OFFSET $4
`

type GetProductsByTypeAndCategoryParams struct {
	ProductType sql.NullString `json:"product_type"`
	Category    sql.NullString `json:"category"`
	Limit       int32          `json:"limit"`
	Offset      int32          `json:"offset"`
}

type GetProductsByTypeAndCategoryRow struct {
	ID          uuid.UUID      `json:"id"`
	Active      string         `json:"active"`
	ProductCode string         `json:"product_code"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (q *Queries) GetProductsByTypeAndCategory(ctx context.Context, arg GetProductsByTypeAndCategoryParams) ([]GetProductsByTypeAndCategoryRow, error) {
	rows, err := q.db.QueryContext(ctx, getProductsByTypeAndCategory,
		arg.ProductType,
		arg.Category,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductsByTypeAndCategoryRow
	for rows.Next() {
		var i GetProductsByTypeAndCategoryRow
		if err := rows.Scan(
			&i.ID,
			&i.Active,
			&i.ProductCode,
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

const getProductsByTypeAndVendor = `-- name: GetProductsByTypeAndVendor :many
SELECT DISTINCT
    id,
    active,
    product_code,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE product_type ILIKE $1
AND vendor ILIKE $2
LIMIT $3 OFFSET $4
`

type GetProductsByTypeAndVendorParams struct {
	ProductType sql.NullString `json:"product_type"`
	Vendor      sql.NullString `json:"vendor"`
	Limit       int32          `json:"limit"`
	Offset      int32          `json:"offset"`
}

type GetProductsByTypeAndVendorRow struct {
	ID          uuid.UUID      `json:"id"`
	Active      string         `json:"active"`
	ProductCode string         `json:"product_code"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (q *Queries) GetProductsByTypeAndVendor(ctx context.Context, arg GetProductsByTypeAndVendorParams) ([]GetProductsByTypeAndVendorRow, error) {
	rows, err := q.db.QueryContext(ctx, getProductsByTypeAndVendor,
		arg.ProductType,
		arg.Vendor,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductsByTypeAndVendorRow
	for rows.Next() {
		var i GetProductsByTypeAndVendorRow
		if err := rows.Scan(
			&i.ID,
			&i.Active,
			&i.ProductCode,
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

const getProductsByVendor = `-- name: GetProductsByVendor :many
SELECT DISTINCT
    id,
    active,
    product_code,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE vendor ILIKE $1
LIMIT $2 OFFSET $3
`

type GetProductsByVendorParams struct {
	Vendor sql.NullString `json:"vendor"`
	Limit  int32          `json:"limit"`
	Offset int32          `json:"offset"`
}

type GetProductsByVendorRow struct {
	ID          uuid.UUID      `json:"id"`
	Active      string         `json:"active"`
	ProductCode string         `json:"product_code"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	UpdatedAt   time.Time      `json:"updated_at"`
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
			&i.Active,
			&i.ProductCode,
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

const getProductsByVendorAndCategory = `-- name: GetProductsByVendorAndCategory :many
SELECT DISTINCT
    id,
    active,
    product_code,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE vendor ILIKE $1
AND category ILIKE $2
LIMIT $3 OFFSET $4
`

type GetProductsByVendorAndCategoryParams struct {
	Vendor   sql.NullString `json:"vendor"`
	Category sql.NullString `json:"category"`
	Limit    int32          `json:"limit"`
	Offset   int32          `json:"offset"`
}

type GetProductsByVendorAndCategoryRow struct {
	ID          uuid.UUID      `json:"id"`
	Active      string         `json:"active"`
	ProductCode string         `json:"product_code"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (q *Queries) GetProductsByVendorAndCategory(ctx context.Context, arg GetProductsByVendorAndCategoryParams) ([]GetProductsByVendorAndCategoryRow, error) {
	rows, err := q.db.QueryContext(ctx, getProductsByVendorAndCategory,
		arg.Vendor,
		arg.Category,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductsByVendorAndCategoryRow
	for rows.Next() {
		var i GetProductsByVendorAndCategoryRow
		if err := rows.Scan(
			&i.ID,
			&i.Active,
			&i.ProductCode,
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

const getProductsFilter = `-- name: GetProductsFilter :many
SELECT DISTINCT
    id,
    active,
    product_code,
    title,
    body_html,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE category ILIKE $1
AND product_type ILIKE $2
AND vendor ILIKE $3
LIMIT $4 OFFSET $5
`

type GetProductsFilterParams struct {
	Category    sql.NullString `json:"category"`
	ProductType sql.NullString `json:"product_type"`
	Vendor      sql.NullString `json:"vendor"`
	Limit       int32          `json:"limit"`
	Offset      int32          `json:"offset"`
}

type GetProductsFilterRow struct {
	ID          uuid.UUID      `json:"id"`
	Active      string         `json:"active"`
	ProductCode string         `json:"product_code"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (q *Queries) GetProductsFilter(ctx context.Context, arg GetProductsFilterParams) ([]GetProductsFilterRow, error) {
	rows, err := q.db.QueryContext(ctx, getProductsFilter,
		arg.Category,
		arg.ProductType,
		arg.Vendor,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductsFilterRow
	for rows.Next() {
		var i GetProductsFilterRow
		if err := rows.Scan(
			&i.ID,
			&i.Active,
			&i.ProductCode,
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

const getProductsSearch = `-- name: GetProductsSearch :many
SELECT
    p.id,
    p.active,
    p.product_code,
    p.title,
    p.category,
    p.vendor,
    p.product_type,
    p.updated_at
FROM products p
INNER JOIN variants v
    ON p.id = v.product_id
WHERE v.sku ILIKE $1
UNION
SELECT
    id,
    active,
    product_code,
    title,
    category,
    vendor,
    product_type,
    updated_at
FROM products
WHERE title ILIKE $1
`

type GetProductsSearchRow struct {
	ID          uuid.UUID      `json:"id"`
	Active      string         `json:"active"`
	ProductCode string         `json:"product_code"`
	Title       sql.NullString `json:"title"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (q *Queries) GetProductsSearch(ctx context.Context, sku string) ([]GetProductsSearchRow, error) {
	rows, err := q.db.QueryContext(ctx, getProductsSearch, sku)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductsSearchRow
	for rows.Next() {
		var i GetProductsSearchRow
		if err := rows.Scan(
			&i.ID,
			&i.Active,
			&i.ProductCode,
			&i.Title,
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

const getVariantOptionsByProductCode = `-- name: GetVariantOptionsByProductCode :many
SELECT
    v.sku,
    v.option1,
    v.option2,
    v.option3
FROM variants v
WHERE v.product_id IN (
    SELECT product_id
    FROM products
    WHERE product_code = $1
)
`

type GetVariantOptionsByProductCodeRow struct {
	Sku     string         `json:"sku"`
	Option1 sql.NullString `json:"option1"`
	Option2 sql.NullString `json:"option2"`
	Option3 sql.NullString `json:"option3"`
}

func (q *Queries) GetVariantOptionsByProductCode(ctx context.Context, productCode string) ([]GetVariantOptionsByProductCodeRow, error) {
	rows, err := q.db.QueryContext(ctx, getVariantOptionsByProductCode, productCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetVariantOptionsByProductCodeRow
	for rows.Next() {
		var i GetVariantOptionsByProductCodeRow
		if err := rows.Scan(
			&i.Sku,
			&i.Option1,
			&i.Option2,
			&i.Option3,
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

const removeProduct = `-- name: RemoveProduct :exec
DELETE FROM products
WHERE id = $1
`

func (q *Queries) RemoveProduct(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, removeProduct, id)
	return err
}

const removeProductByCode = `-- name: RemoveProductByCode :exec
DELETE FROM products
WHERE product_code = $1
`

func (q *Queries) RemoveProductByCode(ctx context.Context, productCode string) error {
	_, err := q.db.ExecContext(ctx, removeProductByCode, productCode)
	return err
}

const updateProduct = `-- name: UpdateProduct :exec
UPDATE products
SET
    active = COALESCE($1, active),
    title = COALESCE($2, title),
    body_html = COALESCE($3, body_html),
    category = COALESCE($4, category),
    vendor = COALESCE($5, vendor),
    product_type = COALESCE($6, product_type),
    updated_at = $7
WHERE product_code = $8
`

type UpdateProductParams struct {
	Active      string         `json:"active"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	UpdatedAt   time.Time      `json:"updated_at"`
	ProductCode string         `json:"product_code"`
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) error {
	_, err := q.db.ExecContext(ctx, updateProduct,
		arg.Active,
		arg.Title,
		arg.BodyHtml,
		arg.Category,
		arg.Vendor,
		arg.ProductType,
		arg.UpdatedAt,
		arg.ProductCode,
	)
	return err
}

const updateProductByID = `-- name: UpdateProductByID :exec
UPDATE products
SET
    active = COALESCE($1, active),
    title = COALESCE($2, title),
    body_html = COALESCE($3, body_html),
    category = COALESCE($4, category),
    vendor = COALESCE($5, vendor),
    product_type = COALESCE($6, product_type),
    updated_at = $7
WHERE id = $8
`

type UpdateProductByIDParams struct {
	Active      string         `json:"active"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	UpdatedAt   time.Time      `json:"updated_at"`
	ID          uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateProductByID(ctx context.Context, arg UpdateProductByIDParams) error {
	_, err := q.db.ExecContext(ctx, updateProductByID,
		arg.Active,
		arg.Title,
		arg.BodyHtml,
		arg.Category,
		arg.Vendor,
		arg.ProductType,
		arg.UpdatedAt,
		arg.ID,
	)
	return err
}

const updateProductBySKU = `-- name: UpdateProductBySKU :exec
UPDATE products
SET
    title = COALESCE($1, title),
    body_html = COALESCE($2, body_html),
    category = COALESCE($3, category),
    vendor = COALESCE($4, vendor),
    product_type = COALESCE($5, product_type),
    updated_at = $6
WHERE id = (
    SELECT
        product_id
    FROM variants
    WHERE sku = $7
)
`

type UpdateProductBySKUParams struct {
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Sku         string         `json:"sku"`
}

func (q *Queries) UpdateProductBySKU(ctx context.Context, arg UpdateProductBySKUParams) error {
	_, err := q.db.ExecContext(ctx, updateProductBySKU,
		arg.Title,
		arg.BodyHtml,
		arg.Category,
		arg.Vendor,
		arg.ProductType,
		arg.UpdatedAt,
		arg.Sku,
	)
	return err
}

const upsertProduct = `-- name: UpsertProduct :one
INSERT INTO products(
    id,
    product_code,
    active,
    title,
    body_html,
    category,
    vendor,
    product_type,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT(product_code)
DO UPDATE 
SET
    active = COALESCE($3, products.active),
    title = COALESCE($4, products.title),
    body_html = COALESCE($5, products.body_html),
    category = COALESCE($6, products.category),
    vendor = COALESCE($7, products.vendor),
    product_type = COALESCE($8, products.product_type),
    updated_at = $9
RETURNING id, active, product_code, title, body_html, category, vendor, product_type, created_at, updated_at, (xmax = 0) AS inserted
`

type UpsertProductParams struct {
	ID          uuid.UUID      `json:"id"`
	ProductCode string         `json:"product_code"`
	Active      string         `json:"active"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type UpsertProductRow struct {
	ID          uuid.UUID      `json:"id"`
	Active      string         `json:"active"`
	ProductCode string         `json:"product_code"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Inserted    bool           `json:"inserted"`
}

func (q *Queries) UpsertProduct(ctx context.Context, arg UpsertProductParams) (UpsertProductRow, error) {
	row := q.db.QueryRowContext(ctx, upsertProduct,
		arg.ID,
		arg.ProductCode,
		arg.Active,
		arg.Title,
		arg.BodyHtml,
		arg.Category,
		arg.Vendor,
		arg.ProductType,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i UpsertProductRow
	err := row.Scan(
		&i.ID,
		&i.Active,
		&i.ProductCode,
		&i.Title,
		&i.BodyHtml,
		&i.Category,
		&i.Vendor,
		&i.ProductType,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Inserted,
	)
	return i, err
}
