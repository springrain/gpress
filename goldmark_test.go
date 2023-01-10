package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"gitee.com/chunanyong/gpress/parser/pageparser"
)

func TestReadMe(t *testing.T) {
	html, _ := conver2Html(datadir + "post/01-nginx-config.md")
	fmt.Println(*html)
}

func TestCreateMarkdown(t *testing.T) {
	f, _ := os.Open(datadir + "post/01-nginx-config.md")
	contentFrontMatter, err := pageparser.ParseFrontMatterAndContent(f)
	if err != nil {
		t.Error(err)
	}

	content, frontmatter := string(contentFrontMatter.Content), contentFrontMatter.FrontMatter
	for key, value := range frontmatter {
		fmt.Printf("%s:%v\n", key, value)
	}

	fmt.Println("content:" + content)
}

func TestCreateMarkdownString(t *testing.T) {
	fb, _ := os.Open(datadir + "post/01-nginx-config.md")
	//content := string(fb)

	metaSlice := []string{}
	br := bufio.NewReader(fb)
	ismeta := true
	content := ""
	for {
		a, _, c := br.ReadLine()
		b := string(a)
		if b == "---" {
			if len(metaSlice) > 0 {
				ismeta = false
			}
			continue
		}
		if ismeta {
			metaSlice = append(metaSlice, b)
		} else if c == io.EOF {
			break
		} else {
			content = content + b + "\n "
		}

		if c == io.EOF {
			break
		}
	}
	metaMap := make(map[string]interface{})
	for i := 0; i < len(metaSlice); i++ {
		v := metaSlice[i]
		strs := strings.Split(v, ": ")
		if len(strs) != 2 {
			continue
		}
		metaMap[strs[0]] = strs[1]

	}
	fmt.Println(metaMap)
	fmt.Println(content)

}
