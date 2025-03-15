package html

import (
	"fmt"
	"html"
	"html/template"
	"net/url"
	"strings"
	"time"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"formatDuration": func(duration time.Duration) string {
			return time.Unix(0, 0).UTC().Add(duration).Format("15:04")
		},
		"rawHTML": func(raw string) template.HTML {
			//nolint:gosec // false positive, string is escaped
			return template.HTML(html.EscapeString(raw))
		},
		"isURL": func(s string) bool {
			_, err := url.ParseRequestURI(s)
			return err == nil
		},
		"hasPrefix": strings.HasPrefix,
		"add": func(a, b any) int {
			return castInt(a) + castInt(b)
		},
		"queryParams": func(values ...any) string {
			if (len(values) % 2) == 1 {
				panic("invalid number of pairs")
			}

			params := url.Values{}
			for i := 0; i < len(values); i += 2 {
				name := values[i].(string)
				value := fmt.Sprintf("%v", values[i+1])

				if name == "" || value == "" {
					continue
				}

				params.Add(name, value)
			}

			return fmt.Sprintf("?%s", params.Encode())
		},
	}
}

func castInt(s any) int {
	switch v := s.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case int32:
		return int(v)
	}

	panic("unsupported cast")
}
