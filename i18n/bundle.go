package i18n

import "fmt"

type Bundle struct {
	Lang     string
	Messages map[string]string
}

func (i *Bundle) T(key string, args ...any) string {
	message, ok := i.Messages[key]
	if !ok {
		panic(fmt.Sprintf("no message for key=%q", key))
	}

	return fmt.Sprintf(message, args...)
}
