package recipe

import (
	"fmt"
	"gluttony/internal/database/pagination"
	"gluttony/internal/i18n"
	"gluttony/internal/ingredient"
	v1 "gluttony/internal/proto/recipe/v1"
)

type Recipe struct {
	ID          int32
	Name        string
	Description string
}

type Ingredient struct {
	ingredient.Ingredient
	Amount int32
	Count  int32
	Note   string
}

type FullRecipe struct {
	Recipe
	Content     string
	Ingredients []Ingredient
}

type CreateRecipe struct {
	Locale      i18n.Locale
	Name        string
	Description string
	Content     string
	Ingredients []CreateIngredient
}

type CreateIngredient struct {
	ID     int32
	Amount int32
	Count  int32
	Note   string
}

type AllRecipesInput struct {
	Locale     i18n.Locale
	Search     string
	Pagination pagination.OffsetPagination
}

func NewAllRecipesInput(locale i18n.Locale, r *v1.AllRecipesRequest) (AllRecipesInput, error) {
	if r == nil {
		return AllRecipesInput{}, fmt.Errorf("all recipes request is nil")
	}

	offset, err := pagination.NewOffsetPagination(r.Offset, r.Limit)
	if err != nil {
		return AllRecipesInput{}, fmt.Errorf("new all recipe input: %w", err)
	}

	return AllRecipesInput{
		Locale:     locale,
		Search:     r.Search,
		Pagination: offset,
	}, nil
}

func NewCreateRecipe(r *v1.CreateRecipeRequest) (CreateRecipe, error) {
	if err := ValidateCreateRecipeRequest(r); err != nil {
		return CreateRecipe{}, err
	}

	locale, err := i18n.NewLocale(r.Locale)
	if err != nil {
		return CreateRecipe{}, fmt.Errorf("new create recipe locale: %w", err)
	}

	ingredients := make([]CreateIngredient, 0, len(r.Ingredients))
	for i := range r.Ingredients {
		row := r.Ingredients[i]
		ingredients = append(ingredients, CreateIngredient{
			ID:     row.Id,
			Amount: row.Amount,
			Count:  row.Count,
			Note:   row.Note,
		})
	}

	return CreateRecipe{
		Name:        r.Name,
		Description: r.Description,
		Content:     r.Content,
		Locale:      locale,
		Ingredients: ingredients,
	}, nil
}
