package app

import (
	"gluttony/internal/config"
	"gluttony/internal/handlers"
	recipehandlers "gluttony/internal/handlers/recipe"
	userhandlers "gluttony/internal/handlers/user"
	"gluttony/internal/service/recipe"
	"gluttony/internal/service/user"
	"gluttony/pkg/router"
	"gluttony/pkg/session"
	"io/fs"
	"net/http"
)

func MountWebRoutes(
	mux *router.Router,
	cfg *config.Config,
	sessionService *session.Service,
	userService *user.Service,
	recipeService *recipe.Service,
) {
	recipeRoutes, err := recipehandlers.NewRoutes(recipeService)
	if err != nil {
		panic(err)
	}

	userRouter, err := userhandlers.NewRoutes(cfg, userService, sessionService)
	if err != nil {
		panic(err)
	}

	mux.Get("/", func(c *router.Context) error {
		c.Redirect("/recipes", http.StatusFound)

		return nil
	})

	userRouter.Mount(mux)
	recipeRoutes.Mount(mux)
}

func MountRoutes(
	mux *router.Router,
	mode config.Environment,
	assetsFS fs.FS,
	mediaFS fs.FS,
) {
	isCacheEnabled := mode == config.EnvProduction
	mux.Get("/assets/{pathname...}", handlers.AssetHandler(assetsFS, isCacheEnabled))
	mux.Get("/media/{pathname...}", handlers.MediaHandler(mediaFS))
}
