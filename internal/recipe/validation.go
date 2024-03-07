package recipe

import (
	"fmt"
)

func ValidateRecipe(recipe Recipe) error {
	if recipe.ID <= 0 {
		return fmt.Errorf("validate recipe: id <= 0")
	}

	if len(recipe.Name) == 0 {
		return fmt.Errorf("validate recipe: empty name")
	}

	return nil
}

func ValidateCreateRecipe(value CreateRecipe) error {
	if len(value.Name) == 0 {
		return fmt.Errorf("validate recipe: empty name")
	}

	return nil
}

func ValidateRecipeStep(step Step) error {
	if step.ID <= 0 {
		return fmt.Errorf("validate recipe step: id <= 0")
	}

	if step.Order <= 0 {
		return fmt.Errorf("validate recipe step: order <= 0")
	}

	if len(step.Description) == 0 {
		return fmt.Errorf("validate recipe step: empty description")
	}

	return nil
}

func ValidateCreateRecipeStep(value CreateStep) error {
	if value.Order <= 0 {
		return fmt.Errorf("validate recipe step: order <= 0")
	}

	if len(value.Description) == 0 {
		return fmt.Errorf("validate recipe step: empty description")
	}

	return nil
}
