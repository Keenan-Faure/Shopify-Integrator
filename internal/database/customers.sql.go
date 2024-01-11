// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: customers.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createCustomer = `-- name: CreateCustomer :one
INSERT INTO customers(
    id,
    first_name,
    last_name,
    email,
    phone,
    created_at,
    updated_at
) VALUES(
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, first_name, last_name, email, phone, created_at, updated_at
`

type CreateCustomerParams struct {
	ID        uuid.UUID      `json:"id"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Email     sql.NullString `json:"email"`
	Phone     sql.NullString `json:"phone"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (q *Queries) CreateCustomer(ctx context.Context, arg CreateCustomerParams) (Customer, error) {
	row := q.db.QueryRowContext(ctx, createCustomer,
		arg.ID,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Phone,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Customer
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Phone,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getCustomerByID = `-- name: GetCustomerByID :one
SELECT
    id,
    first_name,
    last_name,
    email,
    phone,
    updated_at
FROM customers
WHERE id = $1
`

type GetCustomerByIDRow struct {
	ID        uuid.UUID      `json:"id"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Email     sql.NullString `json:"email"`
	Phone     sql.NullString `json:"phone"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (q *Queries) GetCustomerByID(ctx context.Context, id uuid.UUID) (GetCustomerByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getCustomerByID, id)
	var i GetCustomerByIDRow
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Phone,
		&i.UpdatedAt,
	)
	return i, err
}

const getCustomers = `-- name: GetCustomers :many
SELECT
    id,
    first_name,
    last_name,
    email,
    phone,
    updated_at
FROM customers
LIMIT $1 OFFSET $2
`

type GetCustomersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetCustomersRow struct {
	ID        uuid.UUID      `json:"id"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Email     sql.NullString `json:"email"`
	Phone     sql.NullString `json:"phone"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (q *Queries) GetCustomers(ctx context.Context, arg GetCustomersParams) ([]GetCustomersRow, error) {
	rows, err := q.db.QueryContext(ctx, getCustomers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCustomersRow
	for rows.Next() {
		var i GetCustomersRow
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
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

const getCustomersByName = `-- name: GetCustomersByName :many
SELECT
    id,
    first_name,
    last_name,
    email,
    phone,
    updated_at
FROM customers
WHERE CONCAT(first_name, ' ', last_name) SIMILAR TO $1
AND first_name LIKE $1
AND last_name LIKE $1
LIMIT 10
`

type GetCustomersByNameRow struct {
	ID        uuid.UUID      `json:"id"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Email     sql.NullString `json:"email"`
	Phone     sql.NullString `json:"phone"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (q *Queries) GetCustomersByName(ctx context.Context, similarToEscape string) ([]GetCustomersByNameRow, error) {
	rows, err := q.db.QueryContext(ctx, getCustomersByName, similarToEscape)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCustomersByNameRow
	for rows.Next() {
		var i GetCustomersByNameRow
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
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

const removeCustomer = `-- name: RemoveCustomer :exec
DELETE FROM customers
WHERE id = $1
`

func (q *Queries) RemoveCustomer(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, removeCustomer, id)
	return err
}

const updateCustomer = `-- name: UpdateCustomer :exec
UPDATE customers
SET
    first_name = $1,
    last_name = $2,
    email = $3,
    phone = $4,
    updated_at = $5
WHERE id = $6
`

type UpdateCustomerParams struct {
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Email     sql.NullString `json:"email"`
	Phone     sql.NullString `json:"phone"`
	UpdatedAt time.Time      `json:"updated_at"`
	ID        uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateCustomer(ctx context.Context, arg UpdateCustomerParams) error {
	_, err := q.db.ExecContext(ctx, updateCustomer,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Phone,
		arg.UpdatedAt,
		arg.ID,
	)
	return err
}
