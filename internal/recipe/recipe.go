package recipe

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
