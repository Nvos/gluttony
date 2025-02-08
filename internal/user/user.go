package user

import (
	"context"
	"gluttony/internal/security"
	"gluttony/internal/templating"
	"log/slog"
)

type Deps struct {
	service      *Service
	sessionStore SessionStore
	templates    *templating.Templating
	logger       *slog.Logger
}

func NewDeps(
	templates *templating.Templating,
	service *Service,
) *Deps {
	if service == nil {
		panic("nil service")
	}

	if templates == nil {
		panic("nil templates")
	}

	return &Deps{
		templates: templates,
		service:   service,
	}
}

type SessionStore interface {
	Get(ctx context.Context, key string) (security.Session, error)
	New(ctx context.Context) (security.Session, error)
	Save(ctx context.Context, value security.Session) error
	Delete(ctx context.Context, value security.Session) error
}
