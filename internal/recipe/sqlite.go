package recipe

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"gluttony/internal/database"
	"gluttony/internal/ingredient"
	"gluttony/internal/recipe/sqlite"
	"strings"
	"text/template"
	"time"
)

var allRecipesQuery = template.Must(
	template.New("").
		Funcs(template.FuncMap{
			"array": func(count int) string {
				return strings.Repeat(",?", count)[1:]
			},
		}).
		Parse(`
		SELECT id, name, description, thumbnail_url
		FROM recipes
		{{if ne (len .RecipeIDs) 0}}
		WHERE id in ({{array (len .RecipeIDs)}})
		{{end}}
		ORDER BY id DESC
		LIMIT ? OFFSET ?;
	`),
)

type Store struct {
	db      database.DBTX
	queries *sqlite.Queries
}

func NewStore(db database.DBTX) *Store {
	return &Store{
		db:      db,
		queries: sqlite.New(db),
	}
}

func (s *Store) WithTx(tx *sql.Tx) *Store {
	return &Store{
		db:      tx,
		queries: sqlite.New(tx),
	}
}

func (s *Store) AllRecipeSummaries(
	ctx context.Context,
	input SearchInput,
) ([]Summary, error) {
	var buffer bytes.Buffer
	if err := allRecipesQuery.Execute(&buffer, input); err != nil {
		return nil, err
	}

	query := buffer.String()
	params := make([]any, 0, 2+len(input.RecipeIDs))
	for i := range input.RecipeIDs {
		params = append(params, input.RecipeIDs[i])
	}
	params = append(params, input.Limit, input.Page*input.Limit)

	rows, err := s.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	out := make([]Summary, 0, 20)
	for rows.Next() {
		value := Summary{}
		err = rows.Scan(
			&value.ID,
			&value.Name,
			&value.Description,
			&value.ThumbnailImageURL,
		)
		if err != nil {
			return nil, err
		}

		for i := range input.RecipeIDs {
			if value.ID != input.RecipeIDs[i] {
				continue
			}

			out[i] = value
			break
		}

		out = append(out, value)
	}

	return out, nil
}

func (s *Store) GetRecipe(ctx context.Context, ID int64) (Recipe, error) {
	recipe, err := s.queries.GetFullRecipe(ctx, ID)
	if err != nil {
		return Recipe{}, fmt.Errorf("get recipe id=%d: %w", ID, err)
	}

	tags, err := s.AllTagsByRecipeIDs(ctx, ID)
	if err != nil {
		return Recipe{}, fmt.Errorf("get all recipe tags for recipe id=%d: %w", ID, err)
	}

	ingredients, err := s.AllIngredientsByRecipeIDs(ctx, ID)
	if err != nil {
		return Recipe{}, fmt.Errorf("get all ingredients for recipe id=%d: %w", ID, err)
	}

	return Recipe{
		ID:                   ID,
		Name:                 recipe.Name,
		Description:          recipe.Description,
		InstructionsMarkdown: recipe.InstructionsMarkdown,
		ThumbnailImageURL:    recipe.ThumbnailUrl,
		Tags:                 tags[ID],
		Source:               recipe.Source,
		Servings:             int8(recipe.Servings),
		PreparationTime:      time.Duration(recipe.PreparationTimeSeconds) * time.Second,
		CookTime:             time.Duration(recipe.CookTimeSeconds) * time.Second,
		Ingredients:          ingredients[ID],
		Nutrition: Nutrition{
			Calories: float32(recipe.Calories),
			Fat:      float32(recipe.Fat),
			Carbs:    float32(recipe.Carbs),
			Protein:  float32(recipe.Protein),
		},
	}, nil
}

