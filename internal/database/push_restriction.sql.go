// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: push_restriction.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createPushReestriction = `-- name: CreatePushReestriction :exec
INSERT INTO push_restriction(
    id,
    field,
    flag,
    updated_at,
    created_at
) VALUES (
    $1, $2, $3, $4, $5
)
`

type CreatePushReestrictionParams struct {
	ID        uuid.UUID `json:"id"`
	Field     string    `json:"field"`
	Flag      string    `json:"flag"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (q *Queries) CreatePushReestriction(ctx context.Context, arg CreatePushReestrictionParams) error {
	_, err := q.db.ExecContext(ctx, createPushReestriction,
		arg.ID,
		arg.Field,
		arg.Flag,
		arg.UpdatedAt,
		arg.CreatedAt,
	)
	return err
}

const updatePushRestriction = `-- name: UpdatePushRestriction :exec
UPDATE push_restriction
SET
    flag = $1,
    updated_at = $2
WHERE field = $3
`

type UpdatePushRestrictionParams struct {
	Flag      string    `json:"flag"`
	UpdatedAt time.Time `json:"updated_at"`
	Field     string    `json:"field"`
}

func (q *Queries) UpdatePushRestriction(ctx context.Context, arg UpdatePushRestrictionParams) error {
	_, err := q.db.ExecContext(ctx, updatePushRestriction, arg.Flag, arg.UpdatedAt, arg.Field)
	return err
}
