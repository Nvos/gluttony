package unit

type Variant string
type System string

const (
	Metric   System = "METRIC"
	Imperial System = "IMPERIAL"

	Volume      Variant = "VOLUME"
	Temperature Variant = "TEMPERATURE"
	Weight      Variant = "WEIGHT"
)

type Unit struct {
	Name    string
	Symbol  string
	System  string
	Variant string
}
