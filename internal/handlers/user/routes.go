package user

import (
	"gluttony/internal/config"
	"gluttony/internal/user"
	"gluttony/x/httpx"
)

type Routes struct {
	cfg     *config.Config
	service *user.Service
}

func NewRoutes(
	cfg *config.Config,
	service *user.Service,
) (*Routes, error) {
	return &Routes{
		service: service,
		cfg:     cfg,
	}, nil
}

func (r *Routes) Mount(router *httpx.Router) {
	router.Get("/login", r.LoginViewHandler)
	router.Post("/login", r.LoginFormHandler)
	router.Get("/logout", r.LogoutHandler)
}
