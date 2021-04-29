package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/blevesearch/bleve/v2"
)

var indexName string = "userIndex"

//测试用例1：以字符分割
func TestCreate(t *testing.T) {
	// open a new index
	mapping := bleve.NewIndexMapping()
	bleve.New(indexName, mapping)

}

func TestSave(t *testing.T) {
	user := struct {
		Id         string
		Name       string
		Address    string
		age        int
		CreateTime time.Time
	}{
		Id:         "userId",
		Name:       "zcmsmessage",
		Address:    "zhongguo  zhengzhou",
		age:        30,
		CreateTime: time.Now(),
	}

	user2 := struct {
		Id         string
		Name       string
		Address    string
		age        int
		CreateTime time.Time
		Other      string
	}{
		Id:         "userId 2",
		Name:       "zcmsmessage 2",
		Address:    "zhongguo  zhengzhou",
		age:        30,
		CreateTime: time.Now(),
		Other:      "test Other ",
	}

	index, _ := bleve.Open(indexName)
	index.Index(user.Id, user)
	index.Index(user2.Id, user2)

}
func TestSearch1(t *testing.T) {
	index, _ := bleve.Open(indexName)
	query := bleve.NewDocIDQuery([]string{"userId"})
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, _ := index.Search(searchRequest)
	fmt.Println(searchResult)
}
func TestSearch2(t *testing.T) {
	index, _ := bleve.Open(indexName)
	query := bleve.NewQueryStringQuery("zhengzhou")
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, _ := index.Search(searchRequest)
	fmt.Println(searchResult)
}

func TestSearch3(t *testing.T) {
	index, _ := bleve.Open(indexName)
	query := bleve.NewTermQuery("zhengzhou")
	query.SetField("Address")
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, _ := index.Search(searchRequest)
	fmt.Println(searchResult)
}
