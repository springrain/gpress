package main

import (
	"fmt"
	"os"
	"testing"
)

func TestGoldmarkMeta(t *testing.T) {
	source, err := os.ReadFile(datadir + "post/34-es-config.md")
	if err != nil {
		t.Error(err)
		return
	}
	metaData, tocHtml, html, _ := conver2Html(source)
	fmt.Println(metaData["categories"])
	fmt.Println(metaData["tags"])
	fmt.Println(metaData)
	fmt.Println(*tocHtml)
	fmt.Println(*html)
}
