package main

import (
	"context"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
)

//获取表中符合条件字段
//indexName: 表名
//isRequired: 是否可以为空
func findIndexFields(indexName string, isRequired int) (result *bleve.SearchResult, err error) {
	//打开文件
	index := IndexMap[indexFieldIndexName]
	if err != nil {
		FuncLogError(err)
		return nil, err
	}

	var query *query.ConjunctionQuery

	// 查询指定表
	queryIndexCode := bleve.NewTermQuery(indexName)
	//查询指定字段
	queryIndexCode.SetField("IndexCode")
	if isRequired != 1 && isRequired != 0 {
		query = bleve.NewConjunctionQuery(queryIndexCode)

	} else {
		var f float64 = float64(isRequired)
		inclusive := true
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