func (s *Store) AllTagsByRecipeIDs(ctx context.Context, recipeIDs ...int64) (map[int64][]Tag, error) {
	tags, err := s.queries.AllRecipeTags(ctx, recipeIDs)
	if err != nil {
		return nil, fmt.Errorf("get all tags by recipe ids = %+v: %w", recipeIDs, err)
	}

	out := make(map[int64][]Tag, len(tags))
	for i := range tags {
		if out[tags[i].RecipeID] == nil {
			out[tags[i].RecipeID] = []Tag{}
		}

		out[tags[i].RecipeID] = append(out[tags[i].RecipeID], Tag{
			ID:    int(tags[i].ID),
			Order: int(tags[i].RecipeOrder),
			Name:  tags[i].Name,
		})
	}

	return out, nil
}

func (s *Store) AllIngredientsByRecipeIDs(
	ctx context.Context,
	recipeIDs ...int64,
) (map[int64][]Ingredient, error) {
	ingredients, err := s.queries.AllRecipeIngredients(ctx, recipeIDs)
	if err != nil {
		return nil, fmt.Errorf("get all ingredients by recipe ids = %+v: %w", recipeIDs, err)
	}

	out := make(map[int64][]Ingredient, len(ingredients))
	for i := range ingredients {
		if out[ingredients[i].RecipeID] == nil {
			out[ingredients[i].RecipeID] = []Ingredient{}
		}

		out[ingredients[i].RecipeID] = append(out[ingredients[i].RecipeID], Ingredient{
			Ingredient: ingredient.Ingredient{
				ID:   int(ingredients[i].ID),
				Name: ingredients[i].Name,
			},
			Order:    int8(ingredients[i].RecipeOrder),
			Quantity: float32(ingredients[i].Quantity),
			Unit:     ingredients[i].Unit,
		})
	}

	return out, nil
}

func (s *Store) CreateRecipeIngredients(
	ctx context.Context,
	recipeID int64,
	ingredients []Ingredient,
) error {

	names := make([]string, len(ingredients))
	for i := range ingredients {
		names[i] = ingredients[i].Name
	}

	existingIngredients, err := s.queries.AllIngredientsByNames(ctx, names)
	if err != nil {
		return fmt.Errorf("all ingredients by names: %w", err)
	}

	existingIngredientLookup := make(map[string]sqlite.Ingredient, len(existingIngredients))
	for i := range existingIngredients {
		existingIngredientLookup[existingIngredients[i].Name] = existingIngredients[i]
	}

	savedIngredients := make([]Ingredient, len(ingredients))
	for i := range ingredients {
		var ingredientID int64
		if value, ok := existingIngredientLookup[ingredients[i].Name]; ok {
			ingredientID = value.ID
		} else {
			id, err := s.queries.CreateIngredient(ctx, ingredients[i].Name)
			if err != nil {
				return fmt.Errorf("create ingredient: %w", err)
			}

			ingredientID = id
		}

		savedIngredients[i].Ingredient.ID = int(ingredientID)
		savedIngredients[i].Ingredient.Name = ingredients[i].Name
		savedIngredients[i].Order = int8(i)
		savedIngredients[i].Quantity = ingredients[i].Quantity
		savedIngredients[i].Unit = ingredients[i].Unit
	}

	for i := range savedIngredients {
		err = s.queries.CreateRecipeIngredient(ctx, sqlite.CreateRecipeIngredientParams{
			RecipeOrder:  int64(savedIngredients[i].Order),
			RecipeID:     recipeID,
			IngredientID: int64(savedIngredients[i].Ingredient.ID),
			Unit:         savedIngredients[i].Unit,
			Quantity:     int64(savedIngredients[i].Quantity),
		})
		if err != nil {
			return fmt.Errorf("create recipe ingredient: %w", err)
		}
	}

	return nil
}

func (s *Store) CreateRecipeNutrition(
	ctx context.Context,
	recipeID int64,
	nutrition Nutrition,
) error {
	err := s.queries.CreateNutrition(ctx, sqlite.CreateNutritionParams{
		RecipeID: recipeID,
		Calories: float64(nutrition.Calories),
		Fat:      float64(nutrition.Fat),
		Carbs:    float64(nutrition.Carbs),
		Protein:  float64(nutrition.Protein),
	})
	if err != nil {
		return fmt.Errorf("create nutrition: %w", err)
	}

	return nil
}

