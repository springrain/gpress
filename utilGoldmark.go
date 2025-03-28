// Copyright (c) 2023 gpress Authors.
//
// This file is part of gpress.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"

	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/mermaid"
	"go.abhg.dev/goldmark/toc"
	// github.com/OhYee/goldmark-fenced_codeblock_extension
	// github.com/stefanfritsch/goldmark-fences
	// latex "github.com/soypat/goldmark-latex"
)

var markdown goldmark.Markdown

func init() {
	markdown = goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithExtensions(
			extension.CJK,           // 支持中日韩语言
			extension.GFM,           // github标准
			extension.Table,         // 表格
			extension.Strikethrough, // 删除线
			extension.Linkify,       // 链接自动跳转
			extension.TaskList,      // 任务列表
			// extension.Typographer,   //符号替换,替换之后不好用
			// extension.Footnote,//php
			meta.Meta,
			//&toc.Extender{},//不能在这里引用toc插件,手动控制
			emoji.Emoji,        // emoji表情
			initHighlighting(), // 代码高亮
			&mermaid.Extender{MermaidURL: funcBasePath() + "js/mermaid.min.js"},      // mermaid流程图,不使用cdn的js
			&customExtension{NodeName: "video", NodeTag: `controls="controls" src=`}, //video扩展 !video[test.mp4](test.mp4) --> <video controls="controls" src="test.mp4">test.mp4</video>
			&customExtension{NodeName: "audio", NodeTag: `controls="controls" src=`},
		),
		//goldmark.WithRenderer(initLatexRenderer()),
		/*
			goldmark.WithRendererOptions(
				renderer.WithNodeRenderers(
					util.Prioritized(extension.NewTableHTMLRenderer(), 500),
				),
			),
		*/
	)
}

// conver2Html 由Markdown转成html
func conver2Html(source []byte) (map[string]interface{}, *string, *string, error) {

	var htmlBuffer bytes.Buffer
	// 生成id时支持中文
	parserContext := parser.NewContext(parser.WithIDs(newIDs()))
	if err := markdown.Convert(source, &htmlBuffer, parser.WithContext(parserContext)); err != nil {
		return nil, nil, nil, err
	}
	// 生成页面html
	html := htmlBuffer.String()
	// 读取markdown文件中的元属性
	metaData := meta.Get(parserContext)

	// 生成 toc  Table of Contents,文章目录
	var tocBuffer bytes.Buffer

	mdParser := markdown.Parser()

	// 生成id时支持中文
	doc := mdParser.Parse(text.NewReader(source), parser.WithContext(parserContext))
	tocTree, err := toc.Inspect(doc, source)
	if err != nil {
		return metaData, nil, &html, err
	}
	tocNode := toc.RenderList(tocTree)
	if tocNode != nil {
		markdown.Renderer().Render(&tocBuffer, source, tocNode)
	}
	tocHtml := tocBuffer.String()

	return metaData, &tocHtml, &html, err
}

// 重写goldmark的autoHeadingID生成方式 --------------------------
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
	//result = strings.ReplaceAll(result, "(", "_")
	//result = strings.ReplaceAll(result, ")", "_")
	if len(result) == 0 {
		if kind == ast.KindHeading {
			result = "heading"
		} else {
			result = "id"
		}
	}
	// querySelector 不支持 数字开头,去掉所有的数字和符号,提取字符
	//result = extractLetters(result)
	// 如果是数字开头,加上 m-
	//if startsWithDigit(result) {
	//	result = "m-" + result
	//}
	if _, ok := s.values[result]; !ok {
		s.values[result] = true
	}
	return []byte(result)
}

func (s *gpressMarkdownIDS) Put(value []byte) {
	s.values[string(value)] = true
}

/*
	func extractLetters(input string) string {
		var result string

		for _, char := range input {
			// 判断字符是否是字母
			if unicode.IsLetter(char) {
				result += string(char)
			}
		}

		return result
	}

func startsWithDigit(str string) bool {
	if len(str) == 0 {
		return false
	}

	// 使用 unicode.IsDigit 判断第一个字符是否是数字
	firstChar := rune(str[0])
	return unicode.IsDigit(firstChar)
}
*/
//------------------------结束----------------

