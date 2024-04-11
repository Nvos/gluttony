package recipe

import (
	"fmt"
	v1 "gluttony/internal/proto/recipe/v1"
	"gluttony/internal/x/slicex"
	"gluttony/internal/x/validatex"
)

func ValidateCreateRecipeRequest(value *v1.CreateRecipeRequest) error {
	return validatex.NewValidationError(
		validatex.String("name", value.Name, validatex.Empty()),
		validatex.Array(value.Steps, func(index int, value *v1.CreateRecipeRequest_CreateRecipeStep) []validatex.FieldViolation {
			return slicex.Merge(
				validatex.Number(fmt.Sprintf("steps.%d.order", index), value.Order, validatex.Min[int32](0, true)),
				validatex.String(fmt.Sprintf("steps.%d.description", index), value.Description, validatex.Empty()),
			)
		}),
	)
}
