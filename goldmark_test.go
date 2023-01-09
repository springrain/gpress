package main

import (
	"fmt"
	"testing"
)

func TestReadMe(t *testing.T) {
	html, _ := conver2Html("./README.md")
	fmt.Println(*html)
}
