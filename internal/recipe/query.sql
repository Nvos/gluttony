-- name: SingleRecipe :many
SELECT sqlc.embed(recipes), sqlc.embed(recipe_steps) FROM recipes
JOIN recipe_steps on recipes.id = recipe_steps.recipe_id
WHERE recipes.id = $1
ORDER BY recipe_steps.order;

-- name: AllRecipes :many
SELECT * FROM recipes, ts_rank(to_tsvector(name), websearch_to_tsquery($3)) as rank
WHERE CASE WHEN $3 != '' THEN to_tsvector(name) @@ websearch_to_tsquery($3) ELSE TRUE END
ORDER BY
    recipes.id DESC,
    CASE WHEN $3 != '' THEN rank END DESC
OFFSET $1 ROWS
FETCH FIRST $2 ROW ONLY;

-- name: CreateRecipe :one
INSERT INTO recipes (name, description)
VALUES ($1, $2)
RETURNING id;

-- name: CreateRecipeSteps :copyfrom
INSERT INTO recipe_steps (recipe_id, description, "order") VALUES ($1, $2, $3);