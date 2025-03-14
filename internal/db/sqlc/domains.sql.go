// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: domains.sql

package db

import (
	"context"
	"database/sql"

	"github.com/sqlc-dev/pqtype"
)

const createDomain = `-- name: CreateDomain :one
INSERT INTO domains (name, provider, status, user_id, created_at, expires_at)
VALUES ($1, $2, $3, $4, NOW(), $5)
RETURNING id, name, provider, status, created_at, expires_at, is_deleted, settings, user_id
`

type CreateDomainParams struct {
	Name      string       `json:"name"`
	Provider  string       `json:"provider"`
	Status    int32        `json:"status"`
	UserID    int32        `json:"user_id"`
	ExpiresAt sql.NullTime `json:"expires_at"`
}

func (q *Queries) CreateDomain(ctx context.Context, arg CreateDomainParams) (Domain, error) {
	row := q.db.QueryRowContext(ctx, createDomain,
		arg.Name,
		arg.Provider,
		arg.Status,
		arg.UserID,
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
		&i.Settings,
		&i.UserID,
	)
	return i, err
}

const deleteDomain = `-- name: DeleteDomain :exec
UPDATE domains 
SET is_deleted = true 
WHERE id = $1
`

func (q *Queries) DeleteDomain(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteDomain, id)
	return err
}

const getAllDomains = `-- name: GetAllDomains :many
SELECT d.id, d.name, d.provider, d.status, d.created_at, d.expires_at, d.is_deleted, d.settings, d.user_id, COUNT(m.id) as mailboxes_count 
FROM domains d
LEFT JOIN mailboxes m ON d.id = m.domain_id AND m.is_deleted = false
WHERE d.is_deleted = false AND d.user_id = $3
GROUP BY d.id
ORDER BY d.created_at DESC
LIMIT $1 OFFSET $2
`

type GetAllDomainsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
	UserID int32 `json:"user_id"`
}

type GetAllDomainsRow struct {
	ID             int32                 `json:"id"`
	Name           string                `json:"name"`
	Provider       string                `json:"provider"`
	Status         int32                 `json:"status"`
	CreatedAt      sql.NullTime          `json:"created_at"`
	ExpiresAt      sql.NullTime          `json:"expires_at"`
	IsDeleted      sql.NullBool          `json:"is_deleted"`
	Settings       pqtype.NullRawMessage `json:"settings"`
	UserID         int32                 `json:"user_id"`
	MailboxesCount int64                 `json:"mailboxes_count"`
}

func (q *Queries) GetAllDomains(ctx context.Context, arg GetAllDomainsParams) ([]GetAllDomainsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllDomains, arg.Limit, arg.Offset, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllDomainsRow{}
	for rows.Next() {
		var i GetAllDomainsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Provider,
			&i.Status,
			&i.CreatedAt,
			&i.ExpiresAt,
			&i.IsDeleted,
			&i.Settings,
			&i.UserID,
			&i.MailboxesCount,
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
SELECT id, name, provider, status, created_at, expires_at, is_deleted, settings FROM domains
WHERE id = $1 AND is_deleted = false LIMIT 1
`

type GetDomainByIDRow struct {
	ID        int32                 `json:"id"`
	Name      string                `json:"name"`
	Provider  string                `json:"provider"`
	Status    int32                 `json:"status"`
	CreatedAt sql.NullTime          `json:"created_at"`
	ExpiresAt sql.NullTime          `json:"expires_at"`
	IsDeleted sql.NullBool          `json:"is_deleted"`
	Settings  pqtype.NullRawMessage `json:"settings"`
}

func (q *Queries) GetDomainByID(ctx context.Context, id int32) (GetDomainByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getDomainByID, id)
	var i GetDomainByIDRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Provider,
		&i.Status,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.IsDeleted,
		&i.Settings,
	)
	return i, err
}

const getDomainByName = `-- name: GetDomainByName :one
SELECT id, name, provider, status, created_at, expires_at, is_deleted
FROM domains
WHERE name = $1 AND is_deleted = false
`

type GetDomainByNameRow struct {
	ID        int32        `json:"id"`
	Name      string       `json:"name"`
	Provider  string       `json:"provider"`
	Status    int32        `json:"status"`
	CreatedAt sql.NullTime `json:"created_at"`
	ExpiresAt sql.NullTime `json:"expires_at"`
	IsDeleted sql.NullBool `json:"is_deleted"`
}

func (q *Queries) GetDomainByName(ctx context.Context, name string) (GetDomainByNameRow, error) {
	row := q.db.QueryRowContext(ctx, getDomainByName, name)
	var i GetDomainByNameRow
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

const getDomainsByUserID = `-- name: GetDomainsByUserID :many
SELECT id, name, provider, status, created_at, expires_at, is_deleted, settings, user_id FROM domains
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type GetDomainsByUserIDParams struct {
	UserID int32 `json:"user_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetDomainsByUserID(ctx context.Context, arg GetDomainsByUserIDParams) ([]Domain, error) {
	rows, err := q.db.QueryContext(ctx, getDomainsByUserID, arg.UserID, arg.Limit, arg.Offset)
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
			&i.Settings,
			&i.UserID,
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

const getDomainsCount = `-- name: GetDomainsCount :one
SELECT COUNT(*) FROM domains
WHERE is_deleted = false AND user_id = $1
`

func (q *Queries) GetDomainsCount(ctx context.Context, userID int32) (int64, error) {
	row := q.db.QueryRowContext(ctx, getDomainsCount, userID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getUserByDomainID = `-- name: GetUserByDomainID :one
SELECT u.id, u.username, u.email, u.hashed_password, u.full_name, u.role, u.is_active, u.created_at, u.updated_at FROM users u
JOIN domains d ON d.user_id = u.id
WHERE d.id = $1 AND u.is_active = true
LIMIT 1
`

func (q *Queries) GetUserByDomainID(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByDomainID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.HashedPassword,
		&i.FullName,
		&i.Role,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const setDomainStatus = `-- name: SetDomainStatus :exec
UPDATE domains 
SET status = $2
WHERE id = $1
`

type SetDomainStatusParams struct {
	ID     int32 `json:"id"`
	Status int32 `json:"status"`
}

func (q *Queries) SetDomainStatus(ctx context.Context, arg SetDomainStatusParams) error {
	_, err := q.db.ExecContext(ctx, setDomainStatus, arg.ID, arg.Status)
	return err
}

const updateDomainAndMailboxesStatus = `-- name: UpdateDomainAndMailboxesStatus :exec
WITH updated_domain AS (
    UPDATE domains
    SET status = $2
    WHERE id = $1
    AND status != $2
    RETURNING id
)
UPDATE mailboxes
SET status = $2
WHERE domain_id = $1
AND status != $2
`

type UpdateDomainAndMailboxesStatusParams struct {
	DomainID int32 `json:"domain_id"`
	Status   int32 `json:"status"`
}

func (q *Queries) UpdateDomainAndMailboxesStatus(ctx context.Context, arg UpdateDomainAndMailboxesStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateDomainAndMailboxesStatus, arg.DomainID, arg.Status)
	return err
}
