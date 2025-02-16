package recipe

import (
	"fmt"
	"gluttony/internal/recipe"
	"gluttony/internal/web"
	"net/http"
	"strconv"
)

const (
	listView    = "views/recipe/list"
	listContent = "recipes/list"
)

func (r *Routes) ListViewHandler(c *web.Context) error {
	search := c.Request.URL.Query().Get("search")
	pageParam := c.Request.URL.Query().Get("page")
	limit := 20

	page := 0
	if pageParam != "" {
		pageInt, err := strconv.Atoi(pageParam)
		if err != nil {
			return fmt.Errorf("parse page to int: %w", err)
		}
		page = pageInt
	}

	recipePartials, err := r.service.AllSummaries(c.Context(), recipe.SearchInput{
		Search: search,
		Page:   int64(page),
		Limit:  int64(limit),
	})
	if err != nil {
		return fmt.Errorf("could not get recipe partials: %w", err)
	}

	c.Data["Recipes"] = recipePartials
	if c.IsHTMXRequest() {
		return c.RenderViewFragment(listView, listContent, http.StatusOK)
	}

	return c.RenderView(listView, http.StatusOK)
}
