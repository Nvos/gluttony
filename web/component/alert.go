package component

type AlertVariant string

const (
	AlertError AlertVariant = "error"
)

type AlertInput struct {
	Variant AlertVariant
	Title   string
	Message string
}

func NewAlert(variant AlertVariant, title string, message string) AlertInput {
	return AlertInput{Variant: variant, Title: title, Message: message}
}
