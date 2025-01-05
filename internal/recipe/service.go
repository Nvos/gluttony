package recipe

import (
	"context"
	"database/sql"
	"fmt"
	"gluttony/internal/recipe/queries"
	"io"
	"time"
)

type Tag struct {
	ID    int
	Order int
	Name  string
}

type CreateInput struct {
	Name            string
	Description     string
	Source          string
	Instructions    string
	Servings        int8
	PreparationTime time.Duration
	CookTime        time.Duration
	Tags            []string
	Ingredients     []Ingredient
	Nutrition       Nutrition
	ThumbnailImage  io.Reader
}

type Service struct {
	db         *sql.DB
	queries    *queries.Queries
	mediaStore MediaStore
}

func NewService(db *sql.DB, mediaStore MediaStore) *Service {
	if db == nil {
		panic("db is nil")
	}

	if mediaStore == nil {
		panic("mediaStore is nil")
	}

	return &Service{queries: queries.New(db), db: db, mediaStore: mediaStore}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (err error) {
	thumbnailImageURL := ""
	if input.ThumbnailImage != nil {
		thumbnailImageURL, err = s.mediaStore.UploadImage(input.ThumbnailImage)
		if err != nil {
			return fmt.Errorf("upload thumbnail image: %w", err)
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = tx.Commit()
			return
		}

		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = fmt.Errorf("recipe create: %w: %w", err, rollbackErr)

			// TODO: remove image
		}
	}()

	txQueries := s.queries.WithTx(tx)

	ingredients := make([]Ingredient, len(input.Ingredients))
	for i := range input.Ingredients {
		id, err := txQueries.CreateIngredient(ctx, input.Ingredients[i].Name)
		if err != nil {
			return fmt.Errorf("create ingredient: %w", err)
		}

		ingredients[i].ID = int(id)
		ingredients[i].Name = input.Ingredients[i].Name
		ingredients[i].Order = int8(i)
		ingredients[i].Quantity = input.Ingredients[i].Quantity
		ingredients[i].Unit = input.Ingredients[i].Unit
	}

	tags := make([]Tag, len(input.Tags))
	for i := range input.Tags {
		id, err := txQueries.CreateTag(ctx, input.Tags[i])
		if err != nil {
			return fmt.Errorf("create tag: %w", err)
		}

		tags[i].Name = input.Tags[i]
		tags[i].ID = int(id)
		tags[i].Order = i
	}

	createRecipeParams := queries.CreateRecipeParams{
		Name:                   input.Name,
		Description:            input.Description,
		InstructionsMarkdown:   input.Instructions,
		CookTimeSeconds:        int64(input.CookTime.Seconds()),
		PreparationTimeSeconds: int64(input.PreparationTime.Seconds()),
		Source:                 input.Source,
	}
	if thumbnailImageURL != "" {
		createRecipeParams.ThumbnailUrl = sql.NullString{
			String: thumbnailImageURL,
			Valid:  true,
		}
	}

	recipeID, err := txQueries.CreateRecipe(ctx, createRecipeParams)
	if err != nil {
		return fmt.Errorf("create recipe: %w", err)
	}

	err = txQueries.CreateNutrition(ctx, queries.CreateNutritionParams{
		RecipeID: recipeID,
		Calories: float64(input.Nutrition.Calories),
		Fat:      float64(input.Nutrition.Fat),
		Carbs:    float64(input.Nutrition.Carbs),
		Protein:  float64(input.Nutrition.Protein),
	})
	if err != nil {
		return fmt.Errorf("create nutrition: %w", err)
	}

	for i := range tags {
		err = txQueries.CreateRecipeTag(ctx, queries.CreateRecipeTagParams{
			RecipeOrder: int64(tags[i].Order),
			RecipeID:    recipeID,
			TagID:       int64(tags[i].ID),
		})
		if err != nil {
			return fmt.Errorf("create recipe tag: %w", err)
		}
	}

	for i := range ingredients {
		err = txQueries.CreateRecipeIngredient(ctx, queries.CreateRecipeIngredientParams{
			RecipeOrder:  int64(ingredients[i].Order),
			RecipeID:     recipeID,
			IngredientID: int64(ingredients[i].ID),
			Unit:         ingredients[i].Unit,
			Quantity:     int64(ingredients[i].Quantity),
		})
	}

	return nil
}
