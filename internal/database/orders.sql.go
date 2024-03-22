// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: orders.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createOrder = `-- name: CreateOrder :one
INSERT INTO orders(
    id,
    status,
    notes,
    web_code,
    tax_total,
    order_total,
    shipping_total,
    discount_total,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING id, notes, web_code, tax_total, order_total, shipping_total, discount_total, created_at, updated_at, status
`

type CreateOrderParams struct {
	ID            uuid.UUID      `json:"id"`
	Status        string         `json:"status"`
	Notes         sql.NullString `json:"notes"`
	WebCode       string         `json:"web_code"`
	TaxTotal      sql.NullString `json:"tax_total"`
	OrderTotal    sql.NullString `json:"order_total"`
	ShippingTotal sql.NullString `json:"shipping_total"`
	DiscountTotal sql.NullString `json:"discount_total"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, createOrder,
		arg.ID,
		arg.Status,
		arg.Notes,
		arg.WebCode,
		arg.TaxTotal,
		arg.OrderTotal,
		arg.ShippingTotal,
		arg.DiscountTotal,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.Notes,
		&i.WebCode,
		&i.TaxTotal,
		&i.OrderTotal,
		&i.ShippingTotal,
		&i.DiscountTotal,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Status,
	)
	return i, err
}

const fetchOrderStatsNotPaid = `-- name: FetchOrderStatsNotPaid :many
SELECT
	COUNT(id) AS "count",
	to_char(created_at, 'YYYY-MM-DD') AS "day"
FROM orders
WHERE
    created_at > current_date at time zone 'UTC' - interval '7 day' AND
    "status" != 'paid'
GROUP BY "day"
ORDER BY "day" DESC
`

type FetchOrderStatsNotPaidRow struct {
	Count int64  `json:"count"`
	Day   string `json:"day"`
}

