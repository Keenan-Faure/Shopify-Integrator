// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: app_settings.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const addAppSetting = `-- name: AddAppSetting :exec
INSERT INTO app_settings(
    id,
    key,
    field_name,
    description,
    value,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
`

type AddAppSettingParams struct {
	ID          uuid.UUID `json:"id"`
	Key         string    `json:"key"`
	FieldName   string    `json:"field_name"`
	Description string    `json:"description"`
	Value       string    `json:"value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (q *Queries) AddAppSetting(ctx context.Context, arg AddAppSettingParams) error {
	_, err := q.db.ExecContext(ctx, addAppSetting,
		arg.ID,
		arg.Key,
		arg.FieldName,
		arg.Description,
		arg.Value,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const getAppSettingByKey = `-- name: GetAppSettingByKey :one
SELECT
    id,
    key,
    description,
    field_name,
    value,
    updated_at
FROM app_settings
WHERE key = $1
`

type GetAppSettingByKeyRow struct {
	ID          uuid.UUID `json:"id"`
	Key         string    `json:"key"`
	Description string    `json:"description"`
	FieldName   string    `json:"field_name"`
	Value       string    `json:"value"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (q *Queries) GetAppSettingByKey(ctx context.Context, key string) (GetAppSettingByKeyRow, error) {
	row := q.db.QueryRowContext(ctx, getAppSettingByKey, key)
	var i GetAppSettingByKeyRow
	err := row.Scan(
		&i.ID,
		&i.Key,
		&i.Description,
		&i.FieldName,
		&i.Value,
		&i.UpdatedAt,
	)
	return i, err
}

const getAppSettings = `-- name: GetAppSettings :many
SELECT
    id,
    key,
    description,
    field_name,
    value,
    updated_at
FROM app_settings
`

type GetAppSettingsRow struct {
	ID          uuid.UUID `json:"id"`
	Key         string    `json:"key"`
	Description string    `json:"description"`
	FieldName   string    `json:"field_name"`
	Value       string    `json:"value"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (q *Queries) GetAppSettings(ctx context.Context) ([]GetAppSettingsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAppSettings)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAppSettingsRow
	for rows.Next() {
		var i GetAppSettingsRow
		if err := rows.Scan(
			&i.ID,
			&i.Key,
			&i.Description,
			&i.FieldName,
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

const getAppSettingsList = `-- name: GetAppSettingsList :many
SELECT DISTINCT("key") FROM app_settings
`

func (q *Queries) GetAppSettingsList(ctx context.Context) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getAppSettingsList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		items = append(items, key)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeAppSetting = `-- name: RemoveAppSetting :exec
DELETE FROM app_settings WHERE key = $1
`

func (q *Queries) RemoveAppSetting(ctx context.Context, key string) error {
	_, err := q.db.ExecContext(ctx, removeAppSetting, key)
	return err
}

const updateAppSetting = `-- name: UpdateAppSetting :exec
UPDATE app_settings
SET
    value = $1,
    updated_at = $2
WHERE key = $3
`

type UpdateAppSettingParams struct {
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
	Key       string    `json:"key"`
}

func (q *Queries) UpdateAppSetting(ctx context.Context, arg UpdateAppSettingParams) error {
	_, err := q.db.ExecContext(ctx, updateAppSetting, arg.Value, arg.UpdatedAt, arg.Key)
	return err
}
