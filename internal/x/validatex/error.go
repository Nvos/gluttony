package validatex

import (
	"fmt"
	"strings"
)

type FieldViolation struct {
	Path string
	Rule string
}

type ValidationError struct {
	Violations []FieldViolation
}

func (v *ValidationError) Error() string {
	parts := make([]string, 0, len(v.Violations))
	for i := range v.Violations {
		parts = append(parts, fmt.Sprintf("{ path: %s, rule: %s }", v.Violations[i].Path, v.Violations[i].Rule))
	}

	return fmt.Sprintf("validation error, violations=[%s]", strings.Join(parts, ", "))
}

func NewValidationError(violations ...[]FieldViolation) error {
	var flat []FieldViolation
	for i := range violations {
		flat = append(flat, violations[i]...)
	}

	if len(flat) == 0 {
		return nil
	}

	return &ValidationError{
		Violations: flat,
	}
}
