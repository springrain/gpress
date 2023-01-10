package main

import (
	"bytes"
	"os"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/toc"
)

var markdown goldmark.Markdown

func init() {
	markdown = goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithASTTransformers(
				util.Prioritized(&toc.Transformer{
					Title: "目录",
				}, 200),
			),
		),
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			meta.Meta,
			&toc.Extender{},
		),
		/*
			goldmark.WithRendererOptions(
				renderer.WithNodeRenderers(
					util.Prioritized(extension.NewTableHTMLRenderer(), 500),
				),
			),
		*/
	)
}

func conver2Html(mkfile string) (map[string]interface{}, *string, *string, error) {
	source, err := os.ReadFile(mkfile)
	if err != nil {
		return nil, nil, nil, err
	}
	var htmlBuffer bytes.Buffer
	parserContext := parser.NewContext()
	if err := markdown.Convert(source, &htmlBuffer, parser.WithContext(parserContext)); err != nil {
		return nil, nil, nil, err
	}
	//生成页面html
	html := htmlBuffer.String()
	//读取markdown文件中的元属性
	metaData := meta.Get(parserContext)

	//生成 toc  Table of Contents,文章目录
	var tocBuffer bytes.Buffer
	parser := markdown.Parser()
	doc := parser.Parse(text.NewReader(source))
	tocTree, err := toc.Inspect(doc, source)
	if err != nil {
		return metaData, nil, &html, err
	}
	tocNode := toc.RenderList(tocTree)
	markdown.Renderer().Render(&tocBuffer, source, tocNode)
	tocHtml := tocBuffer.String()

	return metaData, &tocHtml, &html, err
}
