// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: register_tokens.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createToken = `-- name: CreateToken :one
INSERT INTO register_tokens(
    id,
    name,
    email,
    token,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id, name, email, token, created_at, updated_at
`

type CreateTokenParams struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Token     uuid.UUID `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) CreateToken(ctx context.Context, arg CreateTokenParams) (RegisterToken, error) {
	row := q.db.QueryRowContext(ctx, createToken,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Token,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i RegisterToken
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteToken = `-- name: DeleteToken :exec
DELETE FROM register_tokens
WHERE
token = $1 AND
email = $2
`

type DeleteTokenParams struct {
	Token uuid.UUID `json:"token"`
	Email string    `json:"email"`
}

func (q *Queries) DeleteToken(ctx context.Context, arg DeleteTokenParams) error {
	_, err := q.db.ExecContext(ctx, deleteToken, arg.Token, arg.Email)
	return err
}

const getToken = `-- name: GetToken :one
SELECT
    name,
    email,
    token
FROM register_tokens
WHERE email = $1
`

type GetTokenRow struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Token uuid.UUID `json:"token"`
}

func (q *Queries) GetToken(ctx context.Context, email string) (GetTokenRow, error) {
	row := q.db.QueryRowContext(ctx, getToken, email)
	var i GetTokenRow
	err := row.Scan(&i.Name, &i.Email, &i.Token)
	return i, err
}

const getTokenValidation = `-- name: GetTokenValidation :one
SELECT
    name,
    email,
    token
FROM register_tokens
WHERE email = $1
`

type GetTokenValidationRow struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Token uuid.UUID `json:"token"`
}

func (q *Queries) GetTokenValidation(ctx context.Context, email string) (GetTokenValidationRow, error) {
	row := q.db.QueryRowContext(ctx, getTokenValidation, email)
	var i GetTokenValidationRow
	err := row.Scan(&i.Name, &i.Email, &i.Token)
	return i, err
}

const updateToken = `-- name: UpdateToken :one
UPDATE register_tokens
SET
   token = $1
where email = $2
RETURNING id, name, email, token, created_at, updated_at
`

type UpdateTokenParams struct {
	Token uuid.UUID `json:"token"`
	Email string    `json:"email"`
}

func (q *Queries) UpdateToken(ctx context.Context, arg UpdateTokenParams) (RegisterToken, error) {
	row := q.db.QueryRowContext(ctx, updateToken, arg.Token, arg.Email)
	var i RegisterToken
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
