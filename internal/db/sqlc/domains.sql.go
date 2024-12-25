// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: domains.sql

package db

import (
	"context"
	"database/sql"
)

const createDomain = `-- name: CreateDomain :one
INSERT INTO domains (name, provider, status, created_at, expires_at)
VALUES ($1, $2, $3, NOW(), $4)
RETURNING id, name, provider, status, created_at, expires_at, is_deleted
`

type CreateDomainParams struct {
	Name      string       `json:"name"`
	Provider  string       `json:"provider"`
	Status    int32        `json:"status"`
	ExpiresAt sql.NullTime `json:"expires_at"`
}

func (q *Queries) CreateDomain(ctx context.Context, arg CreateDomainParams) (Domain, error) {
	row := q.db.QueryRowContext(ctx, createDomain,
		arg.Name,
		arg.Provider,
		arg.Status,
		arg.ExpiresAt,
	)
	var i Domain
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Provider,
		&i.Status,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.IsDeleted,
	)
	return i, err
}

const deleteDomain = `-- name: DeleteDomain :exec
DELETE FROM domains
WHERE id = $1
`

func (q *Queries) DeleteDomain(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteDomain, id)
	return err
}

const getAllDomains = `-- name: GetAllDomains :many
SELECT id, name, provider, status, created_at, expires_at, is_deleted FROM domains
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

type GetAllDomainsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetAllDomains(ctx context.Context, arg GetAllDomainsParams) ([]Domain, error) {
	rows, err := q.db.QueryContext(ctx, getAllDomains, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Domain{}
	for rows.Next() {
		var i Domain
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Provider,
			&i.Status,
			&i.CreatedAt,
			&i.ExpiresAt,
			&i.IsDeleted,
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

const getDomainByID = `-- name: GetDomainByID :one
SELECT id, name, provider, status, created_at, expires_at, is_deleted FROM domains
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetDomainByID(ctx context.Context, id int32) (Domain, error) {
	row := q.db.QueryRowContext(ctx, getDomainByID, id)
	var i Domain
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Provider,
		&i.Status,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.IsDeleted,
	)
	return i, err
}

const getDomainByName = `-- name: GetDomainByName :one
SELECT id, name, provider, status, created_at, expires_at, is_deleted
FROM domains
WHERE name = $1
`

func (q *Queries) GetDomainByName(ctx context.Context, name string) (Domain, error) {
	row := q.db.QueryRowContext(ctx, getDomainByName, name)
	var i Domain
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Provider,
		&i.Status,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.IsDeleted,
	)
	return i, err
}

const getDomainsCount = `-- name: GetDomainsCount :one
SELECT COUNT(*) FROM domains
`

func (q *Queries) GetDomainsCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, getDomainsCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const setDomainStatus = `-- name: SetDomainStatus :exec
UPDATE domains
SET status = $1
WHERE id = $2
`

type SetDomainStatusParams struct {
	Status int32 `json:"status"`
	ID     int32 `json:"id"`
}

func (q *Queries) SetDomainStatus(ctx context.Context, arg SetDomainStatusParams) error {
	_, err := q.db.ExecContext(ctx, setDomainStatus, arg.Status, arg.ID)
	return err
}
