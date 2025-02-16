package handlers

import (
	"gluttony/internal/html"
	"gluttony/internal/security"
	"gluttony/internal/user"
	"gluttony/internal/web"
)

type Routes struct {
	service      *user.Service
	sessionStore *security.SessionStore
	renderer     *html.Renderer
}

func NewRoutes(
	service *user.Service,
	sessionStore *security.SessionStore,
) (*Routes, error) {
	return &Routes{
		service:      service,
		sessionStore: sessionStore,
	}, nil
}

func (r *Routes) Mount(router *web.Router) {
	router.Get("/login", r.LoginViewHandler)
	router.Post("/login/form", r.LoginHTMXFormHandler)
}
