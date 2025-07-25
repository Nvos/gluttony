package recipe

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gluttony/x/markdown"
	"gluttony/x/pagination"
	"gluttony/x/sqlx"
	"log/slog"
	"mime/multipart"
	"slices"
	"time"
)

type Service struct {
	db           *pgxpool.Pool
	imageService ImageService
	searchIndex  Index
	store        Store
	markdown     *markdown.Markdown
	logger       *slog.Logger
}

func (s *Service) Stop() error {
	if err := s.searchIndex.Close(); err != nil {
		return fmt.Errorf("close search index: %w", err)
	}

	return nil
}

func NewService(
	db *pgxpool.Pool,
	store Store,
	imageService ImageService,
	searchIndex Index,
	logger *slog.Logger,
) (*Service, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	if store == nil {
		return nil, errors.New("store is nil")
	}

	if imageService == nil {
		return nil, errors.New("imageService is nil")
	}

	if searchIndex == nil {
		return nil, errors.New("searchIndex is nil")
	}

	return &Service{
		imageService: imageService,
		searchIndex:  searchIndex,
		store:        store,
		db:           db,
		markdown:     markdown.NewMarkdown(),
		logger:       logger,
	}, nil
}

func (s *Service) Create(ctx context.Context, input CreateInput) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin create recipe tx: %w", err)
	}

	var imageURL string
	defer func() {
		if err == nil {
			return
		}

		_ = tx.Rollback(ctx)
		if err := s.imageService.Delete(imageURL); err != nil {
			s.logger.Error(
				"Failed to delete image during recipe create rollback",
				slog.String("url", imageURL),
				slog.String("err", err.Error()),
			)
		}
	}()

	txStore := s.store.WithTx(tx)

	imageURL, imageID, err := s.createImage(ctx, txStore, input.ThumbnailImage)
	if err != nil {
		return fmt.Errorf("create image: %w", err)
	}

	params := CreateRecipe{
		Name:                 input.Name,
		Description:          input.Description,
		ThumbnailImageID:     imageID,
		Source:               input.Source,
		InstructionsMarkdown: input.Instructions,
		Servings:             input.Servings,
		PreparationTime:      input.PreparationTime,
		CookTime:             input.PreparationTime,
		OwnerID:              input.OwnerID,
	}
	recipeID, err := txStore.CreateRecipe(ctx, params)
	if err != nil {
		if sqlx.IsUniqueViolation(err) {
			return ErrUniqueName
		}

		return fmt.Errorf("create recipe: %w", err)
	}

	if err := txStore.CreateRecipeNutrition(ctx, recipeID, input.Nutrition); err != nil {
		return fmt.Errorf("create nutrition: %w", err)
	}

	if err := txStore.CreateRecipeTags(ctx, recipeID, input.Tags); err != nil {
		return fmt.Errorf("create recipe tags: %w", err)
	}

	if err := txStore.CreateRecipeIngredients(ctx, recipeID, input.Ingredients); err != nil {
		return fmt.Errorf("create recipe ingredients: %w", err)
	}

	indexRecipe := IndexRecipeInput{
		ID:          recipeID,
		Name:        input.Name,
		Description: input.Description,
	}

	if err := s.searchIndex.Index(indexRecipe); err != nil {
		return fmt.Errorf("index recipe: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit create recipe tx: %w", err)
	}

	return nil
}

func (s *Service) createImage(
	ctx context.Context,
	txStore Store,
	file *multipart.FileHeader,
) (string, *int32, error) {
	var imageID *int32
	if file == nil {
		return "", nil, nil
	}

	url, err := s.imageService.Upload(file)
	if err != nil {
		return "", nil, fmt.Errorf("upload thumbnail image: %w", err)
	}

	input := CreateRecipeImageInput{URL: url}
	gotID, err := txStore.CreateRecipeImage(ctx, input)
	if err != nil {
		return "", nil, fmt.Errorf("create recipe image: %w", err)
	}
	imageID = &gotID

	return url, imageID, nil
}

func (s *Service) updateTagsIfChanged(
	ctx context.Context,
	txStore Store,
	recipeID int32,
	current []Tag,
	incoming []string,
) error {
	if len(incoming) == 0 {
		return nil
	}

	tags := make([]string, 0, len(current))
	for i := range current {
		tags = append(tags, current[i].Name)
	}
	isTagsChanged := !slices.Equal(tags, incoming)
	if !isTagsChanged {
		return nil
	}

	if err := txStore.DeleteRecipeTags(ctx, recipeID); err != nil {
		return fmt.Errorf("delete recipe tags: %w", err)
	}

	if err := txStore.CreateRecipeTags(ctx, recipeID, incoming); err != nil {
		return fmt.Errorf("create tags: %w", err)
	}

	return nil
}

