package main

import (
	"context"
	"fmt"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
)

//是否包含
var inclusive = true

// findIndexFieldResult 获取表中符合条件字段
// indexName: 表名/索引名
// isRequired: 是否可以为空
func findIndexFieldResult(ctx context.Context, indexName string, isRequired int) (*bleve.SearchResult, error) {

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
	searchRequest := bleve.NewSearchRequestOptions(query, 1000, 0, false)
	//查询所有字段
	searchRequest.Fields = []string{"*"}

	// 按照 SortNo 正序排列.
	// 先将按"SortNo"字段对结果进行排序.如果两个文档在此字段中具有相同的值,则它们将按得分(_score)降序排序,如果文档具有相同的SortNo和得分,则将按文档ID(_id)升序排序.
	searchRequest.SortBy([]string{"SortNo", "-_score", "_id"})

	searchResult, err := index.SearchInContext(ctx, searchRequest)

	if err != nil {
		FuncLogError(err)
		return nil, err
	}
	return searchResult, nil
}

func saveNexIndex(ctx context.Context, newIndex map[string]interface{}, tableName string) (map[string]string, error) {

	SearchResult, err := findIndexFieldResult(ctx, tableName, 1)
	m := make(map[string]string, 2)

	if err != nil {
		FuncLogError(err)
		m["code"] = "303"
		m["msg"] = "查询异常"
		return m, err
	}
	id := FuncGenerateStringID()
	newIndex["ID"] = id
	result := SearchResult.Hits

	for _, v := range result {
		tmp := fmt.Sprintf("%v", v.Fields["FieldCode"]) //转为字符串
		_, ok := newIndex[tmp]
		if ok {
			if newIndex[tmp] == nil || fmt.Sprintf("%v", newIndex[tmp]) == "" {
				m["code"] = "401"
				m["msg"] = tmp + "不能为空"
				return m, nil
			}

		} else {
			m["code"] = "401"
			m["msg"] = tmp + "不能为空"
			return m, nil
		}

	}
	IndexMap[tableName].Index(id, newIndex)

	m["code"] = "200"
	m["msg"] = "保存成功"
	return m, nil

}
