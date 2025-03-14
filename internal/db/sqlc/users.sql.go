// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    hashed_password,
    full_name,
    role
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id, username, email, hashed_password, full_name, role, is_active, created_at, updated_at
`

type CreateUserParams struct {
	Username       string         `json:"username"`
	Email          string         `json:"email"`
	HashedPassword string         `json:"hashed_password"`
	FullName       sql.NullString `json:"full_name"`
	Role           string         `json:"role"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.Email,
		arg.HashedPassword,
		arg.FullName,
		arg.Role,
	)
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

const deactivateUser = `-- name: DeactivateUser :exec
UPDATE users
SET
    is_active = false,
    updated_at = NOW()
WHERE id = $1
`

func (q *Queries) DeactivateUser(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deactivateUser, id)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, username, email, hashed_password, full_name, role, is_active, created_at, updated_at FROM users
WHERE email = $1 AND is_active = true LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
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

const getUserByID = `-- name: GetUserByID :one
SELECT id, username, email, hashed_password, full_name, role, is_active, created_at, updated_at FROM users
WHERE id = $1 AND is_active = true LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
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

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, email, hashed_password, full_name, role, is_active, created_at, updated_at FROM users
WHERE username = $1 AND is_active = true LIMIT 1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
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

const getUserStats = `-- name: GetUserStats :one
SELECT
    COUNT(DISTINCT d.id) as domains_count,
    COUNT(DISTINCT m.id) as mailboxes_count,
    COUNT(DISTINCT CASE WHEN m.status = 1 THEN m.id END) as active_mailboxes_count,
    COUNT(DISTINCT CASE WHEN m.status = 2 THEN m.id END) as inactive_mailboxes_count
FROM users u
LEFT JOIN domains d ON d.user_id = u.id
LEFT JOIN mailboxes m ON m.user_id = u.id
WHERE u.id = $1
`

type GetUserStatsRow struct {
	DomainsCount           int64 `json:"domains_count"`
	MailboxesCount         int64 `json:"mailboxes_count"`
	ActiveMailboxesCount   int64 `json:"active_mailboxes_count"`
	InactiveMailboxesCount int64 `json:"inactive_mailboxes_count"`
}

func (q *Queries) GetUserStats(ctx context.Context, id int32) (GetUserStatsRow, error) {
	row := q.db.QueryRowContext(ctx, getUserStats, id)
	var i GetUserStatsRow
	err := row.Scan(
		&i.DomainsCount,
		&i.MailboxesCount,
		&i.ActiveMailboxesCount,
		&i.InactiveMailboxesCount,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, username, email, hashed_password, full_name, role, is_active, created_at, updated_at FROM users
WHERE is_active = true
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.HashedPassword,
			&i.FullName,
			&i.Role,
			&i.IsActive,
			&i.CreatedAt,
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

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET
    username = COALESCE($1, username),
    email = COALESCE($2, email),
    full_name = COALESCE($3, full_name),
    role = COALESCE($4, role),
    updated_at = NOW()
WHERE id = $5
RETURNING id, username, email, hashed_password, full_name, role, is_active, created_at, updated_at
`

type UpdateUserParams struct {
	Username sql.NullString `json:"username"`
	Email    sql.NullString `json:"email"`
	FullName sql.NullString `json:"full_name"`
	Role     sql.NullString `json:"role"`
	ID       int32          `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.Username,
		arg.Email,
		arg.FullName,
		arg.Role,
		arg.ID,
	)
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

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE users
SET
    hashed_password = $2,
    updated_at = NOW()
WHERE id = $1
`

type UpdateUserPasswordParams struct {
	ID             int32  `json:"id"`
	HashedPassword string `json:"hashed_password"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.ExecContext(ctx, updateUserPassword, arg.ID, arg.HashedPassword)
	return err
}

const verifyUserCredentials = `-- name: VerifyUserCredentials :one
SELECT id, username, email, hashed_password, full_name, role, is_active, created_at, updated_at FROM users
WHERE email = $1 AND is_active = true
LIMIT 1
`

func (q *Queries) VerifyUserCredentials(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, verifyUserCredentials, email)
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
