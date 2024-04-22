package recipe

import (
	v1 "gluttony/internal/proto/recipe/v1"
	"gluttony/internal/x/validatex"
)

func ValidateCreateRecipeRequest(value *v1.CreateRecipeRequest) error {
	return validatex.NewValidationError(
		validatex.String("name", value.Name, validatex.Empty()),
	)
}
