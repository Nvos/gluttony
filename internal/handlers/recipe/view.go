package recipe

import (
	"gluttony/internal/handlers"
	"gluttony/pkg/router"
	"gluttony/web"
	"gluttony/web/component"
	"net/http"
	"strconv"
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

	webCtx := web.NewContext(c.Request, handlers.GetDoer(c), "en")

	return c.TemplComponent(http.StatusOK, component.ViewRecipe(webCtx, recipe))
}
