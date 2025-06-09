package pagination

import "math"

const Limit = 20

type Paginator struct {
	Page       int32
	TotalCount int32
	PrevPage   int32
	NextPage   int32
	HasNext    bool
	HasPrev    bool
}

func New(page int32, totalCount int64) Paginator {
	pageCount := int32(math.Ceil(float64(totalCount) / float64(Limit)))

	hasNext := page+1 < pageCount
	hasPrev := page > 0

	return Paginator{
		Page:       page,
		PrevPage:   page - 1,
		NextPage:   page + 1,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
		TotalCount: pageCount,
	}
}

type Offset struct {
	Offset int32
	Limit  int32
}

func OffsetFromPage(page int32) Offset {
	return Offset{
		Offset: page * Limit,
		Limit:  Limit,
	}
}

func (o Offset) Page() int32 {
	return o.Offset / o.Limit
}

type Page[T any] struct {
	TotalCount int64
	Rows       []T
}
