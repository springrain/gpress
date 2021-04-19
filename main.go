package main

import (
	"fmt"

	"github.com/blevesearch/bleve/v2"
)

func main() {
	// open a new index
	mapping := bleve.NewIndexMapping()
	index, _ := bleve.New("zcmsIndex", mapping)

	// index some data
	_ = index.Index("zcms", "zcms 建站")

	// search for some text
	query := bleve.NewMatchQuery("zcms")
	search := bleve.NewSearchRequest(query)
	searchResults, _ := index.Search(search)
	fmt.Println(searchResults)
}
