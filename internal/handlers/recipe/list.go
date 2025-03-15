package recipe

import (
	"fmt"
	"gluttony/internal/recipe"
	"gluttony/pkg/router"
	"net/http"
	"strconv"
)

const (
	listView    = "view/recipe/list"
	listContent = "recipes/list"
)

func (r *Routes) ListViewHandler(c *router.Context) error {
	search := c.Request.URL.Query().Get("search")
	pageParam := c.Request.URL.Query().Get("page")

	page := int32(0)
	if pageParam != "" {
		pageInt, err := strconv.ParseInt(pageParam, 10, 32)
		if err != nil {
			return fmt.Errorf("parse page to int: %w", err)
		}

		page = int32(pageInt)
	}

	recipePartials, err := r.service.AllSummaries(c.Context(), recipe.SearchInput{
		RecipeIDs: nil,
		Search:    search,
		Page:      page,
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
