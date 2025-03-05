package recipe

import (
	"gluttony/internal/ingredient"
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
	Limit     int32
}

type SearchResult struct {
	IsSearch   bool
	TotalCount uint32
	IDs        []int32
}

// TODO: change to bool (sqlite needed integer)
func (sr SearchResult) IsSearchCondition() int32 {
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
	OwnerID              int32
}

type UpdateRecipe struct {
	ID        int32
	UpdatedAt time.Time

	CreateRecipe
}
