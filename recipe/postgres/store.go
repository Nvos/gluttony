package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gluttony/ingredient"
	"gluttony/recipe"
	"gluttony/x/pagination"
	"gluttony/x/sqlx"
	"time"
)

var _ recipe.Store = (*Store)(nil)

type Store struct {
	queries *Queries
}

func NewStore(db DBTX) *Store {
	return &Store{
		queries: New(db),
	}
}

//nolint:ireturn // has to return self via interface in order to conform to interface
func (s *Store) WithTx(tx pgx.Tx) recipe.Store {
	return &Store{
		queries: New(tx),
	}
}

func (s *Store) CountRecipeSummaries(ctx context.Context) (int64, error) {
	summaries, err := s.queries.CountRecipeSummaries(ctx)
	if err != nil {
		return 0, fmt.Errorf("count recipe summaries: %w", err)
	}

	return summaries, nil
}

func (s *Store) AllRecipeSummaries(
	ctx context.Context,
	input recipe.SearchInput,
) ([]recipe.Summary, error) {
	params := AllRecipeSummariesParams{
		Offset: 0,
		Ids:    nil,
		Limit:  pagination.Limit,
	}
	if len(input.RecipeIDs) > 0 {
		params.Ids = input.RecipeIDs
	} else {
		params.Offset = input.Page * pagination.Limit
	}

	rows, err := s.queries.AllRecipeSummaries(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("fetch all recipe summaries: %w", err)
	}

	out := make([]recipe.Summary, 0, len(rows))
	lookup := make(map[int32]recipe.Summary, len(rows))
	for _, row := range rows {
		value := recipe.Summary{
			ID:                row.ID,
			Name:              row.Name,
			Description:       row.Description,
			ThumbnailImageURL: "",
			Tags:              []recipe.Tag{},
		}

		if row.Url != nil {
			value.ThumbnailImageURL = *row.Url
		}

		out = append(out, value)
		lookup[row.ID] = value
	}

	for i := range input.RecipeIDs {
		out[i] = lookup[input.RecipeIDs[i]]
	}

	return out, nil
}

func (s *Store) GetRecipe(ctx context.Context, id int32) (recipe.Recipe, error) {
	r, err := s.queries.GetFullRecipe(ctx, id)
	if err != nil {
		return recipe.Recipe{}, fmt.Errorf("get recipe id=%d: %w", id, err)
	}

	tags, err := s.AllTagsByRecipeIDs(ctx, []int32{id})
	if err != nil {
		return recipe.Recipe{}, fmt.Errorf("get all recipe tags for recipe id=%d: %w", id, err)
	}

	ingredients, err := s.AllIngredientsByRecipeIDs(ctx, id)
	if err != nil {
		return recipe.Recipe{}, fmt.Errorf("get all ingredients for recipe id=%d: %w", id, err)
	}

	thumbnailImageURL := ""
	if r.Url != nil {
		thumbnailImageURL = *r.Url
	}

	return recipe.Recipe{
		ID:                   id,
		Name:                 r.Name,
		Description:          r.Description,
		InstructionsMarkdown: r.InstructionsMarkdown,
		ThumbnailImageURL:    thumbnailImageURL,
		ThumbnailImageID:     r.ThumbnailID,
		Tags:                 tags[id],
		Source:               r.Source,
		//nolint:gosec // conversion is fine
		Servings:         int8(r.Servings),
		PreparationTime:  time.Duration(r.PreparationTimeSeconds) * time.Second,
		CookTime:         time.Duration(r.CookTimeSeconds) * time.Second,
		Ingredients:      ingredients[id],
		InstructionsHTML: "",
		Nutrition: recipe.Nutrition{
			Calories: r.Calories,
			Fat:      r.Fat,
			Carbs:    r.Carbs,
			Protein:  r.Protein,
		},
	}, nil
}

func (s *Store) AllTagsByRecipeIDs(ctx context.Context, recipeIDs []int32) (map[int32][]recipe.Tag, error) {
	tags, err := s.queries.AllRecipeTags(ctx, recipeIDs)
	if err != nil {
		return nil, fmt.Errorf("get all tags by recipe ids = %+v: %w", recipeIDs, err)
	}

	out := make(map[int32][]recipe.Tag, len(tags))
	for i := range tags {
		if out[tags[i].RecipeID] == nil {
			out[tags[i].RecipeID] = []recipe.Tag{}
		}

		out[tags[i].RecipeID] = append(out[tags[i].RecipeID], recipe.Tag{
			ID:    tags[i].ID,
			Order: tags[i].RecipeOrder,
			Name:  tags[i].Name,
		})
	}

	return out, nil
}

