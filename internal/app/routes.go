package app

import (
	"github.com/spf13/afero"
	"gluttony/internal/config"
	"gluttony/internal/livereload"
	"gluttony/internal/recipe"
	"gluttony/internal/security"
	"gluttony/internal/user"
	"gluttony/internal/web"
	"gluttony/internal/web/handlers"
	recipehandlers "gluttony/internal/web/handlers/recipe"
	"log/slog"
	"net/http"
)

func MountWebRoutes(
	mux *web.Router,
	logger *slog.Logger,
	sessionStore *security.SessionStore,
	userService *user.Service,
	recipeService *recipe.Service,
) {
	recipeRoutes, err := recipehandlers.NewRoutes(recipeService)
	if err != nil {
		panic(err)
	}

	userRouter, err := handlers.NewRoutes(userService, sessionStore)
	if err != nil {
		panic(err)
	}

	userRouter.Mount(mux)
	recipeRoutes.Mount(mux)
}

func MountRoutes(
	mux *web.Router,
	mode config.Mode,
	liveReload *livereload.LiveReload,
	directories *Directories,
) {
	if mode == config.Dev {
		mux.Get("/reload", web.WrapHandlerFunc(liveReload.Handle))
	}

	mux.Get("/assets/{pathname...}", web.WrapHandlerFunc(handleAssets(mode, directories)))
	mux.Get("/media/{pathname...}", web.WrapHandlerFunc(handleMedia(directories)))
}

func handleAssets(mode config.Mode, directories *Directories) func(w http.ResponseWriter, r *http.Request) {
	httpFS := http.FileServerFS(directories.Assets)
	return func(w http.ResponseWriter, r *http.Request) {
		if mode == config.Dev {
			w.Header().Set("Cache-Control", "no-store")
		}

		http.StripPrefix("/assets", httpFS).ServeHTTP(w, r)
	}
}

func handleMedia(directories *Directories) func(w http.ResponseWriter, r *http.Request) {
	httpFS := http.FileServerFS(afero.NewIOFS(directories.Media))
	return func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/media", httpFS).ServeHTTP(w, r)
	}
}