func (q *Queries) FetchOrderStatsNotPaid(ctx context.Context) ([]FetchOrderStatsNotPaidRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchOrderStatsNotPaid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FetchOrderStatsNotPaidRow
	for rows.Next() {
		var i FetchOrderStatsNotPaidRow
		if err := rows.Scan(&i.Count, &i.Day); err != nil {
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

const fetchOrderStatsPaid = `-- name: FetchOrderStatsPaid :many
SELECT
	COUNT(id) AS "count",
	to_char(created_at, 'YYYY-MM-DD') AS "day"
FROM orders
WHERE
    created_at > current_date at time zone 'UTC' - interval '7 day' AND
    "status" = 'paid'
GROUP BY "day"
ORDER BY "day" DESC
`

type FetchOrderStatsPaidRow struct {
	Count int64  `json:"count"`
	Day   string `json:"day"`
}

func (q *Queries) FetchOrderStatsPaid(ctx context.Context) ([]FetchOrderStatsPaidRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchOrderStatsPaid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FetchOrderStatsPaidRow
	for rows.Next() {
		var i FetchOrderStatsPaidRow
		if err := rows.Scan(&i.Count, &i.Day); err != nil {
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

const getOrderByCustomer = `-- name: GetOrderByCustomer :many
SELECT
    o.id,
    o.notes,
    o.status,
    o.web_code,
    o.tax_total,
    o.order_total,
    o.shipping_total,
    o.discount_total,
    o.updated_at
FROM orders o 
WHERE orders.id in (
    SELECT order_id FROM customer_orders
    WHERE customer_id = $1
)
`

type GetOrderByCustomerRow struct {
	ID            uuid.UUID      `json:"id"`
	Notes         sql.NullString `json:"notes"`
	Status        string         `json:"status"`
	WebCode       string         `json:"web_code"`
	TaxTotal      sql.NullString `json:"tax_total"`
	OrderTotal    sql.NullString `json:"order_total"`
	ShippingTotal sql.NullString `json:"shipping_total"`
	DiscountTotal sql.NullString `json:"discount_total"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (q *Queries) GetOrderByCustomer(ctx context.Context, customerID uuid.UUID) ([]GetOrderByCustomerRow, error) {
	rows, err := q.db.QueryContext(ctx, getOrderByCustomer, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOrderByCustomerRow
	for rows.Next() {
		var i GetOrderByCustomerRow
		if err := rows.Scan(
			&i.ID,
			&i.Notes,
			&i.Status,
			&i.WebCode,
			&i.TaxTotal,
			&i.OrderTotal,
			&i.ShippingTotal,
			&i.DiscountTotal,
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

const getOrderByID = `-- name: GetOrderByID :one
SELECT
    id,
    notes,
    status,
    web_code,
    tax_total,
    order_total,
    shipping_total,
    discount_total,
    updated_at,
    created_at
FROM orders
WHERE id = $1
`

type GetOrderByIDRow struct {
	ID            uuid.UUID      `json:"id"`
	Notes         sql.NullString `json:"notes"`
	Status        string         `json:"status"`
	WebCode       string         `json:"web_code"`
	TaxTotal      sql.NullString `json:"tax_total"`
	OrderTotal    sql.NullString `json:"order_total"`
	ShippingTotal sql.NullString `json:"shipping_total"`
	DiscountTotal sql.NullString `json:"discount_total"`
	UpdatedAt     time.Time      `json:"updated_at"`
	CreatedAt     time.Time      `json:"created_at"`
}

func (q *Queries) GetOrderByID(ctx context.Context, id uuid.UUID) (GetOrderByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getOrderByID, id)
	var i GetOrderByIDRow
	err := row.Scan(
		&i.ID,
		&i.Notes,
		&i.Status,
		&i.WebCode,
		&i.TaxTotal,
		&i.OrderTotal,
		&i.ShippingTotal,
		&i.DiscountTotal,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getOrderByWebCode = `-- name: GetOrderByWebCode :one
SELECT
    id,
    notes,
    status,
    web_code,
    tax_total,
    order_total,
    shipping_total,
    discount_total,
    updated_at
FROM orders
WHERE web_code = $1
`

type GetOrderByWebCodeRow struct {
	ID            uuid.UUID      `json:"id"`
	Notes         sql.NullString `json:"notes"`
	Status        string         `json:"status"`
	WebCode       string         `json:"web_code"`
	TaxTotal      sql.NullString `json:"tax_total"`
	OrderTotal    sql.NullString `json:"order_total"`
	ShippingTotal sql.NullString `json:"shipping_total"`
	DiscountTotal sql.NullString `json:"discount_total"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (q *Queries) GetOrderByWebCode(ctx context.Context, webCode string) (GetOrderByWebCodeRow, error) {
	row := q.db.QueryRowContext(ctx, getOrderByWebCode, webCode)
	var i GetOrderByWebCodeRow
	err := row.Scan(
		&i.ID,
		&i.Notes,
		&i.Status,
		&i.WebCode,
		&i.TaxTotal,
		&i.OrderTotal,
		&i.ShippingTotal,
		&i.DiscountTotal,
		&i.UpdatedAt,
	)
	return i, err
}

const getOrderIDByWebCode = `-- name: GetOrderIDByWebCode :one
SELECT
    id
FROM orders
WHERE web_code = $1
`

func (q *Queries) GetOrderIDByWebCode(ctx context.Context, webCode string) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getOrderIDByWebCode, webCode)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getOrders = `-- name: GetOrders :many
SELECT
    id,
    notes,
    status,
    web_code,
    tax_total,
    order_total,
    shipping_total,
    discount_total,
    updated_at
FROM orders
ORDER BY updated_at DESC
LIMIT $1 OFFSET $2
`

type GetOrdersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetOrdersRow struct {
	ID            uuid.UUID      `json:"id"`
	Notes         sql.NullString `json:"notes"`
	Status        string         `json:"status"`
	WebCode       string         `json:"web_code"`
	TaxTotal      sql.NullString `json:"tax_total"`
	OrderTotal    sql.NullString `json:"order_total"`
	ShippingTotal sql.NullString `json:"shipping_total"`
	DiscountTotal sql.NullString `json:"discount_total"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (q *Queries) GetOrders(ctx context.Context, arg GetOrdersParams) ([]GetOrdersRow, error) {
	rows, err := q.db.QueryContext(ctx, getOrders, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOrdersRow
	for rows.Next() {
		var i GetOrdersRow
		if err := rows.Scan(
			&i.ID,
			&i.Notes,
			&i.Status,
			&i.WebCode,
			&i.TaxTotal,
			&i.OrderTotal,
			&i.ShippingTotal,
			&i.DiscountTotal,
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

const getOrdersSearchByCustomer = `-- name: GetOrdersSearchByCustomer :many
SELECT
    o.id,
    o.notes,
    o.status,
    o.web_code,
    o.tax_total,
    o.order_total,
    o.shipping_total,
    o.discount_total,
    o.updated_at
FROM orders o
WHERE o.id in (
    SELECT order_id FROM customer_orders co
    INNER JOIN customers c
    ON co.customer_id = c.id
    WHERE CONCAT(c.first_name, ' ', c.last_name) SIMILAR TO $1
    OR c.first_name LIKE $1
    OR c.last_name LIKE $1
)
`

type GetOrdersSearchByCustomerRow struct {
	ID            uuid.UUID      `json:"id"`
	Notes         sql.NullString `json:"notes"`
	Status        string         `json:"status"`
	WebCode       string         `json:"web_code"`
	TaxTotal      sql.NullString `json:"tax_total"`
	OrderTotal    sql.NullString `json:"order_total"`
	ShippingTotal sql.NullString `json:"shipping_total"`
	DiscountTotal sql.NullString `json:"discount_total"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (q *Queries) GetOrdersSearchByCustomer(ctx context.Context, similarToEscape string) ([]GetOrdersSearchByCustomerRow, error) {
	rows, err := q.db.QueryContext(ctx, getOrdersSearchByCustomer, similarToEscape)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOrdersSearchByCustomerRow
	for rows.Next() {
		var i GetOrdersSearchByCustomerRow
		if err := rows.Scan(
			&i.ID,
			&i.Notes,
			&i.Status,
			&i.WebCode,
			&i.TaxTotal,
			&i.OrderTotal,
			&i.ShippingTotal,
			&i.DiscountTotal,
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

const getOrdersSearchWebCode = `-- name: GetOrdersSearchWebCode :many
SELECT
    id,
    notes,
    status,
    web_code,
    tax_total,
    order_total,
    shipping_total,
    discount_total,
    updated_at
FROM orders
WHERE web_code SIMILAR TO $1
LIMIT 10
`

type GetOrdersSearchWebCodeRow struct {
	ID            uuid.UUID      `json:"id"`
	Notes         sql.NullString `json:"notes"`
	Status        string         `json:"status"`
	WebCode       string         `json:"web_code"`
	TaxTotal      sql.NullString `json:"tax_total"`
	OrderTotal    sql.NullString `json:"order_total"`
	ShippingTotal sql.NullString `json:"shipping_total"`
	DiscountTotal sql.NullString `json:"discount_total"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (q *Queries) GetOrdersSearchWebCode(ctx context.Context, similarToEscape string) ([]GetOrdersSearchWebCodeRow, error) {
	rows, err := q.db.QueryContext(ctx, getOrdersSearchWebCode, similarToEscape)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOrdersSearchWebCodeRow
	for rows.Next() {
		var i GetOrdersSearchWebCodeRow
		if err := rows.Scan(
			&i.ID,
			&i.Notes,
			&i.Status,
			&i.WebCode,
			&i.TaxTotal,
			&i.OrderTotal,
			&i.ShippingTotal,
			&i.DiscountTotal,
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

const removeOrder = `-- name: RemoveOrder :exec
DELETE FROM orders
WHERE id = $1
`

func (q *Queries) RemoveOrder(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, removeOrder, id)
	return err
}

const removeOrderByWebCode = `-- name: RemoveOrderByWebCode :exec
DELETE FROM orders
WHERE web_code = $1
`

func (q *Queries) RemoveOrderByWebCode(ctx context.Context, webCode string) error {
	_, err := q.db.ExecContext(ctx, removeOrderByWebCode, webCode)
	return err
}

const updateOrder = `-- name: UpdateOrder :one
UPDATE orders
SET
    notes = $1,
    status = $2,
    web_code = $3,
    tax_total = $4,
    order_total = $5,
    shipping_total = $6,
    discount_total = $7,
    updated_at = $8
WHERE id = $9
RETURNING id, notes, web_code, tax_total, order_total, shipping_total, discount_total, created_at, updated_at, status
`

type UpdateOrderParams struct {
	Notes         sql.NullString `json:"notes"`
	Status        string         `json:"status"`
	WebCode       string         `json:"web_code"`
	TaxTotal      sql.NullString `json:"tax_total"`
	OrderTotal    sql.NullString `json:"order_total"`
	ShippingTotal sql.NullString `json:"shipping_total"`
	DiscountTotal sql.NullString `json:"discount_total"`
	UpdatedAt     time.Time      `json:"updated_at"`
	ID            uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, updateOrder,
		arg.Notes,
		arg.Status,
		arg.WebCode,
		arg.TaxTotal,
		arg.OrderTotal,
		arg.ShippingTotal,
		arg.DiscountTotal,
		arg.UpdatedAt,
		arg.ID,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.Notes,
		&i.WebCode,
		&i.TaxTotal,
		&i.OrderTotal,
		&i.ShippingTotal,
		&i.DiscountTotal,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Status,
	)
	return i, err
}
