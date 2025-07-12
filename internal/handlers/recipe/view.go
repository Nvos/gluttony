package recipe

import (
	"gluttony/internal/handlers"
	"gluttony/web"
	"gluttony/web/component"
	"gluttony/x/httpx"
	"net/http"
	"strconv"
)

func (r *Routes) DetailsViewHandler(c *httpx.Context) error {
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
