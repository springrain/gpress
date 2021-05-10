package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/yanyiwu/gojieba"

	"github.com/blevesearch/bleve/v2"

	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
)

var indexName string = "userIndex"

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
		Age        int
		CreateTime time.Time
	}{
		Id:         "userId",
		Name:       "测试中文名称",
		Address:    "中国  zhengzhou",
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
		Address:    "zhongguo  zhengzhou",
		Age:        35,
		CreateTime: time.Now(),
		Other:      "test Other ",
	}

	user3 := make(map[string]interface{})
	user3["Id"] = "userId 3"
	user3["Name"] = "测试中文名称 3"
	user3["Address"] = "中国 北京"
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
	searchRequest.Fields = []string{"Age"}

	searchResult, _ := index.Search(searchRequest)
	fmt.Println(searchResult)
}

//精确查询指定的字段,类似SQL语句中的 where name='abc' ,要求name 字段必须使用keyword分词器
func TestSearchJingQue(t *testing.T) {
	index, _ := bleve.Open(indexName)
	//查询的关键字,使用keyword分词器,不对Adress字段分词,精确匹配
	query := bleve.NewTermQuery("zhongguo  zhengzhou")
	//query := bleve.NewTermQuery("zhongguo  zhengzhou")
	//指定查询的字段
	query.SetField("Address")
	searchRequest := bleve.NewSearchRequest(query)

	//searchRequest := bleve.NewSearchRequestOptions(query, 10, 0, true)

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

func TestGojieba(t *testing.T) {
	str := "我是一个测试语句"
	var words []string
	use_hmm := true
	x := gojieba.NewJieba()
	defer x.Free()

	words = x.CutForSearch(str, use_hmm)
	fmt.Println(str)
	fmt.Println("搜索引擎模式:", strings.Join(words, "/"))

	index, _ := bleve.Open(indexName)

	for i := 0; i < len(words); i++ {
		jiebas := make(map[string]interface{})
		jiebas["index"] = words[i]
		jiebas["str"] = str
		index.Index(strconv.FormatInt(time.Now().Unix(), 10), jiebas)
	}

	queryKey := bleve.NewQueryStringQuery("测试")

	searchRequest := bleve.NewSearchRequest(queryKey)

	//指定返回的字段
	searchRequest.Fields = []string{"str"}

	searchResult, _ := index.Search(searchRequest)
	fmt.Println(searchResult)
}
