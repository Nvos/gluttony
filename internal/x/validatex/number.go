package validatex

import "golang.org/x/exp/constraints"

type Numeric interface {
	constraints.Integer | constraints.Float
}

func Number[T Numeric](path string, value T, rules ...Rule[T]) []FieldViolation {
	var violations []FieldViolation

	for i := range rules {
		if ok, violation := rules[i](path, value); ok {
			violations = append(violations, violation)
		}
	}

	return violations
}

func Min[T Numeric](min T, inclusive bool) Rule[T] {
	return func(path string, value T) (bool, FieldViolation) {
		if inclusive && value <= min {
			return true, FieldViolation{
				Path: path,
				Rule: "min",
			}
		}

		if value < min {
			return true, FieldViolation{
				Path: path,
				Rule: "min",
			}
		}

		return false, FieldViolation{}
	}
}

func Max[T constraints.Integer | constraints.Float](max T) Rule[T] {
	return func(path string, value T) (bool, FieldViolation) {
		if value > max {
			return true, FieldViolation{
				Path: path,
				Rule: "max",
			}
		}

		return false, FieldViolation{}
	}
}
