package templates

type AlertVariant string

const (
	AlertError AlertVariant = "error"
)

type Alert struct {
	Variant AlertVariant
	Title   string
	Message string
}

func NewAlert(variant AlertVariant, title string, message string) Alert {
	return Alert{Variant: variant, Title: title, Message: message}
}
