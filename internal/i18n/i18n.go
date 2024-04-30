package i18n

import (
	"encoding/json"
	"fmt"
)

type Locale string

const EnLocale = "en"
const PlLocale = "pl"

func NewLocale(value string) (Locale, error) {
	switch value {
	case EnLocale:
		return EnLocale, nil
	case PlLocale:
		return PlLocale, nil
	}

	return "", fmt.Errorf("unknown locale=%s", value)
}

type Field map[Locale]string

func NewField(locale Locale, text string) Field {
	return Field{
		locale: text,
	}
}

func (f Field) JSONBytes() []byte {
	bytes, err := json.Marshal(f)
	if err != nil {
		panic(fmt.Sprintf("i18n field json marshal: %s", err))
	}

	return bytes
}

func (l Locale) FullName() string {
	switch l {
	case EnLocale:
		return "english"
	case PlLocale:
		return "polish"
	}

	panic(fmt.Sprintf("unknown locale=%s", l))
}
