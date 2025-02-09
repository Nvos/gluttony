package app

import (
	"github.com/spf13/afero"
	"gluttony/internal/config"
	"gluttony/internal/httputil"
	"gluttony/internal/livereload"
	"gluttony/internal/recipe"
	"gluttony/internal/security"
	"gluttony/internal/share"
	"gluttony/internal/templating"
	"gluttony/internal/user"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

func MountWebRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	sessionStore *security.SessionStore,
	userService *user.Service,
	recipeService *recipe.Service,
) {
	userTemplating := templating.New(
		os.DirFS(filepath.Join("internal/user/templates")),
	)
	recipeTemplating := templating.New(
		os.DirFS(filepath.Join("internal/recipe/templates")),
	)

	middlewares := []httputil.MiddlewareFunc{
		security.NewAuthenticationMiddleware(sessionStore),
		share.ContextMiddleware,
		httputil.NewErrorMiddleware(logger),
	}

	userDeps := user.NewDeps(userTemplating, userService)
	user.Routes(userDeps, mux, middlewares...)

	recipeDeps := recipe.NewDeps(recipeService, recipeTemplating)
	recipe.Routes(recipeDeps, mux, middlewares...)
}

func MountRoutes(
	mux *http.ServeMux,
	mode config.Mode,
	liveReload *livereload.LiveReload,
	directories *Directories,
) {
	if mode == config.Dev {
		mux.HandleFunc("GET /reload", liveReload.Handle)
	}

	mux.HandleFunc("GET /assets/{pathname...}", handleAssets(mode, directories))
	mux.HandleFunc("GET /media/{pathname...}", handleMedia(directories))
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
