package recipe

import (
	"bytes"

	"github.com/yuin/goldmark"
)

type MarkdownPreview struct {
	md goldmark.Markdown
}

func NewMarkdownPreview() *MarkdownPreview {
	return &MarkdownPreview{md: goldmark.New()}
}

func (m *MarkdownPreview) Preview(content string) (string, error) {
	var buf bytes.Buffer

	err := m.md.Convert([]byte(content), &buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
