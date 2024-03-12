package validateutil

// Rule returning bool indicates if there was violation
type Rule[T any] func(path string, value T) (bool, FieldViolation)
