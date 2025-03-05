package recipe

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gluttony/internal/ingredient"
	"gluttony/internal/recipe/postgres"
	"time"
)

type Store struct {
	db      postgres.DBTX
	queries *postgres.Queries
}

func NewStore(db postgres.DBTX) *Store {
	return &Store{
		db:      db,
		queries: postgres.New(db),
	}
}

func (s *Store) WithTx(tx pgx.Tx) *Store {
	return &Store{
		db:      tx,
		queries: postgres.New(tx),
	}
}

func (s *Store) AllRecipeSummaries(
	ctx context.Context,
	input SearchInput,
) ([]Summary, error) {

	ids := make([]int32, 0, len(input.RecipeIDs))
	for i := range input.RecipeIDs {
		ids = append(ids, int32(i))
	}

	if len(ids) == 0 {
		ids = nil
	}

	rows, err := s.queries.AllRecipeSummaries(ctx, postgres.AllRecipeSummariesParams{
		Ids:    nil,
		Offset: int32(input.Page * input.Limit),
		Limit:  int32(input.Limit),
	})
	if err != nil {
		return nil, err
	}

	out := make([]Summary, 0, 20)
	for _, row := range rows {
		value := Summary{
			ID:                row.ID,
			Name:              row.Name,
			Description:       row.Description,
			ThumbnailImageURL: row.ThumbnailUrl,
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

func (s *Store) GetRecipe(ctx context.Context, ID int32) (Recipe, error) {
	recipe, err := s.queries.GetFullRecipe(ctx, ID)
	if err != nil {
		return Recipe{}, fmt.Errorf("get recipe id=%d: %w", ID, err)
	}

	tags, err := s.AllTagsByRecipeIDs(ctx, []int32{ID})
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

func (s *Store) AllTagsByRecipeIDs(ctx context.Context, recipeIDs []int32) (map[int32][]Tag, error) {
	tags, err := s.queries.AllRecipeTags(ctx, recipeIDs)
	if err != nil {
		return nil, fmt.Errorf("get all tags by recipe ids = %+v: %w", recipeIDs, err)
	}

	out := make(map[int32][]Tag, len(tags))
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
	recipeIDs ...int32,
) (map[int32][]Ingredient, error) {
	ingredients, err := s.queries.AllRecipeIngredients(ctx, recipeIDs)
	if err != nil {
		return nil, fmt.Errorf("get all ingredients by recipe ids = %+v: %w", recipeIDs, err)
	}

	out := make(map[int32][]Ingredient, len(ingredients))
	for i := range ingredients {
		if out[ingredients[i].RecipeID] == nil {
			out[ingredients[i].RecipeID] = []Ingredient{}
		}

		out[ingredients[i].RecipeID] = append(out[ingredients[i].RecipeID], Ingredient{
			Order:    int8(ingredients[i].RecipeOrder),
			Quantity: float32(ingredients[i].Quantity),
			Note:     ingredients[i].Note,
			Unit:     ingredients[i].Unit,
			Ingredient: ingredient.Ingredient{
				ID:   int(ingredients[i].ID),
				Name: ingredients[i].Name,
			},
		})
	}

	return out, nil
}

func (s *Store) CreateRecipeIngredients(
	ctx context.Context,
	recipeID int32,
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

	existingIngredientLookup := make(map[string]postgres.Ingredient, len(existingIngredients))
	for i := range existingIngredients {
		existingIngredientLookup[existingIngredients[i].Name] = existingIngredients[i]
	}

	savedIngredients := make([]Ingredient, len(ingredients))
	for i := range ingredients {
		var ingredientID int32
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
		savedIngredients[i].Note = ingredients[i].Note
	}

	for i := range savedIngredients {
		err = s.queries.CreateRecipeIngredient(ctx, postgres.CreateRecipeIngredientParams{
			RecipeOrder:  int32(savedIngredients[i].Order),
			RecipeID:     recipeID,
			IngredientID: int32(savedIngredients[i].Ingredient.ID),
			Unit:         savedIngredients[i].Unit,
			Quantity:     savedIngredients[i].Quantity,
			Note:         savedIngredients[i].Note,
		})
		if err != nil {
			return fmt.Errorf("create recipe ingredient: %w", err)
		}
	}

	return nil
}

func (s *Store) CreateRecipeNutrition(
	ctx context.Context,
	recipeID int32,
	nutrition Nutrition,
) error {
	err := s.queries.CreateNutrition(ctx, postgres.CreateNutritionParams{
		RecipeID: recipeID,
		Calories: nutrition.Calories,
		Fat:      nutrition.Fat,
		Carbs:    nutrition.Carbs,
		Protein:  nutrition.Protein,
	})
	if err != nil {
		return fmt.Errorf("create nutrition: %w", err)
	}

	return nil
}

func (s *Store) CreateRecipe(ctx context.Context, input CreateRecipe) (int32, error) {
	createRecipeParams := postgres.CreateRecipeParams{
		Name:                   input.Name,
		Description:            input.Description,
		InstructionsMarkdown:   input.InstructionsMarkdown,
		CookTimeSeconds:        int32(input.CookTime.Seconds()),
		PreparationTimeSeconds: int32(input.PreparationTime.Seconds()),
		Source:                 input.Source,
		ThumbnailUrl:           input.ThumbnailImageURL,
		OwnerID:                int32(input.OwnerID),
		Servings:               int32(input.Servings),
	}

	id, err := s.queries.CreateRecipe(ctx, createRecipeParams)
	if err != nil {
		return 0, fmt.Errorf("create recipe: %w", err)
	}

	return id, nil
}

func (s *Store) CreateRecipeTags(
	ctx context.Context,
	recipeID int32,
	tagNames []string,
) error {
	existingTags, err := s.queries.AllTagsByNames(ctx, tagNames)
	if err != nil {
		return fmt.Errorf("all savedTags by names: %w", err)
	}

	existingTagLookup := make(map[string]postgres.Tag, len(existingTags))
	for i := range existingTags {
		existingTagLookup[existingTags[i].Name] = existingTags[i]
	}

	savedTags := make([]Tag, len(tagNames))
	for i := range tagNames {
		var tagID int32
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
		err = s.queries.CreateRecipeTag(ctx, postgres.CreateRecipeTagParams{
			RecipeOrder: int32(savedTags[i].Order),
			RecipeID:    recipeID,
			TagID:       int32(savedTags[i].ID),
		})
		if err != nil {
			return fmt.Errorf("create recipe tag: %w", err)
		}
	}

	return nil
}

func (s *Store) DeleteRecipeTags(ctx context.Context, recipeID int32) error {
	if err := s.queries.DeleteRecipeTags(ctx, recipeID); err != nil {
		return fmt.Errorf("remove recipe tags: %w", err)
	}

	return nil
}

func (s *Store) DeleteRecipeIngredients(ctx context.Context, recipeID int32) error {
	if err := s.queries.DeleteRecipeIngredients(ctx, recipeID); err != nil {
		return fmt.Errorf("remove recipe ingredients: %w", err)
	}

	return nil
}

func (s *Store) UpdateNutrition(ctx context.Context, recipeID int32, nutrition Nutrition) error {
	err := s.queries.UpdateNutrition(ctx, postgres.UpdateNutritionParams{
		RecipeID: recipeID,
		Calories: nutrition.Calories,
		Fat:      nutrition.Fat,
		Carbs:    nutrition.Carbs,
		Protein:  nutrition.Protein,
	})
	if err != nil {
		return fmt.Errorf("update nutrition: %w", err)
	}

	return nil
}

func (s *Store) UpdateRecipe(ctx context.Context, input UpdateRecipe) error {
	params := postgres.UpdateRecipeParams{
		ID:                     int32(input.ID),
		Name:                   input.Name,
		Description:            input.Description,
		InstructionsMarkdown:   input.InstructionsMarkdown,
		ThumbnailUrl:           input.ThumbnailImageURL,
		CookTimeSeconds:        int32(input.CookTime.Seconds()),
		PreparationTimeSeconds: int32(input.PreparationTime.Seconds()),
		Source:                 input.Source,
		UpdatedAt:              &input.UpdatedAt,
		Servings:               int32(input.Servings),
	}

	if err := s.queries.UpdateRecipe(ctx, params); err != nil {
		return fmt.Errorf("update recipe: %w", err)
	}

	return nil
}
