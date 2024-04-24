-- name: AllIngredients :many
SELECT id, (name ->> sqlc.arg('locale')::text)::text as name
FROM ingredients
WHERE CASE
          WHEN sqlc.arg('search') != '' THEN to_tsvector(name ->> sqlc.arg('locale')::text) @@
                                             websearch_to_tsquery(sqlc.arg('search'))
          ELSE TRUE END
  AND (name ->> sqlc.arg('locale')::text)::text IS NOT NULL
OFFSET sqlc.arg('offset') ROWS FETCH FIRST sqlc.arg('limit') ROW ONLY;

-- name: CreateIngredient :exec
INSERT INTO ingredients (name)
values ($1);