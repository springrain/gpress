package service

import (
	"context"
	"fmt"
	"gitee.com/gpress/gpress/configs"
	"gitee.com/gpress/gpress/logger"
	"gitee.com/gpress/gpress/util"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
)

// 是否包含
var inclusive = true

// findIndexFieldResult 获取表中符合条件字段
// indexName: 表名/索引名
// isRequired: 是否可以为空
func findIndexFieldResult(ctx context.Context, indexName string, isRequired int) (*bleve.SearchResult, error) {
	var queryBleve *query.ConjunctionQuery
	index := configs.IndexMap[configs.INDEX_FIELD_INDEX_NAME]
	// 查询指定表
	queryIndexCode := bleve.NewTermQuery(indexName)
	// 查询指定字段
	queryIndexCode.SetField("IndexCode")
	if isRequired != 1 && isRequired != 0 {
		queryBleve = bleve.NewConjunctionQuery(queryIndexCode)
	} else {
		var f = float64(isRequired)
		queryIsRequired := bleve.NewNumericRangeInclusiveQuery(&f, &f, &inclusive, &inclusive)
		queryIsRequired.SetField("Required")
		queryBleve = bleve.NewConjunctionQuery(queryIndexCode, queryIsRequired)
	}

	// query: 条件  size:大小  from :起始
	searchRequest := bleve.NewSearchRequestOptions(queryBleve, 1000, 0, false)
	// 查询所有字段
	searchRequest.Fields = []string{"*"}

	// 按照 SortNo 正序排列.
	// 先将按"SortNo"字段对结果进行排序.如果两个文档在此字段中具有相同的值,则它们将按得分(_score)降序排序,如果文档具有相同的SortNo和得分,则将按文档ID(_id)升序排序.
	searchRequest.SortBy([]string{"SortNo", "-_score", "_id"})

	searchResult, err := index.SearchInContext(ctx, searchRequest)
	if err != nil {
		logger.FuncLogError(err)
		return nil, err
	}
	return searchResult, nil
}

// 保存新索引
func saveNewIndex(ctx context.Context, newIndex map[string]interface{}, tableName string) (map[string]string, error) {
	SearchResult, err := findIndexFieldResult(ctx, tableName, 1)
	m := make(map[string]string, 2)

	if err != nil {
		logger.FuncLogError(err)
		m["code"] = "303"
		m["msg"] = "查询异常"
		return m, err
	}
	id := util.FuncGenerateStringID()
	newIndex["ID"] = id
	result := SearchResult.Hits

	for _, v := range result {
		tmp := fmt.Sprintf("%v", v.Fields["FieldCode"]) // 转为字符串
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
	err = configs.IndexMap[tableName].Index(id, newIndex)

	if err != nil {
		logger.FuncLogError(err)
		m["code"] = "304"
		m["msg"] = "建立索引异常"
		return m, err
	}

	m["code"] = "200"
	m["msg"] = "保存成功"
	return m, nil
}

func updateIndex(ctx context.Context, tableName string, indexId string, newMap map[string]interface{}) error {
	// 查出原始数据
	index := configs.IndexMap[tableName]                 // 拿到index
	queryIndex := bleve.NewDocIDQuery([]string{indexId}) // 查询索引
	// queryIndex := bleve.NewTermQuery(indexId)            //查询索引
	// queryIndex.SetField("ID")
	searchRequest := bleve.NewSearchRequestOptions(queryIndex, 1000, 0, false)
	searchRequest.Fields = []string{"*"} // 查询所有字段
	result, err := index.SearchInContext(ctx, searchRequest)
	if err != nil {
		logger.FuncLogError(err)
		return err
	}
	// 如果没有查出来数据 证明数据错误
	if len(result.Hits) <= 0 {
		logger.FuncLogError(err)
		return fmt.Errorf("此数据不存在 ,请检查数据")
	}
	oldMap := result.Hits[0].Fields

	for k, v := range oldMap {
		newV := v
		if _, ok := newMap[k]; !ok {
			// 如果key不存在
			newMap[k] = newV
		}
	}
	err = index.Index(indexId, newMap)
	if err != nil {
		return err
	}
	return nil
}

func deleteAll(ctx context.Context, tableName string) error {
	index := configs.IndexMap[tableName]
	count, err := index.DocCount()
	if err != nil {
		return err
	}
	queryBleve := bleve.NewQueryStringQuery("*")
	// 只查一条
	searchRequest := bleve.NewSearchRequestOptions(queryBleve, int(count), 0, false)
	// 只查询id
	searchRequest.Fields = []string{"id"}

	result, err := index.SearchInContext(ctx, searchRequest)
	if err != nil {
		return err
	}

	for i := 0; i < len(result.Hits); i++ {
		err = index.Delete(result.Hits[i].ID)
		if err != nil {
			return err
		}
	}
	return nil
}
