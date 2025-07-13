package app

import (
	"gluttony/internal/config"
	"gluttony/internal/handlers"
	recipehandlers "gluttony/internal/handlers/recipe"
	userhandlers "gluttony/internal/handlers/user"
	"gluttony/internal/recipe"
	"gluttony/internal/user"
	"gluttony/x/httpx"
	"io/fs"
	"net/http"
)

func MountWebRoutes(
	mux *httpx.Router,
	cfg *config.Config,
	userService *user.Service,
	recipeService *recipe.Service,
) {
	recipeRoutes, err := recipehandlers.NewRoutes(recipeService)
	if err != nil {
		panic(err)
	}

	userRouter, err := userhandlers.NewRoutes(cfg, userService)
	if err != nil {
		panic(err)
	}

	mux.Get("/", func(c *httpx.Context) error {
		c.Redirect("/recipes", http.StatusFound)

		return nil
	})

	userRouter.Mount(mux)
	recipeRoutes.Mount(mux)
}

func MountRoutes(
	mux *httpx.Router,
	mode config.Mode,
	assetsFS fs.FS,
	mediaFS fs.FS,
) {
	isCacheEnabled := mode == config.ModeProd
	mux.Get("/assets/{pathname...}", handlers.AssetHandler(assetsFS, isCacheEnabled))
	mux.Get("/media/{pathname...}", handlers.MediaHandler(mediaFS))
}
