package pagination

const DefaultLimit = 20

type Offset struct {
	Offset int32
	Limit  int32
}

func OffsetFromPage(page int32) Offset {
	return Offset{
		Offset: page * DefaultLimit,
		Limit:  DefaultLimit,
	}
}

func (o Offset) Page() int32 {
	return o.Offset / o.Limit
}