// initHighlighting 代码高亮的配置
func initHighlighting() goldmark.Extender {
	// var css bytes.Buffer
	return highlighting.NewHighlighting(
		highlighting.WithStyle("monokai"),
		// 用于处理没有指定语言的情况,例如:
		// ```
		// 不写语言
		// ```
		highlighting.WithGuessLanguage(true),
		// highlighting.WithCSSWriter(&css),
		/*
			highlighting.WithFormatOptions(
				chromahtml.WithClasses(true),
				chromahtml.WithLineNumbers(true),
				chromahtml.TabWidth(4),
				chromahtml.LineNumbersInTable(true),
				//chromahtml.InlineCode(true),
				//chromahtml.LineNumbersInTable(true),
			),
		*/
		highlighting.WithWrapperRenderer(func(w util.BufWriter, c highlighting.CodeBlockContext, entering bool) {
			//language, _ := c.Language()
			if entering {
				w.WriteString(`<div class="highlight">`)
				// 使用 highlighting.WithGuessLanguage(true),
				//if language == nil {
				//	_, _ = w.WriteString(`<pre tabindex="0" class="chroma"><code class="language-fallback"  data-lang="fallback" />`)
				//}
			} else {
				// 使用 highlighting.WithGuessLanguage(true),
				//if language == nil {
				//	w.WriteString(`</code></pre>`)
				//}
				w.WriteString(`</div>`)
			}
		}),

		highlighting.WithCodeBlockOptions(func(c highlighting.CodeBlockContext) []chromahtml.Option {
			languageByte, ok := c.Language()
			if !ok {
				return nil
			}
			language := string(languageByte)
			//if language == "" {
			//	language = "fallback"
			//}
			wrapper := &preWrapper{language: language}
			return []chromahtml.Option{

				// 暂时不显示行号,有些模板不兼容
				//chromahtml.WithLineNumbers(true),
				//chromahtml.LineNumbersInTable(true),

				chromahtml.TabWidth(4),
				chromahtml.WithClasses(true),
				chromahtml.WithPreWrapper(wrapper),
			}
		}),
	)
}

/*
// initLatexRenderer 初始化latex科学符号
func initLatexRenderer() renderer.Renderer {
	r := renderer.NewRenderer(renderer.WithNodeRenderers(util.Prioritized(latex.NewRenderer(latex.Config{
		NoHeadingNumbering: true,                                                                     // No heading numbers
		Preamble:           append(latex.DefaultPreamble(), []byte("\n\\usepackage{MnSymbol}\n")...), // add star symbols to preamble.
		DeclareUnicode: func(r rune) (raw string, isReplaced bool) {
			switch r {
			case '★':
				return `$\filledstar$`, true
			case '☆':
				return `$\smallstar$`, true
			}
			return "", false
		},
	}), 1)))

	return r
}
*/

type preWrapper struct {
	language string
}

func (p *preWrapper) Start(code bool, styleAttr string) string {
	var language string
	if code {
		language = p.language
	} else {
		language = "fallback"
	}
	w := &strings.Builder{}
	WritePreStart(w, language, styleAttr)
	//p.low = p.writeCounter.counter + w.Len()
	return w.String()
}
func WritePreStart(w io.Writer, language, styleAttr string) {
	fmt.Fprintf(w, `<pre tabindex="0"%s>`, styleAttr)
	fmt.Fprint(w, "<code")

	if language != "" {
		fmt.Fprint(w, ` class="language-`+language+`"`)
		fmt.Fprint(w, ` data-lang="`+language+`"`)
	}
	fmt.Fprint(w, ">")
}

const preEnd = "</code></pre>"

func (p *preWrapper) End(code bool) string {
	//p.high = p.writeCounter.counter
	return preEnd
}

// 自定义CustomNode标签解析
type CustomNode struct {
	ast.BaseInline
	NodeName string
	NodeKind ast.NodeKind
	Title    []byte
	URL      []byte
}

