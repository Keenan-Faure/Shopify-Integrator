// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: address.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createAddress = `-- name: CreateAddress :one
INSERT INTO address(
    id,
    customer_id,
    name,
    first_name,
    last_name,
    address1,
    address2,
    suburb,
    city,
    province,
    postal_code,
    company,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
)
RETURNING id, customer_id, name, first_name, last_name, address1, address2, suburb, city, province, postal_code, company, created_at, updated_at
`

type CreateAddressParams struct {
	ID         uuid.UUID      `json:"id"`
	CustomerID uuid.UUID      `json:"customer_id"`
	Name       sql.NullString `json:"name"`
	FirstName  string         `json:"first_name"`
	LastName   string         `json:"last_name"`
	Address1   sql.NullString `json:"address1"`
	Address2   sql.NullString `json:"address2"`
	Suburb     sql.NullString `json:"suburb"`
	City       sql.NullString `json:"city"`
	Province   sql.NullString `json:"province"`
	PostalCode sql.NullString `json:"postal_code"`
	Company    sql.NullString `json:"company"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func (q *Queries) CreateAddress(ctx context.Context, arg CreateAddressParams) (Address, error) {
	row := q.db.QueryRowContext(ctx, createAddress,
		arg.ID,
		arg.CustomerID,
		arg.Name,
		arg.FirstName,
		arg.LastName,
		arg.Address1,
		arg.Address2,
		arg.Suburb,
		arg.City,
		arg.Province,
		arg.PostalCode,
		arg.Company,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Address
	err := row.Scan(
		&i.ID,
		&i.CustomerID,
		&i.Name,
		&i.FirstName,
		&i.LastName,
		&i.Address1,
		&i.Address2,
		&i.Suburb,
		&i.City,
		&i.Province,
		&i.PostalCode,
		&i.Company,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAddressByCustomer = `-- name: GetAddressByCustomer :many
SELECT
    id,
    first_name,
    last_name,
    address1,
    address2,
    suburb,
    city,
    province,
    postal_code,
    company,
    updated_at
FROM address
WHERE customer_id = $1
`

type GetAddressByCustomerRow struct {
	ID         uuid.UUID      `json:"id"`
	FirstName  string         `json:"first_name"`
	LastName   string         `json:"last_name"`
	Address1   sql.NullString `json:"address1"`
	Address2   sql.NullString `json:"address2"`
	Suburb     sql.NullString `json:"suburb"`
	City       sql.NullString `json:"city"`
	Province   sql.NullString `json:"province"`
	PostalCode sql.NullString `json:"postal_code"`
	Company    sql.NullString `json:"company"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func (q *Queries) GetAddressByCustomer(ctx context.Context, customerID uuid.UUID) ([]GetAddressByCustomerRow, error) {
	rows, err := q.db.QueryContext(ctx, getAddressByCustomer, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAddressByCustomerRow
	for rows.Next() {
		var i GetAddressByCustomerRow
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Address1,
			&i.Address2,
			&i.Suburb,
			&i.City,
			&i.Province,
			&i.PostalCode,
			&i.Company,
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

const removeAddress = `-- name: RemoveAddress :exec
DELETE FROM address
WHERE id = $1
`

func (q *Queries) RemoveAddress(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, removeAddress, id)
	return err
}

const updateAddress = `-- name: UpdateAddress :one
UPDATE address
SET
    customer_id = $1,
    first_name = $2,
    last_name = $3,
    address1 = $4,
    address2 = $5,
    suburb = $6,
    city = $7,
    province = $8,
    postal_code = $9,
    company = $10,
    updated_at = $11
WHERE id = $12
RETURNING id, customer_id, name, first_name, last_name, address1, address2, suburb, city, province, postal_code, company, created_at, updated_at
`

type UpdateAddressParams struct {
	CustomerID uuid.UUID      `json:"customer_id"`
	FirstName  string         `json:"first_name"`
	LastName   string         `json:"last_name"`
	Address1   sql.NullString `json:"address1"`
	Address2   sql.NullString `json:"address2"`
	Suburb     sql.NullString `json:"suburb"`
	City       sql.NullString `json:"city"`
	Province   sql.NullString `json:"province"`
	PostalCode sql.NullString `json:"postal_code"`
	Company    sql.NullString `json:"company"`
	UpdatedAt  time.Time      `json:"updated_at"`
	ID         uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateAddress(ctx context.Context, arg UpdateAddressParams) (Address, error) {
	row := q.db.QueryRowContext(ctx, updateAddress,
		arg.CustomerID,
		arg.FirstName,
		arg.LastName,
		arg.Address1,
		arg.Address2,
		arg.Suburb,
		arg.City,
		arg.Province,
		arg.PostalCode,
		arg.Company,
		arg.UpdatedAt,
		arg.ID,
	)
	var i Address
	err := row.Scan(
		&i.ID,
		&i.CustomerID,
		&i.Name,
		&i.FirstName,
		&i.LastName,
		&i.Address1,
		&i.Address2,
		&i.Suburb,
		&i.City,
		&i.Province,
		&i.PostalCode,
		&i.Company,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
