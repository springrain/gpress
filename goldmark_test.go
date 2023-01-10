package main

import (
	"fmt"
	"testing"
)

func TestGoldmarkMeta(t *testing.T) {
	metaData, html, _ := conver2Html(datadir + "post/01-nginx-config.md")
	fmt.Println(metaData)
	fmt.Println(*html)
}
