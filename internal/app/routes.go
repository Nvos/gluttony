package app

import (
	"github.com/spf13/afero"
	recipehandlers "gluttony/internal/handlers/recipe"
	userhandlers "gluttony/internal/handlers/user"
	"gluttony/internal/service/recipe"
	"gluttony/internal/service/user"
	"gluttony/pkg/livereload"
	"gluttony/pkg/router"
	"gluttony/pkg/session"
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
	mode Mode,
	liveReload *livereload.LiveReload,
	directories *Directories,
) {
	if mode == Dev {
		mux.Get("/reload", router.WrapHandlerFunc(liveReload.Handle))
	}

	mux.Get("/assets/{pathname...}", router.WrapHandlerFunc(handleAssets(mode, directories)))
	mux.Get("/media/{pathname...}", router.WrapHandlerFunc(handleMedia(directories)))
}

func handleAssets(mode Mode, directories *Directories) func(w http.ResponseWriter, r *http.Request) {
	httpFS := http.FileServerFS(directories.Assets)
	return func(w http.ResponseWriter, r *http.Request) {
		if mode == Dev {
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
