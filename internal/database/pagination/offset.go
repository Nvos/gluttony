package pagination

import "fmt"

type OffsetPagination struct {
	Offset int32
	Limit  int32
}

func ValidateOffsetPagination(pagination OffsetPagination) error {
	if pagination.Offset < 0 {
		return fmt.Errorf("validate offset pagination: offset < 0")
	}

	if pagination.Limit < 0 {
		return fmt.Errorf("validate offset pagination: limit < 0")
	}

	return nil
}
