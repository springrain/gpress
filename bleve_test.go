package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/blevesearch/bleve/v2"

	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
)

var indexName string = "userIndex"

//创建索引
func TestCreate(t *testing.T) {

	// open a new index
	mapping := bleve.NewIndexMapping()
	//userMapping := bleve.NewDocumentMapping()
	//userMapping.AddFieldMappingsAt("Address", bleve.NewBooleanFieldMapping())
	//mapping.DefaultMapping = userMapping

	// Address的mapping映射,此字段不使用分词,只保存,用于term的绝对精确查询,类似 sql的 where = 条件查询
	addressMapping := bleve.NewTextFieldMapping()
	//不分词,只保存,需要使用keyword分词器,注意是 github.com/blevesearch/bleve/v2/analysis/analyzer/keyword, token下也有一个keyword,别引错了!!!!!
	//addressMapping.Index = false
	//addressMapping.SkipFreqNorm = true
	addressMapping.DocValues = false
	addressMapping.Analyzer = keyword.Name

	//设置字段映射
	mapping.DefaultMapping.AddFieldMappingsAt("Address", addressMapping)
	_, err := bleve.New(indexName, mapping)
	fmt.Println(err)

}

//保存数据
func TestSave(t *testing.T) {
	user := struct {
		Id         string
		Name       string
		Address    string
		age        int
		CreateTime time.Time
	}{
		Id:         "userId",
		Name:       "测试中文名称",
		Address:    "中国  zhengzhou",
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
		Name:       "测试中文名称 2",
		Address:    "zhongguo  zhengzhou",
		age:        30,
		CreateTime: time.Now(),
		Other:      "test Other ",
	}

	index, _ := bleve.Open(indexName)
	index.Index(user.Id, user)
	index.Index(user2.Id, user2)

}

// 根据ID查询
func TestSearch1(t *testing.T) {
	index, _ := bleve.Open(indexName)
	query := bleve.NewDocIDQuery([]string{"userId"})
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, _ := index.Search(searchRequest)
	fmt.Println(searchResult)
}

//根据关键字查询
func TestSearch2(t *testing.T) {
	index, _ := bleve.Open(indexName)
	query := bleve.NewQueryStringQuery("zhongguo  zhengzhou")

	searchRequest := bleve.NewSearchRequest(query)
	searchResult, _ := index.Search(searchRequest)
	fmt.Println(searchResult)
}

//查询指定的字段
func TestSearch3(t *testing.T) {
	index, _ := bleve.Open(indexName)
	//查询的关键字,需要找到绝对匹配的方式,目前还是分词匹配
	query := bleve.NewTermQuery("zhongguo  zhengzhou")
	//query := bleve.NewTermQuery("zhongguo  zhengzhou")
	//指定查询的字段
	query.SetField("Address")
	searchRequest := bleve.NewSearchRequest(query)

	//searchRequest := bleve.NewSearchRequestOptions(query, 10, 0, true)

	searchResult, _ := index.Search(searchRequest)
	fmt.Println(searchResult)
}
