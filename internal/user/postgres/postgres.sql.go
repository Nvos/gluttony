// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: postgres.sql

package postgres

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES ($1, $2)
ON CONFLICT (username) DO UPDATE SET username=EXCLUDED.username
RETURNING id
`

type CreateUserParams struct {
	Username string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (int32, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Username, arg.Password)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const getUser = `-- name: GetUser :one
SELECT id, username, password, role
FROM users
WHERE username = $1
LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Role,
	)
	return i, err
}
