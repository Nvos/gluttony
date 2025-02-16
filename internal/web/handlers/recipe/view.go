package recipe

import (
	"gluttony/internal/web"
	"net/http"
	"strconv"
)

const (
	detailsView = "views/recipe/view"
)

func (r *Routes) DetailsViewHandler(c *web.Context) error {
	recipeIDRaw := c.Request.PathValue("recipe_id")
	recipeID, err := strconv.Atoi(recipeIDRaw)
	if err != nil {
		return c.Error(http.StatusNotFound, nil)
	}

	recipe, err := r.service.GetFull(c.Context(), int64(recipeID))
	if err != nil {
		return c.Error(http.StatusNotFound, nil)
	}

	c.Data["Recipe"] = recipe

	return c.RenderView(detailsView, http.StatusOK)
}
