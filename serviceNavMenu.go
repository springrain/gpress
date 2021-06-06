package main

import (
	"context"
	"fmt"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
)

//获取表中符合条件字段
//indexName: 表名
//isRequired: 是否可以为空
//打开文件
var inclusive = true

func findIndexFields(indexName string, isRequired int) (result *bleve.SearchResult, err error) {
	if err != nil {
		FuncLogError(err)
		return nil, err
	}

	var query *query.ConjunctionQuery
	var index = IndexMap[indexFieldIndexName]
	// 查询指定表
	queryIndexCode := bleve.NewTermQuery(indexName)
	//查询指定字段
	queryIndexCode.SetField("IndexCode")
	if isRequired != 1 && isRequired != 0 {
		query = bleve.NewConjunctionQuery(queryIndexCode)

	} else {

		var f float64 = float64(isRequired)
		queryIsReqired := bleve.NewNumericRangeInclusiveQuery(&f, &f, &inclusive, &inclusive)
		queryIsReqired.SetField("Required")
		query = bleve.NewConjunctionQuery(queryIndexCode, queryIsReqired)
	}

	//query: 条件  size:大小  from :起始
	serarch := bleve.NewSearchRequestOptions(query, 1000, 0, false)
	//查询所有字段
	serarch.Fields = []string{"*"}

	result, err = index.SearchInContext(context.Background(), serarch)

	if err != nil {
		FuncLogError(err)
		return nil, err
	}

	return result, nil
}

//获取菜单树 pid 为0 为一级菜单
var active float64 = 1

func getNavMenu(pid string) (interface{}, error) {

	navIndex := IndexMap[indexNavMenuName]
	//PID 跟 Active 为查询字段
	queryPID := bleve.NewTermQuery(pid)
	queryPID.SetField("PID")
	queryActive := bleve.NewNumericRangeInclusiveQuery(&active, &active, &inclusive, &inclusive)
	queryActive.SetField("Active")
	query := bleve.NewConjunctionQuery(queryPID, queryActive)
	serarch := bleve.NewSearchRequestOptions(query, 1000, 0, false)
	serarch.Fields = []string{"*"}
	result, err := navIndex.SearchInContext(context.Background(), serarch)

	if err != nil {
		FuncLogError(err)
		return nil, err
	}
	data := make([]map[string]interface{}, len(result.Hits))
	for i, v := range result.Hits {
		id := fmt.Sprintf("%v", v.Fields["ID"])

		if id != "" && id != "nil" {
			value, _ := getNavMenu(id)
			v.Fields["childs"] = value
		}
		data[i] = v.Fields
	}

	return data, err

}
