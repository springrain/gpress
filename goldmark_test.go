package main

import (
	"fmt"
	"testing"
)

func TestReadMe(t *testing.T) {
	html, _ := conver2Html("./post/01-nginx-config.md")
	fmt.Println(*html)
}
