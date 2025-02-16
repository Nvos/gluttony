package recipe

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/blevesearch/bleve"
	"gluttony/internal/ingredient"
	"gluttony/internal/recipe/queries"
	"io"
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
}

type Service struct {
	db          *sql.DB
	queries     *queries.Queries
	mediaStore  MediaStore
	searchIndex bleve.Index
	store       *Store
	markdown    *Markdown
}

func (s *Service) Stop() error {
	return s.searchIndex.Close()
}

func NewService(db *sql.DB, mediaStore MediaStore, searchIndex bleve.Index) (*Service, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	if mediaStore == nil {
		return nil, fmt.Errorf("mediaStore is nil")
	}

	store := &Store{db: db}

	return &Service{
		queries:     queries.New(db),
		db:          db,
		mediaStore:  mediaStore,
		searchIndex: searchIndex,
		store:       store,
		markdown:    NewMarkdown(),
	}, nil
}

func (s *Service) index(value indexEntry) error {
	err := s.searchIndex.Index(strconv.Itoa(int(value.ID)), value)
	if err != nil {
		return fmt.Errorf("search index failed: %w", err)
	}

	return nil
}

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

	createRecipeParams := queries.CreateRecipeParams{
		Name:                   input.Name,
		Description:            input.Description,
		InstructionsMarkdown:   input.Instructions,
		CookTimeSeconds:        int64(input.CookTime.Seconds()),
		PreparationTimeSeconds: int64(input.PreparationTime.Seconds()),
		Source:                 input.Source,
		ThumbnailUrl:           thumbnailImageURL,
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

	if err := s.createRecipeTags(ctx, txQueries, recipeID, input.Tags); err != nil {
		return fmt.Errorf("create recipe tags: %w", err)
	}

	if err := s.createRecipeIngredients(ctx, txQueries, recipeID, input.Ingredients); err != nil {
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

	return nil
}

func (s *Service) createRecipeIngredients(
	ctx context.Context,
	txQueries *queries.Queries,
	recipeID int64,
	recipeIngredients []Ingredient,
) error {
	ingredients, err := s.createOrGetIngredients(ctx, txQueries, recipeIngredients)
	if err != nil {
		return fmt.Errorf("create ingredients: %w", err)
	}
	for i := range ingredients {
		err = txQueries.CreateRecipeIngredient(ctx, queries.CreateRecipeIngredientParams{
			RecipeOrder:  int64(ingredients[i].Order),
			RecipeID:     recipeID,
			IngredientID: int64(ingredients[i].Ingredient.ID),
			Unit:         ingredients[i].Unit,
			Quantity:     int64(ingredients[i].Quantity),
		})
	}

	return nil
}

func (s *Service) createRecipeTags(
	ctx context.Context,
	txQueries *queries.Queries,
	recipeID int64,
	tagNames []string,
) error {
	tags, err := s.createOrGetTags(ctx, txQueries, tagNames)
	if err != nil {
		return fmt.Errorf("create tags: %w", err)
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

	return nil
}

func (s *Service) createOrGetTags(
	ctx context.Context,
	txQueries *queries.Queries,
	tags []string,
) ([]Tag, error) {

	existing, err := txQueries.AllTagsByNames(ctx, tags)
	if err != nil {
		return nil, fmt.Errorf("get ingredients: %w", err)
	}

	existingLookup := make(map[string]queries.Tag, len(existing))
	for i := range existing {
		existingLookup[existing[i].Name] = existing[i]
	}

	out := make([]Tag, len(tags))
	for i := range tags {
		var tagID int64
		if value, ok := existingLookup[tags[i]]; ok {
			tagID = value.ID
		} else {
			id, err := txQueries.CreateTag(ctx, tags[i])
			if err != nil {
				return nil, fmt.Errorf("create tag: %w", err)
			}

			tagID = id
		}

		out[i].Name = tags[i]
		out[i].ID = int(tagID)
		out[i].Order = i
	}

	return out, nil
}

func (s *Service) createOrGetIngredients(
	ctx context.Context,
	txQueries *queries.Queries,
	ingredients []Ingredient,
) ([]Ingredient, error) {

	names := make([]string, len(ingredients))
	for i := range ingredients {
		names[i] = ingredients[i].Name
	}

	existingIngredients, err := txQueries.AllIngredientsByNames(ctx, names)
	if err != nil {
		return nil, fmt.Errorf("get ingredients: %w", err)
	}

	existingLookup := make(map[string]queries.Ingredient, len(existingIngredients))
	for i := range existingIngredients {
		existingLookup[existingIngredients[i].Name] = existingIngredients[i]
	}

	out := make([]Ingredient, len(ingredients))
	for i := range ingredients {
		var ingredientID int64
		if value, ok := existingLookup[ingredients[i].Name]; ok {
			ingredientID = value.ID
		} else {
			id, err := txQueries.CreateIngredient(ctx, ingredients[i].Name)
			if err != nil {
				return nil, fmt.Errorf("create ingredient: %w", err)
			}

			ingredientID = id
		}

		out[i].Ingredient.ID = int(ingredientID)
		out[i].Ingredient.Name = ingredients[i].Name
		out[i].Order = int8(i)
		out[i].Quantity = ingredients[i].Quantity
		out[i].Unit = ingredients[i].Unit
	}

	return out, nil
}

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

	var isIngredientsChanged = !slices.EqualFunc(
		current.Ingredients, input.Ingredients,
		func(i Ingredient, i2 Ingredient) bool {
			return i.ID == i2.ID
		},
	)

	// TODO: tx
	txQueries := s.queries
	if isTagsChanged && len(input.Tags) > 0 {
		if err := txQueries.DeleteRecipeTags(ctx, input.ID); err != nil {
			return fmt.Errorf("delete recipe id=%d tags: %w", input.ID, err)
		}

		if err := s.createRecipeTags(ctx, txQueries, input.ID, input.Tags); err != nil {
			return fmt.Errorf("create tags: %w", err)
		}
	}
	if isIngredientsChanged && len(input.Ingredients) > 0 {
		if err := txQueries.DeleteRecipeIngredients(ctx, input.ID); err != nil {
			return fmt.Errorf("delete recipe id=%d ingredients: %w", input.ID, err)
		}

		if err := s.createRecipeIngredients(ctx, txQueries, input.ID, input.Ingredients); err != nil {
			return fmt.Errorf("create ingredients: %w", err)
		}
	}

	err = txQueries.UpdateNutrition(ctx, queries.UpdateNutritionParams{
		Calories: float64(input.Nutrition.Calories),
		Fat:      float64(input.Nutrition.Fat),
		Carbs:    float64(input.Nutrition.Carbs),
		Protein:  float64(input.Nutrition.Protein),
		RecipeID: input.ID,
	})
	if err != nil {
		return fmt.Errorf("update nutrition: %w", err)
	}

	thumbnailImageURL := current.ThumbnailImageURL
	if input.ThumbnailImage != nil {
		thumbnailImageURL, err = s.mediaStore.UploadImage(input.ThumbnailImage)
		if err != nil {
			return fmt.Errorf("upload thumbnail image: %w", err)
		}
	}

	err = txQueries.UpdateRecipe(ctx, queries.UpdateRecipeParams{
		Name:                   input.Name,
		Description:            input.Description,
		InstructionsMarkdown:   input.Instructions,
		ThumbnailUrl:           thumbnailImageURL,
		CookTimeSeconds:        int64(input.CookTime),
		PreparationTimeSeconds: int64(input.PreparationTime),
		Source:                 input.Source,
		UpdatedAt: sql.NullTime{
			Valid: true,
			Time:  time.Now().UTC(),
		},
		ID: input.ID,
	})
	if err != nil {
		return fmt.Errorf("update recipe: %w", err)
	}

	return nil
}

func (s *Service) AllSummaries(ctx context.Context, input SearchInput) ([]Summary, error) {
	searchResult, err := s.search(input)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	if searchResult.IsSearch && searchResult.TotalCount == 0 {
		return nil, nil
	}

	recipes, err := s.store.AllRecipeSummaries(ctx, SearchInput{
		Page:      input.Page,
		Limit:     input.Limit,
		RecipeIDs: searchResult.IDs,
	})
	if err != nil {
		return []Summary{}, fmt.Errorf("get all recipes: %w", err)
	}

	if searchResult.IsSearch {
		sorted := make([]Summary, len(recipes))
		for i := range searchResult.IDs {
			for j := range recipes {
				if int64(recipes[j].ID) == searchResult.IDs[i] {
					sorted[i] = recipes[j]
					break
				}
			}
		}
		recipes = sorted
	}

	tags, err := s.allTagsByRecipeIDs(ctx, searchResult.IDs)
	if err != nil {
		return []Summary{}, err
	}

	for i := range recipes {
		recipeTags, ok := tags[int64(recipes[i].ID)]
		if !ok {
			continue
		}

		recipes[i].Tags = recipeTags
	}

	return recipes, nil
}

func (s *Service) GetFull(ctx context.Context, recipeID int64) (Full, error) {
	recipe, err := s.queries.GetFullRecipe(ctx, recipeID)
	if err != nil {
		return Full{}, fmt.Errorf("get recipe id=%d: %w", recipeID, err)
	}

	tags, err := s.allTagsByRecipeIDs(ctx, []int64{recipeID})
	if err != nil {
		return Full{}, fmt.Errorf("get all recipe tags for recipe id=%d: %w", recipeID, err)
	}

	ingredients, err := s.allIngredientsByRecipeIDs(ctx, []int64{recipeID})
	if err != nil {
		return Full{}, fmt.Errorf("get all ingredients for recipe id=%d: %w", recipeID, err)
	}

	html, err := s.markdown.ConvertToHTML(recipe.InstructionsMarkdown)
	if err != nil {
		return Full{}, fmt.Errorf("convert instructions to HTML: %w", err)
	}

	out := Full{
		ID:                   int(recipeID),
		Name:                 recipe.Name,
		Description:          recipe.Description,
		InstructionsMarkdown: recipe.InstructionsMarkdown,
		InstructionsHTML:     html,
		ThumbnailImageURL:    recipe.ThumbnailUrl,
		Tags:                 tags[recipeID],
		Source:               recipe.Source,
		Servings:             int8(recipe.Servings),
		PreparationTime:      time.Duration(recipe.PreparationTimeSeconds) * time.Second,
		CookTime:             time.Duration(recipe.CookTimeSeconds) * time.Second,
		Ingredients:          ingredients[recipeID],
		Nutrition: Nutrition{
			Calories: float32(recipe.Calories),
			Fat:      float32(recipe.Fat),
			Carbs:    float32(recipe.Carbs),
			Protein:  float32(recipe.Protein),
		},
	}

	return out, nil
}

func (s *Service) allIngredientsByRecipeIDs(
	ctx context.Context,
	recipeIDs []int64,
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

func (s *Service) allTagsByRecipeIDs(ctx context.Context, recipeIDs []int64) (map[int64][]Tag, error) {
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
