package validateutil

func String(path string, value string, rules ...Rule[string]) []FieldViolation {
	var violations []FieldViolation

	for i := range rules {
		if ok, violation := rules[i](path, value); ok {
			violations = append(violations, violation)
		}
	}

	return violations
}

func Empty() Rule[string] {
	return func(path string, value string) (bool, FieldViolation) {
		if value == "" {
			return true, FieldViolation{
				Path: path,
				Rule: "empty",
			}
		}

		return false, FieldViolation{}
	}
}

func MinLength(min int) Rule[string] {
	return func(path string, value string) (bool, FieldViolation) {
		if len(value) < min {
			return true, FieldViolation{
				Path: path,
				Rule: "min-length",
			}
		}

		return false, FieldViolation{}
	}
}

func MaxLength(max int) Rule[string] {
	return func(path string, value string) (bool, FieldViolation) {
		if len(value) > max {
			return true, FieldViolation{
				Path: path,
				Rule: "max-length",
			}
		}

		return false, FieldViolation{}
	}
}
