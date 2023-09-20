// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: orderlines.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createOrderLine = `-- name: CreateOrderLine :one
INSERT INTO order_lines(
    id,
    order_id,
    line_type,
    sku,
    price,
    barcode,
    qty,
    tax_rate,
    tax_total,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING id, order_id, line_type, sku, price, barcode, qty, tax_total, tax_rate, created_at, updated_at
`

type CreateOrderLineParams struct {
	ID        uuid.UUID      `json:"id"`
	OrderID   uuid.UUID      `json:"order_id"`
	LineType  sql.NullString `json:"line_type"`
	Sku       string         `json:"sku"`
	Price     sql.NullString `json:"price"`
	Barcode   sql.NullInt32  `json:"barcode"`
	Qty       sql.NullInt32  `json:"qty"`
	TaxRate   sql.NullString `json:"tax_rate"`
	TaxTotal  sql.NullString `json:"tax_total"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (q *Queries) CreateOrderLine(ctx context.Context, arg CreateOrderLineParams) (OrderLine, error) {
	row := q.db.QueryRowContext(ctx, createOrderLine,
		arg.ID,
		arg.OrderID,
		arg.LineType,
		arg.Sku,
		arg.Price,
		arg.Barcode,
		arg.Qty,
		arg.TaxRate,
		arg.TaxTotal,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i OrderLine
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.LineType,
		&i.Sku,
		&i.Price,
		&i.Barcode,
		&i.Qty,
		&i.TaxTotal,
		&i.TaxRate,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getOrderLinesByOrder = `-- name: GetOrderLinesByOrder :many
SELECT
    sku,
    line_type,
    price,
    barcode,
    qty,
    tax_rate,
    tax_total,
    updated_at
FROM order_lines
WHERE order_id = $1 AND line_type = 'line'
`

type GetOrderLinesByOrderRow struct {
	Sku       string         `json:"sku"`
	LineType  sql.NullString `json:"line_type"`
	Price     sql.NullString `json:"price"`
	Barcode   sql.NullInt32  `json:"barcode"`
	Qty       sql.NullInt32  `json:"qty"`
	TaxRate   sql.NullString `json:"tax_rate"`
	TaxTotal  sql.NullString `json:"tax_total"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (q *Queries) GetOrderLinesByOrder(ctx context.Context, orderID uuid.UUID) ([]GetOrderLinesByOrderRow, error) {
	rows, err := q.db.QueryContext(ctx, getOrderLinesByOrder, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOrderLinesByOrderRow
	for rows.Next() {
		var i GetOrderLinesByOrderRow
		if err := rows.Scan(
			&i.Sku,
			&i.LineType,
			&i.Price,
			&i.Barcode,
			&i.Qty,
			&i.TaxRate,
			&i.TaxTotal,
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

const getShippingLinesByOrder = `-- name: GetShippingLinesByOrder :many
SELECT
    sku,
    line_type,
    price,
    barcode,
    qty,
    tax_rate,
    tax_total,
    updated_at
FROM order_lines
WHERE order_id = $1 AND line_type = 'shipping'
`

type GetShippingLinesByOrderRow struct {
	Sku       string         `json:"sku"`
	LineType  sql.NullString `json:"line_type"`
	Price     sql.NullString `json:"price"`
	Barcode   sql.NullInt32  `json:"barcode"`
	Qty       sql.NullInt32  `json:"qty"`
	TaxRate   sql.NullString `json:"tax_rate"`
	TaxTotal  sql.NullString `json:"tax_total"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (q *Queries) GetShippingLinesByOrder(ctx context.Context, orderID uuid.UUID) ([]GetShippingLinesByOrderRow, error) {
	rows, err := q.db.QueryContext(ctx, getShippingLinesByOrder, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetShippingLinesByOrderRow
	for rows.Next() {
		var i GetShippingLinesByOrderRow
		if err := rows.Scan(
			&i.Sku,
			&i.LineType,
			&i.Price,
			&i.Barcode,
			&i.Qty,
			&i.TaxRate,
			&i.TaxTotal,
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

const updateOrderLine = `-- name: UpdateOrderLine :one
UPDATE order_lines
SET
    order_id = $1,
    line_type = $2,
    sku = $3,
    price = $4,
    barcode = $5,
    qty = $6,
    tax_rate = $7,
    tax_total = $8,
    created_at = $9,
    updated_at = $10
WHERE id = $11
RETURNING id, order_id, line_type, sku, price, barcode, qty, tax_total, tax_rate, created_at, updated_at
`

type UpdateOrderLineParams struct {
	OrderID   uuid.UUID      `json:"order_id"`
	LineType  sql.NullString `json:"line_type"`
	Sku       string         `json:"sku"`
	Price     sql.NullString `json:"price"`
	Barcode   sql.NullInt32  `json:"barcode"`
	Qty       sql.NullInt32  `json:"qty"`
	TaxRate   sql.NullString `json:"tax_rate"`
	TaxTotal  sql.NullString `json:"tax_total"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	ID        uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateOrderLine(ctx context.Context, arg UpdateOrderLineParams) (OrderLine, error) {
	row := q.db.QueryRowContext(ctx, updateOrderLine,
		arg.OrderID,
		arg.LineType,
		arg.Sku,
		arg.Price,
		arg.Barcode,
		arg.Qty,
		arg.TaxRate,
		arg.TaxTotal,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.ID,
	)
	var i OrderLine
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.LineType,
		&i.Sku,
		&i.Price,
		&i.Barcode,
		&i.Qty,
		&i.TaxTotal,
		&i.TaxRate,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
