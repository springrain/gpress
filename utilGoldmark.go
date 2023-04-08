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
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
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

// conver2Html 由markdown转成html
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

// initHighlighting 代码高亮的配置
func initHighlighting() goldmark.Extender {
	// var css bytes.Buffer
	return highlighting.NewHighlighting(
		highlighting.WithStyle("monokai"),
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
			language, _ := c.Language()
			if entering {
				w.WriteString(`<div class="highlight">`)
				if language == nil {
					_, _ = w.WriteString(`<pre class="chroma"><code class="language-fallback"  data-lang="fallback" />`)
				}
			} else {
				if language == nil {
					w.WriteString(`</code></pre>`)
				}
				w.WriteString(`</div>`)
			}
		}),

		highlighting.WithCodeBlockOptions(func(c highlighting.CodeBlockContext) []chromahtml.Option {
			languageByte, _ := c.Language()
			//if !ok {
			//	return nil
			//}
			language := string(languageByte)
			if language == "" {
				language = "fallback"
			}
			wrapper := &preWrapper{language: language}
			return []chromahtml.Option{
				chromahtml.WithClasses(true),
				chromahtml.WithLineNumbers(true),
				chromahtml.TabWidth(4),
				chromahtml.LineNumbersInTable(true),
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
