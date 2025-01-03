package recipe

import (
	"fmt"
	"gluttony/internal/share"
	"net/http"
)

func RecipeCreateViewHandler(deps *Deps) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		appCtx := share.MustGetContext(r.Context())
		get, err := deps.templates.Get("recipe", "recipe_create")
		if err != nil {
			// TODO: proper err
			panic(fmt.Errorf("could not get recipe template: %v", err))
		}

		err = get.View(w, appCtx)
		if err != nil {
			panic(fmt.Errorf("could not get recipe template: %v", err))
		}
	}
}

func RecipesViewHandler(deps *Deps) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		appCtx := share.MustGetContext(r.Context())
		get, err := deps.templates.Get("recipe", "recipes")
		if err != nil {
			// TODO: proper err
			panic(fmt.Errorf("could not get recipe template: %v", err))
		}

		err = get.View(w, appCtx)
		if err != nil {
			panic(fmt.Errorf("could not get recipe template: %v", err))
		}
	}
}

func RecipesCreateHandler(deps *Deps) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		println(fmt.Sprintf("%+v", r.Form["instructions"]))

		appCtx := share.MustGetContext(r.Context())

		get, err := deps.templates.Get("recipe", "recipe_create")
		if err != nil {
			// TODO: proper err
			panic(fmt.Errorf("could not get recipe template: %v", err))
		}

		err = get.Fragment(w, "recipe-create/form", appCtx)
		if err != nil {
			panic(fmt.Errorf("could not get recipe template: %v", err))
		}
	}
}
