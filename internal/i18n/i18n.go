package i18n

import (
	"context"
	"fmt"
)

type contextKey string

const i18nKey contextKey = "i18n"

type I18n struct {
	Bundles map[string]*Bundle
}

func NewI18n() *I18n {
	return &I18n{
		Bundles: map[string]*Bundle{
			"en": &enBundle,
			"pl": &plBundle,
		},
	}
}

func (i *I18n) Get(lang string) (*Bundle, error) {
	got, ok := i.Bundles[lang]
	if !ok {
		return nil, fmt.Errorf("no bundle for lang %q", lang)
	}

	return got, nil
}

func WithI18nBundle(ctx context.Context, bundle *Bundle) context.Context {
	return context.WithValue(ctx, i18nKey, bundle)
}

func T(ctx context.Context, key string, args ...any) string {
	b, ok := ctx.Value(i18nKey).(*Bundle)
	if !ok {
		panic("no i18n bundle in context")
	}

	return b.T(key, args...)
}
