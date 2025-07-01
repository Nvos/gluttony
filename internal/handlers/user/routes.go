package user

import (
	"gluttony/internal/config"
	"gluttony/internal/service/user"
	"gluttony/pkg/router"
	"gluttony/pkg/session"
)

type Routes struct {
	cfg            *config.Config
	service        *user.Service
	sessionService *session.Service
}

func NewRoutes(
	cfg *config.Config,
	service *user.Service,
	sessionStore *session.Service,
) (*Routes, error) {
	return &Routes{
		service:        service,
		sessionService: sessionStore,
		cfg:            cfg,
	}, nil
}

func (r *Routes) Mount(router *router.Router) {
	router.Get("/login", r.LoginViewHandler)
	router.Post("/login", r.LoginFormHandler)
	router.Get("/logout", r.LogoutHandler)
}
