// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
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
    "name",
    user_type,
    email,
    "password",
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, name, email, user_type, password, api_key, webhook_token, created_at, updated_at
`

type CreateUserParams struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	UserType  string    `json:"user_type"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Name,
		arg.UserType,
		arg.Email,
		arg.Password,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.UserType,
		&i.Password,
		&i.ApiKey,
		&i.WebhookToken,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getApiKeyByCookieSecret = `-- name: GetApiKeyByCookieSecret :one
SELECT users.id, name, users.email, user_type, password, api_key, webhook_token, users.created_at, users.updated_at, google_oauth.id, user_id, cookie_secret, google_id, google_oauth.email, picture, google_oauth.created_at, google_oauth.updated_at FROM users
INNER JOIN google_oauth
ON users.id = google_oauth.user_id
WHERE google_oauth.cookie_secret = $1
`

type GetApiKeyByCookieSecretRow struct {
	ID           uuid.UUID      `json:"id"`
	Name         string         `json:"name"`
	Email        string         `json:"email"`
	UserType     string         `json:"user_type"`
	Password     string         `json:"password"`
	ApiKey       string         `json:"api_key"`
	WebhookToken string         `json:"webhook_token"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	ID_2         uuid.UUID      `json:"id_2"`
	UserID       uuid.UUID      `json:"user_id"`
	CookieSecret string         `json:"cookie_secret"`
	GoogleID     string         `json:"google_id"`
	Email_2      string         `json:"email_2"`
	Picture      sql.NullString `json:"picture"`
	CreatedAt_2  time.Time      `json:"created_at_2"`
	UpdatedAt_2  time.Time      `json:"updated_at_2"`
}

func (q *Queries) GetApiKeyByCookieSecret(ctx context.Context, cookieSecret string) (GetApiKeyByCookieSecretRow, error) {
	row := q.db.QueryRowContext(ctx, getApiKeyByCookieSecret, cookieSecret)
	var i GetApiKeyByCookieSecretRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.UserType,
		&i.Password,
		&i.ApiKey,
		&i.WebhookToken,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ID_2,
		&i.UserID,
		&i.CookieSecret,
		&i.GoogleID,
		&i.Email_2,
		&i.Picture,
		&i.CreatedAt_2,
		&i.UpdatedAt_2,
	)
	return i, err
}

const getUserByApiKey = `-- name: GetUserByApiKey :one
SELECT id, name, email, user_type, password, api_key, webhook_token, created_at, updated_at FROM users
WHERE api_key = $1
LIMIT 1
`

func (q *Queries) GetUserByApiKey(ctx context.Context, apiKey string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByApiKey, apiKey)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.UserType,
		&i.Password,
		&i.ApiKey,
		&i.WebhookToken,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, name, email, user_type, password, api_key, webhook_token, created_at, updated_at FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.UserType,
		&i.Password,
		&i.ApiKey,
		&i.WebhookToken,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmailType = `-- name: GetUserByEmailType :one
SELECT
    email
FROM users
WHERE email = $1 AND user_type = $2
LIMIT 1
`

type GetUserByEmailTypeParams struct {
	Email    string `json:"email"`
	UserType string `json:"user_type"`
}

func (q *Queries) GetUserByEmailType(ctx context.Context, arg GetUserByEmailTypeParams) (string, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmailType, arg.Email, arg.UserType)
	var email string
	err := row.Scan(&email)
	return email, err
}

const getUserByName = `-- name: GetUserByName :one
SELECT
    "name"
FROM users
WHERE "name" = $1
LIMIT 1
`

func (q *Queries) GetUserByName(ctx context.Context, name string) (string, error) {
	row := q.db.QueryRowContext(ctx, getUserByName, name)
	err := row.Scan(&name)
	return name, err
}

const getUserCredentials = `-- name: GetUserCredentials :one
SELECT
    "name",
    "password",
    api_key
FROM users
WHERE "name" = $1
AND "password" = $2
LIMIT 1
`

type GetUserCredentialsParams struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type GetUserCredentialsRow struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	ApiKey   string `json:"api_key"`
}

func (q *Queries) GetUserCredentials(ctx context.Context, arg GetUserCredentialsParams) (GetUserCredentialsRow, error) {
	row := q.db.QueryRowContext(ctx, getUserCredentials, arg.Name, arg.Password)
	var i GetUserCredentialsRow
	err := row.Scan(&i.Name, &i.Password, &i.ApiKey)
	return i, err
}

const getUsers = `-- name: GetUsers :one
SELECT id, name, email, user_type, password, api_key, webhook_token, created_at, updated_at FROM users LIMIT 1
`

func (q *Queries) GetUsers(ctx context.Context) (User, error) {
	row := q.db.QueryRowContext(ctx, getUsers)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.UserType,
		&i.Password,
		&i.ApiKey,
		&i.WebhookToken,
		&i.CreatedAt,
		&i.UpdatedAt,
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
    "name" = $1,
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
    "name"
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
