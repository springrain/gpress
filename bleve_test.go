package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/blevesearch/bleve/v2"
)

var indexName string = "testIndex"

//创建索引
func TestCreate(t *testing.T) {

	os.RemoveAll(indexName)
	// open a new index
	mapping := bleve.NewIndexMapping()
	//userMapping := bleve.NewDocumentMapping()
	//userMapping.AddFieldMappingsAt("Address", bleve.NewBooleanFieldMapping())
	//mapping.DefaultMapping = userMapping

	// Address的mapping映射,此字段不使用分词,只保存,用于term的绝对精确查询,类似 sql的 where = 条件查询
	addressMapping := bleve.NewTextFieldMapping()

	//addressMapping.Index = false
	//addressMapping.SkipFreqNorm = true
	addressMapping.DocValues = false
	//不分词,只保存,需要使用keyword分词器,注意是 github.com/blevesearch/bleve/v2/analysis/analyzer/keyword, token下也有一个keyword,别引错了!!!!!
	//addressMapping.Analyzer = keyword.Name
	//gse中文分词器
	//addressMapping.Analyzer = gseAnalyzerName
	// 逗号(,)分词器
	addressMapping.Analyzer = commaAnalyzerName

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
		Age        int
		CreateTime time.Time
	}{
		Id:         "userId",
		Name:       "测试中文名称",
		Address:    "中国,zhengzhou",
		Age:        30,
		CreateTime: time.Now(),
	}

	user2 := struct {
		Id         string
		Name       string
		Address    string
		Age        int
		CreateTime time.Time
		Other      string
	}{
		Id:         "userId 2",
		Name:       "测试中文名称 2",
		Address:    "zhongguo,zhengzhou",
		Age:        35,
		CreateTime: time.Now(),
		Other:      "test Other ",
	}

	user3 := make(map[string]interface{})
	user3["Id"] = "userId 3"
	user3["Name"] = "测试中文名称 3"
	user3["Address"] = "完,美"
	user3["Age"] = 36
	user3["CreateTime"] = time.Now()

	index, _ := bleve.Open(indexName)
	index.Index(user.Id, user)
	index.Index(user2.Id, user2)
	index.Index("userId 3", user3)

}

// 根据ID查询
func TestSearchID(t *testing.T) {
	index, _ := bleve.Open(indexName)
	query := bleve.NewDocIDQuery([]string{"userId"})
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, _ := index.Search(searchRequest)
	fmt.Println(searchResult)
}

//根据关键字查询
func TestSearchKey(t *testing.T) {
	index, _ := bleve.Open(indexName)
	queryKey := bleve.NewQueryStringQuery("中文2")

	searchRequest := bleve.NewSearchRequest(queryKey)

	//指定返回的字段
	searchRequest.Fields = []string{"*"}

	searchResult, _ := index.Search(searchRequest)
	fmt.Println(searchResult)
}

//精确查询指定的字段,类似SQL语句中的 where name='abc' ,要求name 字段必须使用keyword分词器
func TestSearchJingQue(t *testing.T) {
	index, _ := bleve.Open(indexName)
	//查询的关键字,使用keyword分词器,不对Adress字段分词,精确匹配
	query := bleve.NewTermQuery("完美")
	//query := bleve.NewTermQuery("zhongguo  zhengzhou")
	//指定查询的字段
	query.SetField("Address")
	searchRequest := bleve.NewSearchRequest(query)

	//searchRequest := bleve.NewSearchRequestOptions(query, 10, 0, true)
	//查询所有的字段
	searchRequest.Fields = []string{"*"}
	searchResult, _ := index.Search(searchRequest)
	fmt.Println(searchResult)
}

//数字范围查询
func TestSearchNum(t *testing.T) {
	index, _ := bleve.Open(indexName)
	//查询的关键字,使用keyword分词器,不对Adress字段分词,精确匹配
	var min float64 = 20.00
	var max float64 = 32.00
	querynum := bleve.NewNumericRangeQuery(&min, &max)
	//指定查询的字段
	querynum.SetField("Age")

	searchRequest := bleve.NewSearchRequest(querynum)

	//searchRequest := bleve.NewSearchRequestOptions(query, 10, 0, true)

	searchResult, _ := index.Search(searchRequest)
	fmt.Println(searchResult)
}

//日期范围查询
func TestSearchDate(t *testing.T) {
	index, _ := bleve.Open(indexName)
	//查询的关键字,使用keyword分词器,不对Adress字段分词,精确匹配
	start := time.Now().Add(time.Hour * -1)
	end := time.Now()
	querynum := bleve.NewDateRangeQuery(start, end)
	//指定查询的字段
	querynum.SetField("CreateTime")

	searchRequest := bleve.NewSearchRequest(querynum)

	//searchRequest := bleve.NewSearchRequestOptions(query, 10, 0, true)

	searchResult, _ := index.Search(searchRequest)
	fmt.Println(searchResult)
}

//复合查询,类似SQL中的 WHERE 后的条件语句
func TestSearchWhere(t *testing.T) {
	index, _ := bleve.Open(indexName)
	//查询的关键字,使用keyword分词器,不对Adress字段分词,精确匹配
	var min float64 = 20.00
	var max float64 = 40.00
	queryNum := bleve.NewNumericRangeQuery(&min, &max)
	//指定查询的字段
	queryNum.SetField("Age")

	queryKey := bleve.NewQueryStringQuery("中文")

	//多个条件联查
	query := bleve.NewConjunctionQuery(queryNum, queryKey)

	searchRequest := bleve.NewSearchRequest(query)

	//查询所有的字段
	searchRequest.Fields = []string{"*"}
	//指定返回的字段
	//searchRequest.Fields = []string{"Name", "Age"}

	searchResult, _ := index.Search(searchRequest)

	fmt.Println(searchResult)
}

func TestSearchOrder(t *testing.T) {
	index := IndexMap[indexFieldIndexName]
	//查询的关键字,使用keyword分词器,不对Adress字段分词,精确匹配
	query := bleve.NewTermQuery(userIndexName)
	//query := bleve.NewTermQuery("zhongguo  zhengzhou")
	//指定查询的字段
	query.SetField("IndexCode")
	searchRequest := bleve.NewSearchRequest(query)

	// 按照 SortNo 正序排列.
	// 先将按"SortNo"字段对结果进行排序.如果两个文档在此字段中具有相同的值,则它们将按得分(_score)降序排序,如果文档具有相同的SortNo和得分,则将按文档ID(_id)升序排序.
	searchRequest.SortBy([]string{"SortNo", "-_score", "_id"})
	// 按照 SortNo 降序排列.
	//searchRequest.SortBy([]string{"-SortNo", "-_score", "_id"})

	//searchRequest := bleve.NewSearchRequestOptions(query, 10, 0, true)
	//查询所有的字段
	searchRequest.Fields = []string{"SortNo"}
	searchResult, _ := index.Search(searchRequest)
	fmt.Println(searchResult)
}
