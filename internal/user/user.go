package user

import (
	"context"
	"gluttony/internal/security"
	"gluttony/internal/templates"
	"log/slog"
)

type Deps struct {
	service      *Service
	sessionStore SessionStore
	templates    *templates.Templates
	logger       *slog.Logger
}

func NewDeps(
	sessionStore SessionStore,
	templates *templates.Templates,
	logger *slog.Logger,
	service *Service,
) *Deps {
	if service == nil {
		panic("nil service")
	}

	if logger == nil {
		panic("nil logger")
	}

	if sessionStore == nil {
		panic("nil sessionStore")
	}

	if templates == nil {
		panic("nil templates")
	}

	return &Deps{
		sessionStore: sessionStore,
		templates:    templates,
		logger:       logger,
		service:      service,
	}
}

type SessionStore interface {
	Get(ctx context.Context, key string) (security.Session, error)
	New(ctx context.Context) (security.Session, error)
	Save(ctx context.Context, value security.Session) error
	Delete(ctx context.Context, value security.Session) error
}
