package handlers

import (
	"gluttony/internal/user"
	"gluttony/pkg/router"
	"io/fs"
	"net/http"
)

func GetDoer(c *router.Context) *user.User {
	value, ok := c.Data["User"]
	if !ok {
		return nil
	}

	doer, ok := value.(user.User)
	if !ok {
		panic("Invalid type set in 'User' router ctx key")
	}

	return &doer
}

func AssetHandler(assetsFS fs.FS, isCache bool) router.HandlerFunc {
	httpFS := http.FileServerFS(assetsFS)

	return func(c *router.Context) error {
		if !isCache {
			c.Response.Header().Set("Cache-Control", "no-store")
		}

		http.StripPrefix("/assets/", httpFS).ServeHTTP(c.Response, c.Request)

		return nil
	}
}

func MediaHandler(mediaFS fs.FS) router.HandlerFunc {
	httpFS := http.FileServerFS(mediaFS)

	return func(c *router.Context) error {
		http.StripPrefix("/media", httpFS).ServeHTTP(c.Response, c.Request)

		return nil
	}
}
