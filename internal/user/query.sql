-- name: GetUser :one
SELECT *
FROM users
WHERE username = ?
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES (?, ?)
ON CONFLICT (username) DO UPDATE SET username=EXCLUDED.username
RETURNING id;