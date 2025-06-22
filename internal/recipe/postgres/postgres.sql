-- name: GetFullRecipe :one
SELECT *
FROM recipes
         JOIN recipe_nutrition rn on recipes.id = rn.recipe_id
         JOIN images im on recipes.thumbnail_id = im.id
WHERE recipes.id = $1
LIMIT 1;

-- name: AllRecipeTags :many
SELECT *
FROM tags
         JOIN recipe_tags rt on tags.id = rt.tag_id
WHERE recipe_id = ANY (sqlc.slice('ids')::int[])
ORDER BY recipe_id, recipe_order;

-- name: AllRecipeIngredients :many
SELECT *
FROM ingredients
         JOIN recipe_ingredients ri on ingredients.id = ri.ingredient_id
WHERE recipe_id = ANY (sqlc.slice('ids')::int[])
ORDER BY recipe_id, recipe_order;


-- name: CreateRecipe :one
INSERT INTO recipes (name, description, instructions_markdown, thumbnail_id,
                     cook_time_seconds, preparation_time_seconds, source, owner_id, servings)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id;

-- name: UpdateRecipe :exec
UPDATE recipes
SET name                     = $1,
    description              = $2,
    instructions_markdown    = $3,
    thumbnail_id             = $4,
    cook_time_seconds        = $5,
    preparation_time_seconds = $6,
    source                   = $7,
    updated_at               = $8,
    servings                 = $9
WHERE id = $10;

-- name: CreateNutrition :exec
INSERT INTO recipe_nutrition (recipe_id, calories, fat, carbs, protein)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateNutrition :exec
UPDATE recipe_nutrition
SET calories = $1,
    fat      = $2,
    carbs    = $3,
    protein  = $4
WHERE recipe_id = $5;

-- name: AllTagsByNames :many
SELECT id, name
FROM tags
WHERE name = ANY (sqlc.slice('names')::text[]);

-- name: CreateTag :one
INSERT INTO tags (name)
VALUES ($1)
RETURNING id;

-- name: AllIngredientsByNames :many
SELECT id, name
FROM ingredients
WHERE name = ANY (sqlc.slice('names')::text[]);

-- name: CreateIngredient :one
INSERT INTO ingredients (name)
VALUES ($1)
RETURNING id;

-- name: CreateRecipeIngredient :exec
INSERT INTO recipe_ingredients (recipe_order, recipe_id, ingredient_id, unit, quantity, note)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: CreateRecipeTag :exec
INSERT INTO recipe_tags (recipe_order, recipe_id, tag_id)
VALUES ($1, $2, $3);

-- name: DeleteRecipeTags :exec
DELETE
FROM recipe_tags
WHERE recipe_id = sqlc.arg('recipe_id');

-- name: DeleteRecipeIngredients :exec
DELETE
FROM recipe_ingredients
WHERE recipe_id = sqlc.arg('recipe_id');

-- name: AllRecipeSummaries :many
SELECT recipes.id, recipes.name, recipes.description, images.url
FROM recipes
         LEFT JOIN images on recipes.thumbnail_id = images.id
WHERE (sqlc.slice(ids)::int[] IS NULL OR recipes.id = ANY (sqlc.slice(ids)::int[]))
ORDER BY recipes.id DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountRecipeSummaries :one
SELECT count(*)
FROM recipes;

-- name: CreateRecipeImage :one
INSERT INTO images (url)
VALUES ($1)
RETURNING id;

