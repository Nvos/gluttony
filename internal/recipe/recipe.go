package recipe

type Recipe struct {
	ID          int
	Name        string
	Description string
}

type Step struct {
	ID          int
	Order       int
	Name        string
	Description string
}
