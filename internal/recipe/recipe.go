package recipe

import (
	"context"
	"github.com/jackc/pgx/v5"
	"gluttony/internal/ingredient"
	"gluttony/pkg/pagination"
	"io"
	"time"
)

type Ingredient struct {
	Order    int8
	Quantity float32
	Note     string
	// TODO: unit enum
	Unit string

	ingredient.Ingredient
}

type Nutrition struct {
	Calories float32
	Fat      float32
	Carbs    float32
	Protein  float32
}

type Summary struct {
	ID                int32
	Name              string
	Description       string
	ThumbnailImageURL string
	Tags              []Tag
}

type Recipe struct {
	ID                   int32
	Name                 string
	Description          string
	ThumbnailImageURL    string
	Source               string
	InstructionsMarkdown string
	InstructionsHTML     string
	Servings             int8
	PreparationTime      time.Duration
	CookTime             time.Duration

	Tags        []Tag
	Ingredients []Ingredient
	Nutrition   Nutrition
}

type SearchInput struct {
	Search    string
	RecipeIDs []int32
	Page      int32
}

type SearchResult struct {
	TotalCount int64
	IDs        []int32
}

type Tag struct {
	ID    int32
	Order int32
	Name  string
}

type UpdateInput struct {
	ID int32
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
	OwnerID         int32
}

type CreateRecipe struct {
	Name                 string
	Description          string
	ThumbnailImageURL    string
	Source               string
	InstructionsMarkdown string
	Servings             int8
	PreparationTime      time.Duration
	CookTime             time.Duration
	OwnerID              int32
}

type UpdateRecipe struct {
	ID        int32
	UpdatedAt time.Time

	CreateRecipe
}

type MediaStore interface {
	UploadImage(file io.Reader) (string, error)
}

type Store interface {
	WithTx(tx pgx.Tx) Store
	//TODO: verify search input
	AllRecipeSummaries(ctx context.Context, input SearchInput) ([]Summary, error)
	GetRecipe(ctx context.Context, id int32) (Recipe, error)
	AllTagsByRecipeIDs(ctx context.Context, recipeIDs []int32) (map[int32][]Tag, error)
	AllIngredientsByRecipeIDs(
		ctx context.Context,
		recipeIDs ...int32,
	) (map[int32][]Ingredient, error)
	CreateRecipeIngredients(
		ctx context.Context,
		recipeID int32,
		ingredients []Ingredient,
	) error
	CreateRecipeNutrition(
		ctx context.Context,
		recipeID int32,
		nutrition Nutrition,
	) error
	CreateRecipe(ctx context.Context, input CreateRecipe) (int32, error)
	CreateRecipeTags(
		ctx context.Context,
		recipeID int32,
		tagNames []string,
	) error
	DeleteRecipeTags(ctx context.Context, recipeID int32) error
	DeleteRecipeIngredients(ctx context.Context, recipeID int32) error
	UpdateNutrition(ctx context.Context, recipeID int32, nutrition Nutrition) error
	UpdateRecipe(ctx context.Context, input UpdateRecipe) error
}

type Index interface {
	Index(value Recipe) error
	Search(ctx context.Context, query string, offset pagination.Offset) (SearchResult, error)
	Close() error
}
