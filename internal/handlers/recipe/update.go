package recipe

import (
	"errors"
	datastar "github.com/starfederation/datastar/sdk/go"
	"gluttony/internal/handlers"
	"gluttony/internal/recipe"
	"gluttony/web"
	"gluttony/web/component"
	"gluttony/x/httpx"
	"net/http"
	"strconv"
)

const (
	updateView = "view/recipe/update"
	updateForm = "recipe/form"
)

func (r *Routes) UpdateViewHandler(c *httpx.Context) error {
	recipeIDPathParam := c.Request.PathValue("recipe_id")
	recipeID, err := strconv.ParseInt(recipeIDPathParam, 10, 32)
	if err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	fullRecipe, err := r.service.GetFull(c.Context(), int32(recipeID))
	if err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	tags := make([]string, 0, len(fullRecipe.Tags))
	for _, tag := range fullRecipe.Tags {
		tags = append(tags, tag.Name)
	}

	form := recipe.Form{
		ID:              fullRecipe.ID,
		Name:            fullRecipe.Name,
		Description:     fullRecipe.Description,
		Source:          fullRecipe.Source,
		Instructions:    fullRecipe.InstructionsMarkdown,
		ThumbnailImage:  nil,
		Servings:        fullRecipe.Servings,
		PreparationTime: fullRecipe.PreparationTime,
		CookTime:        fullRecipe.CookTime,
		Tags:            tags,
		Ingredients:     fullRecipe.Ingredients,
		Nutrition:       fullRecipe.Nutrition,
	}

	webCtx := web.NewContext(c.Request, handlers.GetDoer(c), "en")
	return c.TemplComponent(
		http.StatusOK,
		component.ViewRecipeUpdate(webCtx, fullRecipe.ThumbnailImageURL, form),
	)
}

func (r *Routes) UpdateFormHandler(c *httpx.Context) error {
	if err := c.FormParse(); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	form, err := recipe.NewRecipeForm(c.Request.MultipartForm)
	if err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	input := form.ToInput(handlers.GetDoer(c).ID)

	sse := datastar.NewSSE(c.Response, c.Request)
	err = r.service.Update(c.Context(), recipe.UpdateInput{
		ID:          form.ID,
		CreateInput: input,
	})
	if err == nil {
		if err := sse.Redirect("/recipes"); err != nil {
			return c.Error(http.StatusInternalServerError, err)
		}

		return nil
	}

	uniqueName := ""
	if errors.Is(err, recipe.ErrUniqueName) {
		uniqueName = "Recipe with such name exists"
	}
	err = sse.MarshalAndMergeSignals(map[string]string{
		"errors.name": uniqueName,
	})
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return nil
}
