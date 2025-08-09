package recipe

import (
	"errors"
	datastar "github.com/starfederation/datastar/sdk/go"
	"gluttony/ingredient"
	"gluttony/recipe"
	"gluttony/web"
	"gluttony/web/component"
	"gluttony/web/handlers"
	"gluttony/x/httpx"
	"net/http"
)

func (r *Routes) CreateViewHandler(c *httpx.Context) error {
	form := recipe.Form{
		ID:              0,
		Name:            "",
		Description:     "",
		Source:          "",
		Instructions:    "",
		ThumbnailImage:  nil,
		Servings:        1,
		PreparationTime: 0,
		CookTime:        0,
		Tags:            []string{},
		Ingredients: []recipe.Ingredient{
			{
				Order:    0,
				Quantity: 0,
				Note:     "",
				Unit:     "g",
				Ingredient: ingredient.Ingredient{
					ID:   0,
					Name: "",
				},
			},
		},
		Nutrition: recipe.Nutrition{
			Calories: 0,
			Fat:      0,
			Carbs:    0,
			Protein:  0,
		},
	}

	webCtx := web.NewContext(c.Request, handlers.GetDoer(c), "en")
	return c.TemplComponent(
		http.StatusOK,
		component.ViewRecipeCreate(webCtx, form),
	)
}

func (r *Routes) CreateFormHandler(c *httpx.Context) error {
	if err := c.FormParse(); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	form, err := recipe.NewRecipeForm(c.Request.MultipartForm)
	if err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	input := form.ToInput(handlers.GetDoer(c).ID)

	sse := datastar.NewSSE(c.Response, c.Request)

	err = r.service.Create(c.Context(), input)
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
