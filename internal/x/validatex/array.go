package validatex

func Array[T any](arr []T, fn func(index int, value T) []FieldViolation) []FieldViolation {
	var violations []FieldViolation
	for i := range arr {
		if result := fn(i, arr[i]); len(result) > 0 {
			violations = append(violations, result...)
		}
	}

	return violations
}