func (s *Store) AllIngredientsByRecipeIDs(
	ctx context.Context,
	recipeIDs ...int32,
) (map[int32][]recipe.Ingredient, error) {
	ingredients, err := s.queries.AllRecipeIngredients(ctx, recipeIDs)
	if err != nil {
		return nil, fmt.Errorf("get all ingredients by recipe ids = %+v: %w", recipeIDs, err)
	}

	out := make(map[int32][]recipe.Ingredient, len(ingredients))
	for i := range ingredients {
		if out[ingredients[i].RecipeID] == nil {
			out[ingredients[i].RecipeID] = []recipe.Ingredient{}
		}

		out[ingredients[i].RecipeID] = append(out[ingredients[i].RecipeID], recipe.Ingredient{
			//nolint:gosec // conversion is fine, there's no reason for more than 127 ingredients
			Order:    int8(ingredients[i].RecipeOrder),
			Quantity: ingredients[i].Quantity,
			Note:     ingredients[i].Note,
			Unit:     ingredients[i].Unit,
			Ingredient: ingredient.Ingredient{
				ID:   ingredients[i].ID,
				Name: ingredients[i].Name,
			},
		})
	}

	return out, nil
}

func (s *Store) CreateRecipeIngredients(
	ctx context.Context,
	recipeID int32,
	ingredients []recipe.Ingredient,
) error {
	names := make([]string, len(ingredients))
	for i := range ingredients {
		names[i] = ingredients[i].Name
	}

	existingIngredients, err := s.queries.AllIngredientsByNames(ctx, names)
	if err != nil {
		return fmt.Errorf("all ingredients by names: %w", err)
	}

	existingIngredientLookup := make(map[string]Ingredient, len(existingIngredients))
	for i := range existingIngredients {
		existingIngredientLookup[existingIngredients[i].Name] = existingIngredients[i]
	}

	savedIngredients := make([]recipe.Ingredient, len(ingredients))
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

		savedIngredients[i].ID = ingredientID
		savedIngredients[i].Name = ingredients[i].Name
		//nolint:gosec // conversion is fine, there's no reason for more than 127 ingredients
		savedIngredients[i].Order = int8(i)
		savedIngredients[i].Quantity = ingredients[i].Quantity
		savedIngredients[i].Unit = ingredients[i].Unit
		savedIngredients[i].Note = ingredients[i].Note
	}

	for i := range savedIngredients {
		err = s.queries.CreateRecipeIngredient(ctx, CreateRecipeIngredientParams{
			RecipeOrder:  int32(savedIngredients[i].Order),
			RecipeID:     recipeID,
			IngredientID: savedIngredients[i].ID,
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
	nutrition recipe.Nutrition,
) error {
	err := s.queries.CreateNutrition(ctx, CreateNutritionParams{
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

func (s *Store) CreateRecipe(ctx context.Context, input recipe.CreateRecipe) (int32, error) {
	createRecipeParams := CreateRecipeParams{
		Name:                   input.Name,
		Description:            input.Description,
		InstructionsMarkdown:   input.InstructionsMarkdown,
		CookTimeSeconds:        int32(input.CookTime.Seconds()),
		PreparationTimeSeconds: int32(input.PreparationTime.Seconds()),
		Source:                 input.Source,
		ThumbnailID:            input.ThumbnailImageID,
		OwnerID:                input.OwnerID,
		Servings:               int32(input.Servings),
	}

	id, err := s.queries.CreateRecipe(ctx, createRecipeParams)
	if err != nil {
		return 0, fmt.Errorf("create recipe: %w", sqlx.TransformSQLError(err))
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

	existingTagLookup := make(map[string]Tag, len(existingTags))
	for i := range existingTags {
		existingTagLookup[existingTags[i].Name] = existingTags[i]
	}

	savedTags := make([]recipe.Tag, len(tagNames))
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
		savedTags[i].ID = tagID
		//nolint:gosec // not realistic for this to overflow
		savedTags[i].Order = int32(i)
	}

	for i := range savedTags {
		err = s.queries.CreateRecipeTag(ctx, CreateRecipeTagParams{
			RecipeOrder: savedTags[i].Order,
			RecipeID:    recipeID,
			TagID:       savedTags[i].ID,
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

func (s *Store) UpdateNutrition(ctx context.Context, recipeID int32, nutrition recipe.Nutrition) error {
	err := s.queries.UpdateNutrition(ctx, UpdateNutritionParams{
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

func (s *Store) UpdateRecipe(ctx context.Context, input recipe.UpdateRecipe) error {
	params := UpdateRecipeParams{
		ID:                     input.ID,
		Name:                   input.Name,
		Description:            input.Description,
		InstructionsMarkdown:   input.InstructionsMarkdown,
		ThumbnailID:            input.ThumbnailImageID,
		CookTimeSeconds:        int32(input.CookTime.Seconds()),
		PreparationTimeSeconds: int32(input.PreparationTime.Seconds()),
		Source:                 input.Source,
		UpdatedAt:              &input.UpdatedAt,
		Servings:               int32(input.Servings),
	}

	if err := s.queries.UpdateRecipe(ctx, params); err != nil {
		return fmt.Errorf("update recipe: %w", sqlx.TransformSQLError(err))
	}

	return nil
}

func (s *Store) CreateRecipeImage(ctx context.Context, input recipe.CreateRecipeImageInput) (int32, error) {
	id, err := s.queries.CreateRecipeImage(ctx, input.URL)
	if err != nil {
		return 0, err
	}

	return id, nil
}
