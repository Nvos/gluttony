package pagination

import "fmt"

type OffsetPagination struct {
	Offset int32
	Limit  int32
}

func NewOffsetPagination(offset, limit int32) (OffsetPagination, error) {
	if offset < 0 {
		return OffsetPagination{}, fmt.Errorf("offset=%d has to be positive", offset)
	}

	if limit <= 0 {
		return OffsetPagination{}, fmt.Errorf("limit=%d has to be greater than 0", limit)
	}

	return OffsetPagination{
		Offset: offset,
		Limit:  limit,
	}, nil
}
