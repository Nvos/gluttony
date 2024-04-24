package ingredient

import (
	"errors"
	"gluttony/internal/database/pagination"
	"gluttony/internal/i18n"
	ingredientv1 "gluttony/internal/proto/ingredient/v1"
)

type CreateIngredientInput struct {
	Locale i18n.Locale
	Name   string
}

func NewCreateIngredientInput(r *ingredientv1.CreateRequest) (CreateIngredientInput, error) {
	if r == nil {
		return CreateIngredientInput{}, errors.New("nil create request")
	}

	locale, err := i18n.NewLocale(r.Locale)
	if err != nil {
		return CreateIngredientInput{}, err
	}

	return CreateIngredientInput{
		Locale: locale,
		Name:   r.Name,
	}, nil
}

type AllIngredientsInput struct {
	Locale     i18n.Locale
	Search     string
	Pagination pagination.OffsetPagination
}

func NewAllIngredientsInput(locale i18n.Locale, r *ingredientv1.AllRequest) (AllIngredientsInput, error) {
	if r == nil {
		return AllIngredientsInput{}, errors.New("nil all request")
	}

	offset, err := pagination.NewOffsetPagination(r.Offset, r.Limit)
	if err != nil {
		return AllIngredientsInput{}, err
	}

	return AllIngredientsInput{
		Locale:     locale,
		Search:     r.Search,
		Pagination: offset,
	}, nil
}

type Ingredient struct {
	ID   int32
	Name string
}
