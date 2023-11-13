// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: users.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    id,
    name,
    email,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5)
RETURNING id, webhook_token, created_at, updated_at, name, email, api_key
`

type CreateUserParams struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.WebhookToken,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Email,
		&i.ApiKey,
	)
	return i, err
}

const getUserByApiKey = `-- name: GetUserByApiKey :one
SELECT id, webhook_token, created_at, updated_at, name, email, api_key FROM users
WHERE api_key = $1
LIMIT 1
`

func (q *Queries) GetUserByApiKey(ctx context.Context, apiKey string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByApiKey, apiKey)
	var i User
	err := row.Scan(
		&i.ID,
		&i.WebhookToken,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Email,
		&i.ApiKey,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, webhook_token, created_at, updated_at, name, email, api_key FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.WebhookToken,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Email,
		&i.ApiKey,
	)
	return i, err
}

const getUserByName = `-- name: GetUserByName :one
SELECT
    name
FROM users
WHERE name = $1
LIMIT 1
`

func (q *Queries) GetUserByName(ctx context.Context, name string) (string, error) {
	row := q.db.QueryRowContext(ctx, getUserByName, name)
	err := row.Scan(&name)
	return name, err
}

const getUsers = `-- name: GetUsers :one
SELECT id, webhook_token, created_at, updated_at, name, email, api_key FROM users LIMIT 1
`

func (q *Queries) GetUsers(ctx context.Context) (User, error) {
	row := q.db.QueryRowContext(ctx, getUsers)
	var i User
	err := row.Scan(
		&i.ID,
		&i.WebhookToken,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Email,
		&i.ApiKey,
	)
	return i, err
}

const removeUser = `-- name: RemoveUser :exec
DELETE FROM users
WHERE api_key = $1
`

func (q *Queries) RemoveUser(ctx context.Context, apiKey string) error {
	_, err := q.db.ExecContext(ctx, removeUser, apiKey)
	return err
}

const updateUser = `-- name: UpdateUser :execresult
UPDATE users 
SET
    name = $1,
    email = $2,
    updated_at = $3
WHERE id = $4
`

type UpdateUserParams struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updated_at"`
	ID        uuid.UUID `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateUser,
		arg.Name,
		arg.Email,
		arg.UpdatedAt,
		arg.ID,
	)
}

const validateWebhookByUser = `-- name: ValidateWebhookByUser :one
SELECT
    name
FROM users
WHERE 
webhook_token = $1 AND api_key = $2
`

type ValidateWebhookByUserParams struct {
	WebhookToken string `json:"webhook_token"`
	ApiKey       string `json:"api_key"`
}

func (q *Queries) ValidateWebhookByUser(ctx context.Context, arg ValidateWebhookByUserParams) (string, error) {
	row := q.db.QueryRowContext(ctx, validateWebhookByUser, arg.WebhookToken, arg.ApiKey)
	var name string
	err := row.Scan(&name)
	return name, err
}
