package recipe

import (
	"fmt"
	"gluttony/internal/handlers"
	"gluttony/internal/recipe"
	"gluttony/pkg/router"
	"gluttony/web"
	"gluttony/web/component"
	"mime/multipart"
	"net/http"
	"strconv"
)

const (
	updateView = "view/recipe/update"
	updateForm = "recipe/form"
)

func (r *Routes) UpdateViewHandler(c *router.Context) error {
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
		ID:                fullRecipe.ID,
		Name:              fullRecipe.Name,
		Description:       fullRecipe.Description,
		Source:            fullRecipe.Source,
		Instructions:      fullRecipe.InstructionsMarkdown,
		ThumbnailImageURL: fullRecipe.ThumbnailImageURL,
		Servings:          fullRecipe.Servings,
		PreparationTime:   fullRecipe.PreparationTime,
		CookTime:          fullRecipe.CookTime,
		Tags:              tags,
		Ingredients:       fullRecipe.Ingredients,
		Nutrition:         fullRecipe.Nutrition,
	}

	webCtx := web.NewContext(c.Request, handlers.GetDoer(c), "en")
	return c.TemplComponent(
		http.StatusOK,
		component.ViewRecipeCreate(webCtx, form),
	)
}

func (r *Routes) UpdateFormHandler(c *router.Context) error {
	if err := c.FormParse(); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	form, err := recipe.NewRecipeForm(c.Request.MultipartForm.Value)
	if err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	input := form.ToInput(handlers.GetDoer(c).ID)

	coverImage := c.Request.MultipartForm.File["thumbnail-image"]
	if len(coverImage) == 1 {
		file, err := coverImage[0].Open()
		if err != nil {
			// TODO: handle err
			panic(fmt.Errorf("could not open cover image: %w", err))
		}
		defer func(file multipart.File) {
			_ = file.Close()
		}(file)

		input.ThumbnailImage = file
	}

	err = r.service.Update(c.Context(), recipe.UpdateInput{
		ID:          form.ID,
		CreateInput: input,
	})
	if err == nil {
		c.Redirect(fmt.Sprintf("/recipes/%d", form.ID), http.StatusFound)
		return nil
	}

	// TODO: Handle errors
	return c.Error(http.StatusBadRequest, err)
}
