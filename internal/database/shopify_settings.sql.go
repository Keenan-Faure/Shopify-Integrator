// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: shopify_settings.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const addShopifySetting = `-- name: AddShopifySetting :exec
INSERT INTO shopify_settings(
    id,
    key,
    description,
    value,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)
`

type AddShopifySettingParams struct {
	ID          uuid.UUID `json:"id"`
	Key         string    `json:"key"`
	Description string    `json:"description"`
	Value       string    `json:"value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (q *Queries) AddShopifySetting(ctx context.Context, arg AddShopifySettingParams) error {
	_, err := q.db.ExecContext(ctx, addShopifySetting,
		arg.ID,
		arg.Key,
		arg.Description,
		arg.Value,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const getShopifySettingByKey = `-- name: GetShopifySettingByKey :one
SELECT
    id,
    key,
    description,
    value,
    updated_at
FROM shopify_settings
WHERE key = $1
`

type GetShopifySettingByKeyRow struct {
	ID          uuid.UUID `json:"id"`
	Key         string    `json:"key"`
	Description string    `json:"description"`
	Value       string    `json:"value"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (q *Queries) GetShopifySettingByKey(ctx context.Context, key string) (GetShopifySettingByKeyRow, error) {
	row := q.db.QueryRowContext(ctx, getShopifySettingByKey, key)
	var i GetShopifySettingByKeyRow
	err := row.Scan(
		&i.ID,
		&i.Key,
		&i.Description,
		&i.Value,
		&i.UpdatedAt,
	)
	return i, err
}

const getShopifySettings = `-- name: GetShopifySettings :many
SELECT
    id,
    key,
    description,
    value,
    updated_at
FROM shopify_settings
`

type GetShopifySettingsRow struct {
	ID          uuid.UUID `json:"id"`
	Key         string    `json:"key"`
	Description string    `json:"description"`
	Value       string    `json:"value"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (q *Queries) GetShopifySettings(ctx context.Context) ([]GetShopifySettingsRow, error) {
	rows, err := q.db.QueryContext(ctx, getShopifySettings)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetShopifySettingsRow
	for rows.Next() {
		var i GetShopifySettingsRow
		if err := rows.Scan(
			&i.ID,
			&i.Key,
			&i.Description,
			&i.Value,
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

const removeShopifySetting = `-- name: RemoveShopifySetting :exec
DELETE FROM shopify_settings
WHERE key = $1
`

func (q *Queries) RemoveShopifySetting(ctx context.Context, key string) error {
	_, err := q.db.ExecContext(ctx, removeShopifySetting, key)
	return err
}

const updateShopifySetting = `-- name: UpdateShopifySetting :exec
UPDATE shopify_settings
SET
    value = $1,
    updated_at = $2
WHERE key = $3
`

type UpdateShopifySettingParams struct {
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
	Key       string    `json:"key"`
}

func (q *Queries) UpdateShopifySetting(ctx context.Context, arg UpdateShopifySettingParams) error {
	_, err := q.db.ExecContext(ctx, updateShopifySetting, arg.Value, arg.UpdatedAt, arg.Key)
	return err
}