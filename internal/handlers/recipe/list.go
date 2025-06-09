package recipe

import (
	"errors"
	"fmt"
	"gluttony/internal/handlers"
	"gluttony/internal/recipe"
	"gluttony/pkg/pagination"
	"gluttony/pkg/router"
	"gluttony/web"
	"gluttony/web/component"
	"net/http"
	"strconv"
)

type RecipesQueryParams struct {
	Query string
	Page  int32
}

func (r *Routes) RecipesHandler(c *router.Context) error {
	params, err := readRecipesURLParams(c.Request)
	if err != nil {
		return router.NewHTTPError(
			http.StatusBadRequest,
			router.WithError(err),
		)
	}

	summariesPage, err := r.service.AllSummaries(c.Context(), recipe.SearchInput{
		RecipeIDs: nil,
		Search:    params.Query,
		Page:      params.Page,
	})
	if err != nil {
		return fmt.Errorf("could not get recipe partials: %w", err)
	}

	paginator := pagination.New(params.Page, summariesPage.TotalCount)
	if params.Page > paginator.TotalCount {
		return router.NewHTTPError(http.StatusNotFound)
	}

	webCtx := web.NewContext(c.Request, handlers.GetDoer(c), "en")
	return c.TemplComponent(
		http.StatusOK,
		component.ViewRecipes(webCtx, params.Query, paginator, summariesPage.Rows),
	)
}

func readRecipesURLParams(r *http.Request) (RecipesQueryParams, error) {
	query := r.URL.Query().Get("query")
	pageParam := r.URL.Query().Get("page")

	page := int32(0)
	if pageParam != "" {
		pageInt, err := strconv.ParseInt(pageParam, 10, 32)
		if err != nil {
			return RecipesQueryParams{}, fmt.Errorf("parse page to int: %w", err)
		}

		page = int32(pageInt)
		if page < 0 {
			return RecipesQueryParams{}, errors.New("page must be >= 0")
		}
	}

	return RecipesQueryParams{
		Query: query,
		Page:  page,
	}, nil
}
