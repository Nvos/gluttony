package recipe

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/blevesearch/bleve"
	"io"
	"log/slog"
	"slices"
	"strconv"
	"time"
)

type indexEntry struct {
	ID          int64
	Name        string
	Description string
}

type Tag struct {
	ID    int
	Order int
	Name  string
}

type UpdateInput struct {
	ID int64
	CreateInput
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
	ThumbnailURL    string
}

type Service struct {
	db          *sql.DB
	mediaStore  MediaStore
	searchIndex bleve.Index
	store       *Store
	markdown    *Markdown
	logger      *slog.Logger
}

func (s *Service) Stop() error {
	return s.searchIndex.Close()
}

func NewService(
	db *sql.DB,
	mediaStore MediaStore,
	searchIndex bleve.Index,
	logger *slog.Logger,
) (*Service, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	if mediaStore == nil {
		return nil, fmt.Errorf("mediaStore is nil")
	}

	if searchIndex == nil {
		return nil, fmt.Errorf("searchIndex is nil")
	}

	return &Service{
		db:          db,
		mediaStore:  mediaStore,
		searchIndex: searchIndex,
		store:       NewStore(db),
		markdown:    NewMarkdown(),
		logger:      logger,
	}, nil
}

// TODO: move to index
func (s *Service) index(value indexEntry) error {
	err := s.searchIndex.Index(strconv.Itoa(int(value.ID)), value)
	if err != nil {
		return fmt.Errorf("search index failed: %w", err)
	}

	return nil
}

