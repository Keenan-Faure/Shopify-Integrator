-- name: CreateWarehouse :exec
INSERT INTO warehouses(
    id,
    name,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4
);

-- name: GetWarehouses :many
SELECT
    id,
    name,
    updated_at
FROM warehouses;

-- name: RemoveWarehouse :exec
DELETE FROM warehouses
WHERE id = $1;

-- name: GetWarehouseByID :one
SELECT * FROM warehouses
WHERE id = $1;

-- name: GetWarehouseByName :one
SELECT * FROM warehouses
WHERE name = $1;