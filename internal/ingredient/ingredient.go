package ingredient

import (
	"errors"
	"fmt"
	"gluttony/internal/database/pagination"
	"gluttony/internal/i18n"
	ingredientv1 "gluttony/internal/proto/ingredient/v1"
)

type Unit string

func (u Unit) Connect() (ingredientv1.Unit, error) {
	switch u {
	case Weight:
		return ingredientv1.Unit_UNIT_WEIGHT, nil
	case Volume:
		return ingredientv1.Unit_UNIT_VOLUME, nil
	}

	return 0, fmt.Errorf("invalid unit=%s", u)
}

const (
	Weight = "weight"
	Volume = "volume"
)

func NewUnit(value string) (Unit, error) {
	switch value {
	case Weight:
		return Weight, nil
	case Volume:
		return Volume, nil
	}

	return "", fmt.Errorf("unknown ingredient unit=%s", value)
}

func NewUnitFromConnect(value ingredientv1.Unit) (Unit, error) {
	switch value {
	case ingredientv1.Unit_UNIT_VOLUME:
		return Volume, nil
	case ingredientv1.Unit_UNIT_WEIGHT:
		return Weight, nil
	}

	return "", fmt.Errorf("unknown ingredient unit=%d", value)
}

type CreateIngredientInput struct {
	Locale i18n.Locale
	Name   string
	Unit   Unit
}

func NewCreateIngredientInput(r *ingredientv1.CreateRequest) (CreateIngredientInput, error) {
	if r == nil {
		return CreateIngredientInput{}, errors.New("nil create request")
	}

	locale, err := i18n.NewLocale(r.Locale)
	if err != nil {
		return CreateIngredientInput{}, err
	}

	unit, err := NewUnitFromConnect(r.Unit)
	if err != nil {
		return CreateIngredientInput{}, err
	}

	return CreateIngredientInput{
		Locale: locale,
		Name:   r.Name,
		Unit:   unit,
	}, nil
}

type AllIngredientsInput struct {
	Locale     i18n.Locale
	Search     string
	Pagination pagination.OffsetPagination
}

type SingleInput struct {
	ID     int32
	Locale i18n.Locale
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
	Unit Unit
}
