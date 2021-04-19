package main

import (
	"fmt"

	"github.com/blevesearch/bleve/v2"
	"github.com/gin-gonic/gin"
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

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
