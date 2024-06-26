// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package postgresql

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users(name, password)
VALUES ($1, $2)
RETURNING id
`

type CreateUserParams struct {
	Name     string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (int32, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Name, arg.Password)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const singleUser = `-- name: SingleUser :one
SELECT id, name, password FROM users
WHERE id = $1
LIMIT 1
`

func (q *Queries) SingleUser(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRow(ctx, singleUser, id)
	var i User
	err := row.Scan(&i.ID, &i.Name, &i.Password)
	return i, err
}

const singleUserByName = `-- name: SingleUserByName :one
SELECT id, name, password FROM users
WHERE users.name = $1
LIMIT 1
`

func (q *Queries) SingleUserByName(ctx context.Context, name string) (User, error) {
	row := q.db.QueryRow(ctx, singleUserByName, name)
	var i User
	err := row.Scan(&i.ID, &i.Name, &i.Password)
	return i, err
}
