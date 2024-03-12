package recipe

import (
	"fmt"
	v1 "gluttony/internal/proto/recipe/v1"
	"gluttony/internal/util/sliceutil"
	"gluttony/internal/util/validateutil"
)

func ValidateCreateRecipeRequest(value *v1.CreateRecipeRequest) error {
	return validateutil.NewValidationError(
		validateutil.String("name", value.Name, validateutil.Empty()),
		validateutil.Array(value.Steps, func(index int, value *v1.CreateRecipeStep) []validateutil.FieldViolation {
			return sliceutil.Merge(
				validateutil.Number(fmt.Sprintf("steps.%d.order", index), value.Order, validateutil.Min[int32](0, true)),
				validateutil.String(fmt.Sprintf("steps.%d.description", index), value.Description, validateutil.Empty()),
			)
		}),
	)
}
