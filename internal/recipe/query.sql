-- name: GetRecipe :one
SELECT *
FROM recipes
WHERE id = ?
LIMIT 1;

-- name: CreateRecipe :one

INSERT INTO recipes (name, content_markdown, created_at)
VALUES (?, ?, ?)
RETURNING *;