func (s *Service) updateIngredientsIfChanged(
	ctx context.Context,
	txStore Store,
	recipeID int32,
	current []Ingredient,
	incoming []Ingredient,
) error {
	if len(incoming) == 0 {
		return nil
	}

	isIngredientsChanged := !slices.EqualFunc(
		current, incoming,
		func(v1, v2 Ingredient) bool {
			return v1 == v2
		},
	)

	if !isIngredientsChanged {
		return nil
	}

	if err := txStore.DeleteRecipeIngredients(ctx, recipeID); err != nil {
		return fmt.Errorf("delete ingredients: %w", err)
	}

	if err := txStore.CreateRecipeIngredients(ctx, recipeID, incoming); err != nil {
		return fmt.Errorf("create ingredients: %w", err)
	}

	return nil
}

//nolint:nonamedreturns // used for rollback detection
func (s *Service) Update(ctx context.Context, input UpdateInput) (err error) {
	current, err := s.GetFull(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("get recipe by id=%v: %w", input.ID, err)
	}

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin update recipe tx: %w", err)
	}

	var imageURL string
	defer func() {
		if err == nil {
			return
		}

		_ = tx.Rollback(ctx)
		if err := s.imageService.Delete(imageURL); err != nil {
			s.logger.Error(
				"Failed to delete image during recipe update rollback",
				slog.String("url", imageURL),
				slog.String("err", err.Error()),
			)
		}
	}()

	txStore := s.store.WithTx(tx)

	err = s.updateTagsIfChanged(
		ctx,
		txStore,
		current.ID,
		current.Tags,
		input.Tags,
	)
	if err != nil {
		return fmt.Errorf("update recipe id=%v tags: %w", input.ID, err)
	}

	err = s.updateIngredientsIfChanged(
		ctx,
		txStore,
		current.ID,
		current.Ingredients,
		input.Ingredients,
	)
	if err != nil {
		return fmt.Errorf("update recipe id=%v ingredients : %w", input.ID, err)
	}

	if err = txStore.UpdateNutrition(ctx, input.ID, input.Nutrition); err != nil {
		return fmt.Errorf("update nutrition: %w", err)
	}

	thumbnailImageID := current.ThumbnailImageID

	imageURL, imageID, err := s.createImage(ctx, txStore, input.ThumbnailImage)
	if err != nil {
		return fmt.Errorf("create image: %w", err)
	}

	if imageID != nil {
		thumbnailImageID = imageID
	}

	err = txStore.UpdateRecipe(ctx, UpdateRecipe{
		CreateRecipe: CreateRecipe{
			Name:                 input.Name,
			Description:          input.Description,
			ThumbnailImageID:     thumbnailImageID,
			Source:               input.Source,
			InstructionsMarkdown: input.Instructions,
			Servings:             input.Servings,
			PreparationTime:      input.PreparationTime,
			CookTime:             input.CookTime,
			// TODO: Not used by update, likely need to separate model fully
			OwnerID: 0,
		},
		UpdatedAt: time.Now().UTC(),
		ID:        input.ID,
	})
	if err != nil {
		if sqlx.IsUniqueViolation(err) {
			return ErrUniqueName
		}

		return fmt.Errorf("update recipe: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit update recipe tx: %w", err)
	}

	return nil
}

func (s *Service) AllSummaries(
	ctx context.Context,
	input SearchInput,
) (pagination.Page[Summary], error) {
	out := pagination.Page[Summary]{
		TotalCount: 0,
		Rows:       nil,
	}
	if input.Search != "" {
		result, err := s.searchIndex.Search(ctx, input.Search, pagination.OffsetFromPage(input.Page))
		if err != nil {
			return out, fmt.Errorf("search index: %w", err)
		}

		if result.TotalCount == 0 {
			return out, nil
		}

		input.RecipeIDs = result.IDs
		out.TotalCount = result.TotalCount
	} else {
		count, err := s.store.CountRecipeSummaries(ctx)
		if err != nil {
			return pagination.Page[Summary]{}, fmt.Errorf("count recipe summaries: %w", err)
		}

		out.TotalCount = count
	}

	recipeSummaries, err := s.store.AllRecipeSummaries(ctx, input)
	if err != nil {
		return out, fmt.Errorf("all recipe summaries: %w", err)
	}

	tags, err := s.store.AllTagsByRecipeIDs(ctx, input.RecipeIDs)
	if err != nil {
		return out, fmt.Errorf("all recipe tags by ids=%+v: %w", input.RecipeIDs, err)
	}

	for i := range recipeSummaries {
		recipeTags, ok := tags[recipeSummaries[i].ID]
		if !ok {
			continue
		}

		recipeSummaries[i].Tags = recipeTags
	}

	out.Rows = recipeSummaries

	return out, nil
}

func (s *Service) GetFull(ctx context.Context, recipeID int32) (Recipe, error) {
	r, err := s.store.GetRecipe(ctx, recipeID)
	if err != nil {
		return Recipe{}, fmt.Errorf("get recipe id=%d: %w", recipeID, err)
	}

	html, err := s.markdown.ConvertToHTML(r.InstructionsMarkdown)
	if err != nil {
		return Recipe{}, fmt.Errorf("convert instructions to HTML: %w", err)
	}

	r.InstructionsHTML = html

	return r, nil
}
