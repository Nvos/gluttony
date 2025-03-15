package recipe

import (
	"gluttony/pkg/router"
	"net/http"
	"strconv"
)

const (
	detailsView = "view/recipe/view"
)

func (r *Routes) DetailsViewHandler(c *router.Context) error {
	recipeIDRaw := c.Request.PathValue("recipe_id")
	recipeID, err := strconv.ParseInt(recipeIDRaw, 10, 32)
	if err != nil {
		return c.Error(http.StatusNotFound, nil)
	}

	recipe, err := r.service.GetFull(c.Context(), int32(recipeID))
	if err != nil {
		return c.Error(http.StatusNotFound, nil)
	}

	c.Data["Recipe"] = recipe

	return c.RenderView(detailsView, http.StatusOK)
}
