package main

import (
	"bytes"
	"os"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

var md goldmark.Markdown

func init() {
	md = goldmark.New(
		goldmark.WithExtensions(extension.GFM,
			extension.Footnote,
			meta.Meta,
		),
	)
}

func conver2Html(mkfile string) (map[string]interface{}, *string, error) {
	source, err := os.ReadFile(mkfile)
	html := ""
	if err != nil {
		return nil, &html, err
	}
	var htmlBuffer bytes.Buffer
	context := parser.NewContext()
	if err := md.Convert([]byte(source), &htmlBuffer, parser.WithContext(context)); err != nil {
		return nil, &html, err
	}
	metaData := meta.Get(context)
	html = htmlBuffer.String()
	return metaData, &html, err
}
