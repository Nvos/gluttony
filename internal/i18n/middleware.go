package i18n

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type localeKey struct{}

func GetLocale(ctx context.Context) (Locale, error) {
	value, ok := ctx.Value(localeKey{}).(string)
	if !ok {
		return "", fmt.Errorf("context has no locale")
	}

	locale, err := NewLocale(value)
	if err != nil {
		return "", err
	}

	return locale, nil
}

func LocaleInjectionMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			value := r.Header.Get("Accept-Language")
			languages, ok := parseAcceptLanguage(value)
			if !ok || len(languages) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			resolved := EnLocale
			if len(languages) > 0 {
				sort.Slice(languages, func(i, j int) bool {
					return languages[i].Weight < languages[i].Weight
				})

				resolved = languages[0].Locale
			}

			nextCtx := context.WithValue(r.Context(), localeKey{}, resolved)
			next.ServeHTTP(w, r.WithContext(nextCtx))
		})
	}
}

type AcceptLanguage struct {
	Locale string
	Weight float32
}

// TODO, 24/04/2024: add parser test
// TODO, 24/04/2024: improve parser
func parseAcceptLanguage(value string) ([]AcceptLanguage, bool) {
	parts := strings.Split(value, ",")
	if value == "" {
		return nil, false
	}

	all := make([]AcceptLanguage, 0, len(parts))
	for i := range parts {
		split := strings.Split(parts[i], ";")
		locale := split[0][0:2]
		var weight float32 = 1

		if len(split) == 2 {
			parsed, err := strconv.ParseFloat(split[1], 32)
			if err != nil {
				return nil, false
			}

			weight = float32(parsed)
		}

		// Skip unknown locale
		switch locale {
		case EnLocale:
		case PlLocale:
		default:
			continue
		}

		all = append(all, AcceptLanguage{
			Locale: locale,
			Weight: weight,
		})
	}

	return all, true
}
