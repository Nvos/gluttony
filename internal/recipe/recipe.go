package recipe

import (
	"gluttony/internal/ingredient"
	"gluttony/internal/templating"
	"io"
	"log/slog"
	"time"
)

type Ingredient struct {
	ingredient.Ingredient

	Order    int8
	Quantity float32
	// TODO: unit enum
	Unit string
}

type Nutrition struct {
	Calories float32
	Fat      float32
	Carbs    float32
	Protein  float32
}

type Partial struct {
	ID                int
	Name              string
	Description       string
	ThumbnailImageURL string
	Tags              []Tag
}

type Full struct {
	ID                int
	Name              string
	Description       string
	ThumbnailImageURL string
	Tags              []Tag
	Source            string
	Instructions      string
	Servings          int8
	PreparationTime   time.Duration
	CookTime          time.Duration
	Ingredients       []Ingredient
	Nutrition         Nutrition
}

type SearchInput struct {
	Query string
}

type MediaStore interface {
	UploadImage(file io.Reader) (string, error)
}

type Deps struct {
	service    *Service
	logger     *slog.Logger
	templates  *templating.Templating
	mediaStore MediaStore
}

func NewDeps(
	service *Service,
	logger *slog.Logger,
	templateManager *templating.Templating,
	mediaStore MediaStore,
) *Deps {
	if logger == nil {
		panic("logger must not be nil")
	}

	if service == nil {
		panic("service must not be nil")
	}

	if templateManager == nil {
		panic("templateManager must not be nil")
	}

	if mediaStore == nil {
		panic("mediaStore must not be nil")
	}

	return &Deps{
		service:    service,
		logger:     logger,
		templates:  templateManager,
		mediaStore: mediaStore,
	}
}

type CreateRecipeInput struct {
	Name      string
	Steps     string
	CreatedAt time.Time
}
