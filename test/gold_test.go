package test

import (
	"fmt"
	"gitee.com/gpress/gpress/constant"
	"gitee.com/gpress/gpress/gmarkdown"
	"testing"
)

func TestGoldmarkMeta(t *testing.T) {
	metaData, tocHtml, html, _ := gmarkdown.Conver2Html(constant.DATA_DIR + "post/01-nginx-config.md")
	fmt.Println(metaData)
	fmt.Println(*tocHtml)
	fmt.Println(*html)
}
