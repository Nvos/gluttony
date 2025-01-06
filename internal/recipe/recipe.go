package recipe

import (
	"gluttony/internal/templates"
	"log/slog"
	"time"
)

type Partial struct {
	ID                int
	Name              string
	Description       string
	ThumbnailImageURL string
	Tags              []Tag
}

type SearchInput struct {
	Query string
}

type Deps struct {
	service    *Service
	logger     *slog.Logger
	templates  *templates.Templates
	mediaStore MediaStore
}

func NewDeps(
	service *Service,
	logger *slog.Logger,
	templateManager *templates.Templates,
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
