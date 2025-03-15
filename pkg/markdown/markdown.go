package markdown

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark"
)

type Markdown struct {
	md goldmark.Markdown
}

func NewMarkdown() *Markdown {
	return &Markdown{
		md: goldmark.New(),
	}
}

func (m *Markdown) ConvertToHTML(markdown string) (string, error) {
	var out bytes.Buffer
	if err := m.md.Convert([]byte(markdown), &out); err != nil {
		return "", fmt.Errorf("convert markdown to html: %w", err)
	}

	return out.String(), nil
}
