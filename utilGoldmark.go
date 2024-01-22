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
			&mermaid.Extender{MermaidURL: funcBasePath() + "js/mermaid.min.js"}, // mermaid流程图,不使用cdn的js
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
