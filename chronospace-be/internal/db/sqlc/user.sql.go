// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    username,
    full_name,
    email,
    password
) VALUES (
    $1, $2, $3, $4
) RETURNING id, username, full_name, email, password, created_at
`

type CreateUserParams struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Username,
		arg.FullName,
		arg.Email,
		arg.Password,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.FullName,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, username, full_name, email, password, created_at FROM users
WHERE id = $1
`

func (q *Queries) GetUser(ctx context.Context, id pgtype.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.FullName,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, username, full_name, email, password, created_at FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.FullName,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, full_name, email, password, created_at FROM users
WHERE username = $1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.FullName,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, username, full_name, email, password, created_at FROM users
ORDER BY created_at
LIMIT $1
OFFSET $2
`

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers, arg.Limit, arg.Offset)
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
			&i.FullName,
			&i.Email,
			&i.Password,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET 
    username = COALESCE($2, username),
    full_name = COALESCE($3, full_name),
    email = COALESCE($4, email),
    password = COALESCE($5, password)
WHERE id = $1
RETURNING id, username, full_name, email, password, created_at
`

type UpdateUserParams struct {
	ID       pgtype.UUID `json:"id"`
	Username string      `json:"username"`
	FullName string      `json:"full_name"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.ID,
		arg.Username,
		arg.FullName,
		arg.Email,
		arg.Password,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.FullName,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :one
UPDATE users
SET password = $2
WHERE id = $1
RETURNING id, username, full_name, email, password, created_at
`

type UpdateUserPasswordParams struct {
	ID       pgtype.UUID `json:"id"`
	Password string      `json:"password"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUserPassword, arg.ID, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.FullName,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
	)
	return i, err
}
