package main

import (
	"context"
	"fmt"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
)

//获取表中符合条件字段
//tableName: 表名
//isRequired: 是否可以为空
func getFields(tableName string, isRequired float64) (err error, result *bleve.SearchResult) {
	//打开文件
	index := IndexMap[indexFieldIndexName]
	if err != nil {
		fmt.Errorf("getFields() 异常: %v", err)
		return err, nil
	}

	var query *query.ConjunctionQuery

	// 查询指定表
	queryIndexCode := bleve.NewTermQuery(tableName)
	//查询指定字段
	queryIndexCode.SetField("FieldCode")

	if isRequired != 1 && isRequired != 0 {
		query = bleve.NewConjunctionQuery(queryIndexCode)

	} else {
		var f float64 = 0.00
		var f2 float64 = 2.00
		queryIsReqired := bleve.NewNumericRangeQuery(&f, &f2)
		//queryIsReqired := bleve.NewTermQuery("1")
		queryIndexCode.SetField("Required")
		//queryIsReqired.SetField("Required")
		query = bleve.NewConjunctionQuery(queryIndexCode, queryIsReqired)
	}

	//query: 条件  size:大小  from :起始
	serarch := bleve.NewSearchRequestOptions(query, 1000, 0, false)
	//查询所有字段
	serarch.Fields = []string{"*"}

	result, err = index.SearchInContext(context.Background(), serarch)

	if err != nil {
		fmt.Errorf("getFields() 异常: %v", err)
		return err, nil
	}

	return nil, result

}
