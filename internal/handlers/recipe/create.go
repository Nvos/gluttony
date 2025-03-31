package recipe

import (
	"gluttony/internal/handlers"
	"gluttony/internal/ingredient"
	"gluttony/internal/recipe"
	"gluttony/pkg/router"
	"mime/multipart"
	"net/http"
)

const (
	createView = "view/recipe/create"
	createForm = "recipe/form"
)

func (r *Routes) CreateViewHandler(c *router.Context) error {
	c.Data["Form"] = Form{
		ID:                0,
		Name:              "",
		Description:       "",
		Source:            "",
		Instructions:      "",
		ThumbnailImageURL: "",
		Servings:          1,
		PreparationTime:   0,
		CookTime:          0,
		Tags:              []string{},
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

	return c.RenderView(createView, http.StatusOK)
}

func (r *Routes) CreateFormHandler(c *router.Context) error {
	if err := c.FormParse(); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	form, err := NewRecipeForm(c.Request.MultipartForm.Value)
	if err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	input := form.ToInput(handlers.GetDoer(c).ID)

	coverImage := c.Request.MultipartForm.File["thumbnail-image"]
	if len(coverImage) == 1 {
		file, err := coverImage[0].Open()
		if err != nil {
			return c.Error(http.StatusInternalServerError, err)
		}
		defer func(file multipart.File) {
			_ = file.Close()
		}(file)

		input.ThumbnailImage = file
	}

	err = r.service.Create(c.Context(), input)
	if err == nil {
		c.HTMXRedirect("/recipes")
		return nil
	}

	return c.Error(http.StatusBadRequest, err)
}
