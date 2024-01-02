// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: fetch_worker.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFetchWorker = `-- name: CreateFetchWorker :exec
INSERT INTO fetch_worker(
    id,
    status,
    local_count,
    shopify_product_count,
    fetch_url,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7 
)
`

type CreateFetchWorkerParams struct {
	ID                  uuid.UUID `json:"id"`
	Status              string    `json:"status"`
	LocalCount          int32     `json:"local_count"`
	ShopifyProductCount int32     `json:"shopify_product_count"`
	FetchUrl            string    `json:"fetch_url"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (q *Queries) CreateFetchWorker(ctx context.Context, arg CreateFetchWorkerParams) error {
	_, err := q.db.ExecContext(ctx, createFetchWorker,
		arg.ID,
		arg.Status,
		arg.LocalCount,
		arg.ShopifyProductCount,
		arg.FetchUrl,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const getFetchWorker = `-- name: GetFetchWorker :one
SELECT id, status, local_count, shopify_product_count, fetch_url, created_at, updated_at FROM fetch_worker
LIMIT 1
`

func (q *Queries) GetFetchWorker(ctx context.Context) (FetchWorker, error) {
	row := q.db.QueryRowContext(ctx, getFetchWorker)
	var i FetchWorker
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.LocalCount,
		&i.ShopifyProductCount,
		&i.FetchUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const resetFetchWorker = `-- name: ResetFetchWorker :exec
UPDATE fetch_worker SET status = $1
`

func (q *Queries) ResetFetchWorker(ctx context.Context, status string) error {
	_, err := q.db.ExecContext(ctx, resetFetchWorker, status)
	return err
}

const updateFetchWorker = `-- name: UpdateFetchWorker :exec
UPDATE fetch_worker
SET
    status = $1,
    local_count = $2,
    shopify_product_count = $3,
    fetch_url = $4,
    updated_at = $5
WHERE id = $6
`

type UpdateFetchWorkerParams struct {
	Status              string    `json:"status"`
	LocalCount          int32     `json:"local_count"`
	ShopifyProductCount int32     `json:"shopify_product_count"`
	FetchUrl            string    `json:"fetch_url"`
	UpdatedAt           time.Time `json:"updated_at"`
	ID                  uuid.UUID `json:"id"`
}

func (q *Queries) UpdateFetchWorker(ctx context.Context, arg UpdateFetchWorkerParams) error {
	_, err := q.db.ExecContext(ctx, updateFetchWorker,
		arg.Status,
		arg.LocalCount,
		arg.ShopifyProductCount,
		arg.FetchUrl,
		arg.UpdatedAt,
		arg.ID,
	)
	return err
}