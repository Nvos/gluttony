-- name: SingleRecipe :one
SELECT id,
       (name ->> sqlc.arg('locale')::text)::text as name,
       (description ->> sqlc.arg('locale')::text)::text as description,
       (content ->> sqlc.arg('locale')::text)::text    as content
FROM recipes
WHERE recipes.id = sqlc.arg('recipe_id');

-- name: AllRecipeIngredients :many
SELECT id, (name ->> sqlc.arg('locale')::text)::text as name, amount, count, note, unit FROM ingredients
JOIN recipes_ingredients ri on ingredients.id = ri.ingredient_id
where ri.recipe_id = sqlc.arg('recipe_id');

-- name: AllRecipes :many
SELECT id, (name ->> sqlc.arg('locale'))::text as name, (description ->> sqlc.arg('locale'))::text as description
FROM recipes as rank
WHERE CASE
          WHEN sqlc.arg('search') != '' THEN to_tsvector(name->>sqlc.arg('locale')::text) @@ websearch_to_tsquery(sqlc.arg('search'))
          ELSE TRUE END
OFFSET sqlc.arg('offset') ROWS FETCH FIRST sqlc.arg('limit') ROW ONLY;

-- name: CreateRecipe :one
INSERT INTO recipes (name, description, content)
VALUES ($1, $2, $3)
RETURNING id;

-- name: CreateRecipeIngredientEdges :copyfrom
INSERT INTO recipes_ingredients (recipe_id, ingredient_id, amount, count, note)
VALUES ($1, $2, $3, $4, $5);
