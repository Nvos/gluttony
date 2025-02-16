package recipe

import (
	"gluttony/internal/ingredient"
	"gluttony/internal/recipe"
	"gluttony/internal/web"
	"net/http"
)

const (
	createView = "views/recipe/create"
	createForm = "recipe/form"
)

func (r *Routes) CreateViewHandler(c *web.Context) error {
	c.Data["Form"] = Form{
		Servings: 1,
		Ingredients: []recipe.Ingredient{
			{
				Ingredient: ingredient.Ingredient{
					Name: "",
				},
				Quantity: 0,
				Unit:     "g",
			},
		},
	}

	return c.RenderView(createView, http.StatusOK)
}

func (r *Routes) CreateFormHandler(c *web.Context) error {
	if err := c.FormParse(); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	form, err := NewRecipeForm(c.Request.MultipartForm.Value)
	if err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	input := form.ToInput()

	coverImage := c.Request.MultipartForm.File["thumbnail-image"]
	if len(coverImage) == 1 {
		file, err := coverImage[0].Open()
		if err != nil {
			return c.Error(http.StatusInternalServerError, err)
		}
		defer file.Close()

		input.ThumbnailImage = file
	}

	err = r.service.Create(c.Context(), input)
	if err == nil {
		c.HTMXRedirect("/recipes")
		return nil
	}

	c.Data["Form"] = form
	return c.RenderViewFragment(createView, createForm, http.StatusOK)
}
