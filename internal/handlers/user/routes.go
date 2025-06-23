package user

import (
	"gluttony/internal/service/user"
	"gluttony/pkg/router"
	"gluttony/pkg/session"
)

type Routes struct {
	service        *user.Service
	sessionService *session.Service
}

func NewRoutes(
	service *user.Service,
	sessionStore *session.Service,
) (*Routes, error) {
	return &Routes{
		service:        service,
		sessionService: sessionStore,
	}, nil
}

func (r *Routes) Mount(router *router.Router) {
	router.Get("/login", r.LoginViewHandler)
	router.Post("/login", r.LoginFormHandler)
	router.Get("/logout", r.LogoutHandler)
}
