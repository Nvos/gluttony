package app

import (
	"gluttony/internal/handlers"
	recipehandlers "gluttony/internal/handlers/recipe"
	userhandlers "gluttony/internal/handlers/user"
	"gluttony/internal/service/recipe"
	"gluttony/internal/service/user"
	"gluttony/pkg/livereload"
	"gluttony/pkg/router"
	"gluttony/pkg/session"
	"io/fs"
	"net/http"
)

func MountWebRoutes(
	mux *router.Router,
	sessionService *session.Service,
	userService *user.Service,
	recipeService *recipe.Service,
) {
	recipeRoutes, err := recipehandlers.NewRoutes(recipeService)
	if err != nil {
		panic(err)
	}

	userRouter, err := userhandlers.NewRoutes(userService, sessionService)
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
	mode Environment,
	liveReload *livereload.LiveReload,
	assetsFS fs.FS,
	mediaFS fs.FS,
) {
	if mode == EnvDevelopment {
		mux.Get("/reload", router.WrapHandlerFunc(liveReload.Handle))
	}

	isCacheEnabled := mode == EnvProduction
	mux.Get("/assets/{pathname...}", handlers.AssetHandler(assetsFS, isCacheEnabled))
	mux.Get("/media/{pathname...}", handlers.MediaHandler(mediaFS))
}
