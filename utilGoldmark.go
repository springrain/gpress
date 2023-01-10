package main

import (
	"bytes"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
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
		),
		goldmark.WithExtensions(
			extension.GFM,
			//extension.Footnote,
			meta.Meta,
			//&toc.Extender{},//不能在这里引用toc插件,手动控制
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
	parserContext := parser.NewContext(parser.WithIDs(newIDs()))
	if err := markdown.Convert(source, &htmlBuffer, parser.WithContext(parserContext)); err != nil {
		return nil, nil, nil, err
	}
	//生成页面html
	html := htmlBuffer.String()
	//读取markdown文件中的元属性
	metaData := meta.Get(parserContext)

	//生成 toc  Table of Contents,文章目录
	var tocBuffer bytes.Buffer
	tocMD := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithASTTransformers(
				util.Prioritized(&toc.Transformer{
					Title: "目录",
				}, 200),
			),
		),
		goldmark.WithExtensions(
			meta.Meta,
		),
	)
	mdParser := tocMD.Parser()

	doc := mdParser.Parse(text.NewReader(source), parser.WithContext(parserContext))
	//doc := mdParser.Parse(text.NewReader(source))
	tocTree, err := toc.Inspect(doc, source)
	if err != nil {
		return metaData, nil, &html, err
	}
	tocNode := toc.RenderList(tocTree)
	tocMD.Renderer().Render(&tocBuffer, source, tocNode)
	tocHtml := tocBuffer.String()
	return metaData, &tocHtml, &html, err
}

// 重写goldmark的autoHeadingID生成方式,兼容中文 --------------------------
type gpressMarkdownIDS struct {
	values map[string]bool
}

func newIDs() parser.IDs {
	return &gpressMarkdownIDS{
		values: map[string]bool{},
	}
}

func (s *gpressMarkdownIDS) Generate(value []byte, kind ast.NodeKind) []byte {
	value = util.TrimLeftSpace(value)
	value = util.TrimRightSpace(value)
	result := string(value)
	result = strings.ReplaceAll(result, " ", "")
	result = strings.ReplaceAll(result, ".", "-")
	if len(result) == 0 {
		if kind == ast.KindHeading {
			result = "heading"
		} else {
			result = "id"
		}
	}
	if _, ok := s.values[result]; !ok {
		s.values[result] = true
	}
	return []byte(result)
}

func (s *gpressMarkdownIDS) Put(value []byte) {
	s.values[string(value)] = true
}

//------------------------结束----------------
