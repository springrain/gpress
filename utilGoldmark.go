package main

import (
	"bytes"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var markdown goldmark.Markdown

func init() {
	markdown = goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Footnote),
	)
}

func conver2Html(mkfile string) (*string, error) {
	source, err := os.ReadFile(mkfile)
	html := ""
	if err != nil {
		return &html, err
	}
	var buf bytes.Buffer
	if err := markdown.Convert(source, &buf); err != nil {
		return &html, err
	}
	html = buf.String()
	return &html, err
}