func (s *Store) CreateRecipe(ctx context.Context, input CreateRecipe) (int64, error) {
	createRecipeParams := sqlite.CreateRecipeParams{
		Name:                   input.Name,
		Description:            input.Description,
		InstructionsMarkdown:   input.InstructionsMarkdown,
		CookTimeSeconds:        int64(input.CookTime.Seconds()),
		PreparationTimeSeconds: int64(input.PreparationTime.Seconds()),
		Source:                 input.Source,
		ThumbnailUrl:           input.ThumbnailImageURL,
	}

	id, err := s.queries.CreateRecipe(ctx, createRecipeParams)
	if err != nil {
		return 0, fmt.Errorf("create recipe: %w", err)
	}

	return id, nil
}

func (s *Store) CreateRecipeTags(
	ctx context.Context,
	recipeID int64,
	tagNames []string,
) error {
	existingTags, err := s.queries.AllTagsByNames(ctx, tagNames)
	if err != nil {
		return fmt.Errorf("all savedTags by names: %w", err)
	}

	existingTagLookup := make(map[string]sqlite.Tag, len(existingTags))
	for i := range existingTags {
		existingTagLookup[existingTags[i].Name] = existingTags[i]
	}

	savedTags := make([]Tag, len(tagNames))
	for i := range tagNames {
		var tagID int64
		if value, ok := existingTagLookup[tagNames[i]]; ok {
			tagID = value.ID
		} else {
			id, err := s.queries.CreateTag(ctx, tagNames[i])
			if err != nil {
				return fmt.Errorf("create tag: %w", err)
			}

			tagID = id
		}

		savedTags[i].Name = tagNames[i]
		savedTags[i].ID = int(tagID)
		savedTags[i].Order = i
	}

	for i := range savedTags {
		err = s.queries.CreateRecipeTag(ctx, sqlite.CreateRecipeTagParams{
			RecipeOrder: int64(savedTags[i].Order),
			RecipeID:    recipeID,
			TagID:       int64(savedTags[i].ID),
		})
		if err != nil {
			return fmt.Errorf("create recipe tag: %w", err)
		}
	}

	return nil
}

func (s *Store) DeleteRecipeTags(ctx context.Context, recipeID int64) error {
	if err := s.queries.DeleteRecipeTags(ctx, recipeID); err != nil {
		return fmt.Errorf("remove recipe tags: %w", err)
	}

	return nil
}

func (s *Store) DeleteRecipeIngredients(ctx context.Context, recipeID int64) error {
	if err := s.queries.DeleteRecipeIngredients(ctx, recipeID); err != nil {
		return fmt.Errorf("remove recipe ingredients: %w", err)
	}

	return nil
}

func (s *Store) UpdateNutrition(ctx context.Context, recipeID int64, nutrition Nutrition) error {
	err := s.queries.UpdateNutrition(ctx, sqlite.UpdateNutritionParams{
		RecipeID: recipeID,
		Calories: float64(nutrition.Calories),
		Fat:      float64(nutrition.Fat),
		Carbs:    float64(nutrition.Carbs),
		Protein:  float64(nutrition.Protein),
	})
	if err != nil {
		return fmt.Errorf("update nutrition: %w", err)
	}

	return nil
}

func (s *Store) UpdateRecipe(ctx context.Context, input UpdateRecipe) error {
	params := sqlite.UpdateRecipeParams{
		ID:                     input.ID,
		Name:                   input.Name,
		Description:            input.Description,
		InstructionsMarkdown:   input.InstructionsMarkdown,
		ThumbnailUrl:           input.ThumbnailImageURL,
		CookTimeSeconds:        int64(input.CookTime.Seconds()),
		PreparationTimeSeconds: int64(input.PreparationTime.Seconds()),
		Source:                 input.Source,
		UpdatedAt: sql.NullTime{
			Valid: true,
			Time:  input.UpdatedAt,
		},
	}

	if err := s.queries.UpdateRecipe(ctx, params); err != nil {
		return fmt.Errorf("update recipe: %w", err)
	}

	return nil
}
