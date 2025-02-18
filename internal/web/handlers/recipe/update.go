package recipe

import (
	"fmt"
	"gluttony/internal/recipe"
	"gluttony/internal/web"
	"net/http"
	"strconv"
)

const (
	updateView = "views/recipe/update"
	updateForm = "recipe/form"
)

func (r *Routes) UpdateViewHandler(c *web.Context) error {
	recipeIDPathParam := c.Request.PathValue("recipe_id")
	recipeID, err := strconv.ParseInt(recipeIDPathParam, 10, 64)
	if err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	recipe, err := r.service.GetFull(c.Context(), recipeID)
	if err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	tags := make([]string, 0, len(recipe.Tags))
	for _, tag := range recipe.Tags {
		tags = append(tags, tag.Name)
	}

	c.Data["Form"] = Form{
		ID:                int64(recipe.ID),
		Name:              recipe.Name,
		Description:       recipe.Description,
		Source:            recipe.Source,
		Instructions:      recipe.InstructionsMarkdown,
		ThumbnailImageURL: recipe.ThumbnailImageURL,
		Servings:          recipe.Servings,
		PreparationTime:   recipe.PreparationTime,
		CookTime:          recipe.CookTime,
		Tags:              tags,
		Ingredients:       recipe.Ingredients,
		Nutrition:         recipe.Nutrition,
	}

	return c.RenderView(updateView, http.StatusOK)
}

func (r *Routes) UpdateFormHandler(c *web.Context) error {
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
			// TODO: handle err
			panic(fmt.Errorf("could not open cover image: %w", err))
		}
		defer file.Close()

		input.ThumbnailImage = file
	}

	err = r.service.Update(c.Context(), recipe.UpdateInput{
		ID:          form.ID,
		CreateInput: input,
	})
	if err == nil {
		c.HTMXRedirect(fmt.Sprintf("/recipes/%d", form.ID))
		return nil
	}

	c.Data["Form"] = form
	return c.RenderViewFragment(updateView, updateForm, http.StatusOK)
}
