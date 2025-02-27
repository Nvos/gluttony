-- name: GetFullRecipe :one
SELECT *
FROM recipes
         JOIN main.recipe_nutrition rn on recipes.id = rn.recipe_id
WHERE recipes.id = ?
LIMIT 1;

-- name: AllRecipeSummary :many
SELECT id, name, description, thumbnail_url
FROM recipes
WHERE (CAST(sqlc.arg('is_search') as INTEGER) = FALSE OR id in (sqlc.slice('ids')))
ORDER BY id DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: AllRecipeTags :many
SELECT *
FROM tags
         JOIN recipe_tags rt on tags.id = rt.tag_id
WHERE recipe_id in (sqlc.slice('ids'))
ORDER BY recipe_id, recipe_order;

-- name: AllRecipeIngredients :many
SELECT *
FROM ingredients
         JOIN recipe_ingredients ri on ingredients.id = ri.ingredient_id
WHERE recipe_id in (sqlc.slice('ids'))
ORDER BY recipe_id, recipe_order;


-- name: CreateRecipe :one
INSERT INTO recipes (name, description, instructions_markdown, thumbnail_url,
                     cook_time_seconds, preparation_time_seconds, source, owner_id, servings)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id;

-- name: UpdateRecipe :exec
UPDATE recipes
SET name                     = ?,
    description              = ?,
    instructions_markdown    = ?,
    thumbnail_url            = ?,
    cook_time_seconds        = ?,
    preparation_time_seconds = ?,
    source                   = ?,
    updated_at               = ?,
    servings                 = ?
WHERE id = ?;

-- name: CreateNutrition :exec
INSERT INTO recipe_nutrition (recipe_id, calories, fat, carbs, protein)
VALUES (?, ?, ?, ?, ?);

-- name: UpdateNutrition :exec
UPDATE recipe_nutrition
SET calories = ?,
    fat      = ?,
    carbs    = ?,
    protein  = ?
WHERE recipe_id = ?;

-- name: AllTagsByNames :many
SELECT id, name
FROM tags
WHERE name in (sqlc.slice('names'));

-- name: CreateTag :one
INSERT INTO tags (name)
VALUES (?)
RETURNING id;

-- name: AllIngredientsByNames :many
SELECT id, name
FROM ingredients
WHERE name in (sqlc.slice('names'));

-- name: CreateIngredient :one
INSERT INTO ingredients (name)
VALUES (?)
RETURNING id;

-- name: CreateRecipeIngredient :exec
INSERT INTO recipe_ingredients (recipe_order, recipe_id, ingredient_id, unit, quantity, note)
VALUES (?, ?, ?, ?, ?, ?);

-- name: CreateRecipeTag :exec
INSERT INTO recipe_tags (recipe_order, recipe_id, tag_id)
VALUES (?, ?, ?);

-- name: DeleteRecipeTags :exec
DELETE
FROM recipe_tags
WHERE recipe_id = sqlc.arg('recipe_id');

-- name: DeleteRecipeIngredients :exec
DELETE
FROM recipe_ingredients
WHERE recipe_id = sqlc.arg('recipe_id');
