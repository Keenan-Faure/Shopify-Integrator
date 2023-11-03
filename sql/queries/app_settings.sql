-- name: AddAppSetting :exec
INSERT INTO app_settings(
    id,
    key,
    description,
    value,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: UpdateAppSetting :exec
UPDATE app_settings
SET
    value = $1,
    updated_at = $2
WHERE key = $3;

-- name: GetAppSettingByKey :one
SELECT
    id,
    key,
    description,
    value,
    updated_at
FROM app_settings
WHERE key = $1;

-- name: GetAppSettings :many
SELECT
    id,
    key,
    description,
    value,
    updated_at
FROM app_settings;

-- name: RemoveAppSetting :exec
DELETE FROM app_settings WHERE key = $1;