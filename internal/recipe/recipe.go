package recipe

import (
	"gluttony/internal/templates"
	"log/slog"
	"time"
)

type Deps struct {
	service   *Service
	logger    *slog.Logger
	templates *templates.Templates
}

func NewDeps(service *Service, logger *slog.Logger, templateManager *templates.Templates) *Deps {
	if logger == nil {
		panic("logger must not be nil")
	}

	if service == nil {
		panic("service must not be nil")
	}

	if templateManager == nil {
		panic("templateManager must not be nil")
	}

	return &Deps{
		service:   service,
		logger:    logger,
		templates: templateManager,
	}
}

type CreateRecipeInput struct {
	Name      string
	Steps     string
	CreatedAt time.Time
}
