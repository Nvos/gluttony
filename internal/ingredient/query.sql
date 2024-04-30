-- name: AllIngredients :many
SELECT ingredients.id,
       (ingredients.name ->> sqlc.arg('locale')::text)::text as name
FROM ingredients,
     to_tsvector(sqlc.arg('search_locale')::regconfig, ingredients.name ->> sqlc.arg('locale')::text) document,
     websearch_to_tsquery(sqlc.arg('search_locale')::regconfig, sqlc.arg('search')) query,
     similarity(sqlc.arg('search'), ingredients.name ->> sqlc.arg('locale')::text) similarity,
     NULLIF(ts_rank(document, query), 0)                  rank_name
WHERE ingredients.name ->> sqlc.arg('locale') IS NOT NULL
  AND CASE
          WHEN sqlc.arg('search') != '' THEN (query @@ document OR similarity > 0.5)
          ELSE TRUE END
ORDER BY rank_name, similarity DESC NULLS LAST
OFFSET sqlc.arg('offset') ROWS FETCH FIRST sqlc.arg('limit') ROW ONLY;

-- name: CreateIngredient :exec
INSERT INTO ingredients (name)
values ($1);