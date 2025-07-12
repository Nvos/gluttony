package handlers

import (
	"gluttony/internal/user"
	"gluttony/x/httpx"
	"io/fs"
	"net/http"
)

func GetDoer(c *httpx.Context) *user.User {
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

func AssetHandler(assetsFS fs.FS, isCache bool) httpx.HandlerFunc {
	httpFS := http.FileServerFS(assetsFS)

	return func(c *httpx.Context) error {
		if !isCache {
			c.Response.Header().Set("Cache-Control", "no-store")
		}

		http.StripPrefix("/assets/", httpFS).ServeHTTP(c.Response, c.Request)

		return nil
	}
}

func MediaHandler(mediaFS fs.FS) httpx.HandlerFunc {
	httpFS := http.FileServerFS(mediaFS)

	return func(c *httpx.Context) error {
		http.StripPrefix("/media", httpFS).ServeHTTP(c.Response, c.Request)

		return nil
	}
}