func (n *CustomNode) Dump(source []byte, level int) {
	// 可选的调试方法
}

func (n *CustomNode) Kind() ast.NodeKind {
	return n.NodeKind
}

type customParser struct {
	NodeName string
	NodeKind ast.NodeKind
}

func (s *customParser) Trigger() []byte {
	return []byte{'!'} // 检测以 '!' 开头的文本
}

func (s *customParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if len(line) < len(s.NodeName)+2 || string(line[0:len(s.NodeName)+1]) != ("!"+s.NodeName) {
		return nil // 不是 !video 语法
	}
	block.Advance(len(s.NodeName) + 1) // 跳过 "!video"

	// 解析标题 [title]
	title, ok := parseDelimitedContent(block, '[', ']')
	if !ok {
		return nil
	}

	// 解析 URL (url)
	url, ok := parseDelimitedContent(block, '(', ')')
	if !ok {
		return nil
	}
	// 创建 Video 节点
	return &CustomNode{
		NodeName: s.NodeName,
		NodeKind: s.NodeKind,
		Title:    title,
		URL:      url,
	}
}

// FindClosureOptions 配置（根据旧参数 false, false 设置）
var (
	noNestingOptions = text.FindClosureOptions{ // 不允许嵌套
		Nesting: false,
		Newline: true, // 支持跨行查找
	}
)

// parseDelimitedContent 通用函数：解析类似 [content] 或 (content) 的语法
// opener: 起始字符（如 '[', '('）
// closure: 闭合字符（如 ']', ')'）
func parseDelimitedContent(block text.Reader, opener, closer byte) ([]byte, bool) {
	line, _ := block.PeekLine()
	if len(line) == 0 || line[0] != opener {
		return nil, false
	}
	block.Advance(1) // 跳过起始符（如 '[' 或 '('）
	// 查找闭合符
	segments, found := block.FindClosure(opener, closer, noNestingOptions)
	if !found || segments.Len() == 0 {
		return nil, false
	}
	// 计算闭合符的绝对位置
	lastSegment := segments.At(segments.Len() - 1)
	content := block.Value(lastSegment)
	block.Advance(lastSegment.Len() + 1)
	return content, true
}

type customHTMLRenderer struct {
	NodeName string
	NodeKind ast.NodeKind
	NodeTag  string
	html.Config
}

func newCustomHTMLRenderer(nodeName string, nodeKind ast.NodeKind, nodeTag string, opts ...html.Option) renderer.NodeRenderer {
	r := &customHTMLRenderer{
		NodeName: nodeName,
		NodeKind: nodeKind,
		NodeTag:  nodeTag,
		Config:   html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

func (r *customHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(r.NodeKind, r.renderCustomNode)
}

func (r *customHTMLRenderer) renderCustomNode(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*CustomNode)
	title := string(n.Title)
	url := string(n.URL)

	// 生成 HTML
	_, _ = w.WriteString("<" + n.NodeName + " ")
	_, _ = w.WriteString(r.NodeTag)
	_, _ = w.WriteString(`"`)
	_, _ = w.Write(util.EscapeHTML(util.URLEscape([]byte(url), true)))
	_, _ = w.WriteString(`">`)
	_, _ = w.Write(util.EscapeHTML([]byte(title)))
	_, _ = w.WriteString("</" + n.NodeName + ">")

	return ast.WalkContinue, nil
}

// customExtension video扩展 !video[test.mp4](test.mp4) --> <video controls="controls" src="test.mp4">test.mp4</video>
type customExtension struct {
	NodeName string
	NodeTag  string
	nodeKind ast.NodeKind
}

func (ce *customExtension) Extend(m goldmark.Markdown) {
	if ce.nodeKind == 0 {
		ce.nodeKind = ast.NewNodeKind(ce.NodeName)
	}
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(&customParser{NodeName: ce.NodeName, NodeKind: ce.nodeKind}, 100), // 高优先级
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(newCustomHTMLRenderer(ce.NodeName, ce.nodeKind, ce.NodeTag), 100),
		),
	)
}
