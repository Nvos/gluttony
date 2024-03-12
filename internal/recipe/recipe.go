package recipe

import v1 "gluttony/internal/proto/recipe/v1"

type Recipe struct {
	ID          int32
	Name        string
	Description string

	Steps []Step
}

type CreateRecipe struct {
	Name        string
	Description string
	Steps       []CreateStep
}

type CreateStep struct {
	Order       int32
	Description string
}

type Step struct {
	ID          int32
	Order       int32
	Description string
}

func NewCreateRecipe(r *v1.CreateRecipeRequest) (CreateRecipe, error) {
	if err := ValidateCreateRecipeRequest(r); err != nil {
		return CreateRecipe{}, err
	}

	steps := make([]CreateStep, 0, len(r.Steps))
	for i := range r.Steps {
		v := r.Steps[i]
		step := CreateStep{
			Order:       v.Order,
			Description: v.Description,
		}

		steps = append(steps, step)
	}

	return CreateRecipe{
		Name:        r.Name,
		Description: r.Description,
		Steps:       steps,
	}, nil
}