// TODO: move to index
func (s *Service) indexAll(ctx context.Context) error {
	partial, err := s.store.AllRecipeSummaries(ctx, SearchInput{
		Limit: 9999,
		Page:  0,
	})
	if err != nil {
		return err
	}

	batch := s.searchIndex.NewBatch()
	for i := range partial {
		value := indexEntry{
			ID:          int64(partial[i].ID),
			Name:        partial[i].Name,
			Description: partial[i].Description,
		}

		if err := batch.Index(strconv.Itoa(int(value.ID)), value); err != nil {
			return fmt.Errorf("index recipe, batch add: %w", err)
		}
	}

	if err := s.searchIndex.Batch(batch); err != nil {
		return fmt.Errorf("index recipes execute batch: %w", err)
	}

	return nil
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
		err := tx.Rollback()
		if errors.Is(err, sql.ErrTxDone) {
			return
		}

		if err != nil {
			s.logger.Error("Rolling back transaction", slog.String("err", err.Error()))
		}

		// TODO: remove image
	}()

	txStore := s.store.WithTx(tx)

	params := CreateRecipe{
		Name:                 input.Name,
		Description:          input.Description,
		ThumbnailImageURL:    thumbnailImageURL,
		Source:               input.Source,
		InstructionsMarkdown: input.Instructions,
		Servings:             input.Servings,
		PreparationTime:      input.PreparationTime,
		CookTime:             input.PreparationTime,
	}
	recipeID, err := txStore.CreateRecipe(ctx, params)
	if err != nil {
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

	err = s.index(indexEntry{
		ID:          recipeID,
		Name:        input.Name,
		Description: input.Description,
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}

// TODO: move to index
func (s *Service) search(input SearchInput) (SearchResult, error) {
	if input.Search == "" {
		return SearchResult{}, nil
	}

	query := bleve.NewQueryStringQuery(input.Search)
	offset := input.Page * input.Limit
	searchRequest := bleve.NewSearchRequestOptions(query, int(input.Limit), int(offset), false)
	searchResult, err := s.searchIndex.Search(searchRequest)
	if err != nil {
		return SearchResult{}, fmt.Errorf("search recipe index: %w", err)
	}
	recipeIDs := make([]int64, len(searchResult.Hits))
	for i := range searchResult.Hits {
		id, err := strconv.ParseInt(searchResult.Hits[i].ID, 10, 64)
		if err != nil {
			// All indexed ids are int64, any other id is unexpected and shouldn't happen
			panic(fmt.Sprintf("unexpected recipe index id: %v", err))
		}

		recipeIDs[i] = id
	}

	return SearchResult{
		IsSearch:   true,
		TotalCount: searchResult.Total,
		IDs:        recipeIDs,
	}, nil
}

func (s *Service) Update(ctx context.Context, input UpdateInput) error {
	current, err := s.GetFull(ctx, input.ID)
	if err != nil {
		return err
	}

	tags := make([]string, 0, len(current.Tags))
	for i := range current.Tags {
		tags = append(tags, current.Tags[i].Name)
	}
	isTagsChanged := !slices.Equal(tags, input.Tags)

	isIngredientsChanged := !slices.EqualFunc(
		current.Ingredients, input.Ingredients,
		func(i Ingredient, i2 Ingredient) bool {
			return i.ID == i2.ID
		},
	)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		err := tx.Rollback()
		if errors.Is(err, sql.ErrTxDone) {
			return
		}

		if err != nil {
			s.logger.Error("Rolling back transaction", slog.String("err", err.Error()))
		}

		// TODO: remove image
	}()

	txStore := s.store.WithTx(tx)
	if isTagsChanged && len(input.Tags) > 0 {
		if err := txStore.DeleteRecipeTags(ctx, input.ID); err != nil {
			return fmt.Errorf("delete recipe id=%d tags: %w", input.ID, err)
		}

		if err := txStore.CreateRecipeTags(ctx, input.ID, input.Tags); err != nil {
			return fmt.Errorf("create tags: %w", err)
		}
	}

	if isIngredientsChanged && len(input.Ingredients) > 0 {
		if err := txStore.DeleteRecipeIngredients(ctx, input.ID); err != nil {
			return fmt.Errorf("delete recipe id=%d ingredients: %w", input.ID, err)
		}

		if err := txStore.CreateRecipeIngredients(ctx, input.ID, input.Ingredients); err != nil {
			return fmt.Errorf("create ingredients: %w", err)
		}
	}

	if err = txStore.UpdateNutrition(ctx, input.ID, input.Nutrition); err != nil {
		return fmt.Errorf("update nutrition: %w", err)
	}

	thumbnailImageURL := current.ThumbnailImageURL
	if input.ThumbnailImage != nil {
		thumbnailImageURL, err = s.mediaStore.UploadImage(input.ThumbnailImage)
		if err != nil {
			return fmt.Errorf("upload thumbnail image: %w", err)
		}

		// TODO: remove previous image
	}

	err = txStore.UpdateRecipe(ctx, UpdateRecipe{
		CreateRecipe: CreateRecipe{
			Name:                 input.Name,
			Description:          input.Description,
			ThumbnailImageURL:    thumbnailImageURL,
			Source:               input.Source,
			InstructionsMarkdown: input.Instructions,
			Servings:             input.Servings,
			PreparationTime:      input.PreparationTime,
			CookTime:             input.CookTime,
		},
		UpdatedAt: time.Now().UTC(),
		ID:        input.ID,
	})
	if err != nil {
		return fmt.Errorf("update recipe: %w", err)
	}

	return tx.Commit()
}

func (s *Service) AllSummaries(ctx context.Context, input SearchInput) ([]Summary, error) {
	searchResult, err := s.search(input)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	if searchResult.IsSearch && searchResult.TotalCount == 0 {
		return nil, nil
	}

	recipeSummaries, err := s.store.AllRecipeSummaries(ctx, input)
	if err != nil {
		return nil, err
	}

	tags, err := s.store.AllTagsByRecipeIDs(ctx, searchResult.IDs...)
	if err != nil {
		return []Summary{}, err
	}

	for i := range recipeSummaries {
		recipeTags, ok := tags[recipeSummaries[i].ID]
		if !ok {
			continue
		}

		recipeSummaries[i].Tags = recipeTags
	}

	return recipeSummaries, nil
}

func (s *Service) GetFull(ctx context.Context, recipeID int64) (Recipe, error) {
	recipe, err := s.store.GetRecipe(ctx, recipeID)
	if err != nil {
		return Recipe{}, fmt.Errorf("get recipe id=%d: %w", recipeID, err)
	}

	html, err := s.markdown.ConvertToHTML(recipe.InstructionsMarkdown)
	if err != nil {
		return Recipe{}, fmt.Errorf("convert instructions to HTML: %w", err)
	}

	recipe.InstructionsHTML = html

	return recipe, nil
}
