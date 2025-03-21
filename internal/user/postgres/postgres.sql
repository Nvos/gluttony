-- name: GetUser :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES ($1, $2)
ON CONFLICT (username) DO UPDATE SET username=EXCLUDED.username
RETURNING id;