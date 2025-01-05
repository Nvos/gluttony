-- name: GetRecipe :one
SELECT *
FROM recipes
WHERE id = ?
LIMIT 1;

-- name: CreateRecipe :one
INSERT INTO recipes (name, description, instructions_markdown, thumbnail_url,
                     cook_time_seconds, preparation_time_seconds, source, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id;

-- name: CreateNutrition :exec
INSERT INTO recipe_nutrition (recipe_id, calories, fat, carbs, protein)
VALUES (?, ?, ?, ?, ?);

-- name: CreateTag :one
INSERT INTO tags (name)
VALUES (?)
ON CONFLICT DO NOTHING
RETURNING id;

-- name: CreateIngredient :one
INSERT INTO ingredients (name)
VALUES (?)
ON CONFLICT DO NOTHING
RETURNING id;

-- name: CreateRecipeIngredient :exec
INSERT INTO recipe_ingredients (recipe_order, recipe_id, ingredient_id, unit, quantity)
VALUES (?, ?, ?, ?, ?);

-- name: CreateRecipeTag :exec
INSERT INTO recipe_tags (recipe_order, recipe_id, tag_id)
VALUES (?, ?, ?);