package ingredient

import (
	"gluttony/internal/database/pagination"
	"gluttony/internal/i18n"
)

type CreateIngredientInput struct {
	Locale i18n.Locale
	Name   string
}

type AllIngredientsInput struct {
	Locale     i18n.Locale
	Search     string
	Pagination pagination.OffsetPagination
}

type Ingredient struct {
	ID   int32
	Name string
}
