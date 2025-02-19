package recipe

import (
	"gluttony/internal/ingredient"
	"io"
	"time"
)

type Ingredient struct {
	ID       int64
	Order    int8
	Quantity float32
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
	ID                int64
	Name              string
	Description       string
	ThumbnailImageURL string
	Tags              []Tag
}

type Recipe struct {
	ID                   int64
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
	RecipeIDs []int64
	Page      int64
	Limit     int64
}

type SearchResult struct {
	IsSearch   bool
	TotalCount uint64
	IDs        []int64
}

func (sr SearchResult) IsSearchCondition() int64 {
	if sr.IsSearch {
		return 1
	}

	return 0
}

type MediaStore interface {
	UploadImage(file io.Reader) (string, error)
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
	OwnerID              int64
}

type UpdateRecipe struct {
	ID        int64
	UpdatedAt time.Time

	CreateRecipe
}
