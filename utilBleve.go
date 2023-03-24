package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
)

// FuncGenerateStringID 默认生成字符串ID的函数.方便自定义扩展
// FuncGenerateStringID Function to generate string ID by default. Convenient for custom extension
var FuncGenerateStringID func() string = generateStringID

// generateStringID 生成主键字符串
// generateStringID Generate primary key string
func generateStringID() string {
	// 使用 crypto/rand 真随机9位数
	randNum, randErr := rand.Int(rand.Reader, big.NewInt(1000000000))
	if randErr != nil {
		return ""
	}
	// 获取9位数,前置补0,确保9位数
	rand9 := fmt.Sprintf("%09d", randNum)

	// 获取纳秒 按照 年月日时分秒毫秒微秒纳秒 拼接为长度23位的字符串
	pk := time.Now().Format("2006.01.02.15.04.05.000000000")
	pk = strings.ReplaceAll(pk, ".", "")

	// 23位字符串+9位随机数=32位字符串,这样的好处就是可以使用ID进行排序
	pk = pk + rand9
	return pk
}

// pathExists 文件或者目录是否存在
func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// result2Map 单个查询结果转map
func result2Map(indexName string, result *bleve.SearchResult) (map[string]interface{}, error) {
	if result == nil {
		return nil, errors.New("结果集为空")
	}
	if result.Total == 0 { //没有记录
		return nil, nil
	}
	if result.Total > 1 { // 大于1条记录
		return nil, errors.New("查询出多条记录")
	}
	//获取到查询的对象
	value := result.Hits[0]
	m := make(map[string]interface{}, 0)
	for k, v := range value.Fields {
		m[k] = v
	}
	return m, nil

}

// result2SliceMap 多条结果转map数组
func result2SliceMap(indexName string, result *bleve.SearchResult) ([]map[string]interface{}, error) {
	if result == nil {
		return nil, errors.New("结果集为空")
	}
	if result.Total == 0 { //没有记录
		return nil, nil
	}
	ms := make([]map[string]interface{}, 0)
	//获取到查询的对象
	for _, value := range result.Hits {
		m := make(map[string]interface{}, 0)
		for k, v := range value.Fields {
			m[k] = v
		}
		ms = append(ms, m)
	}
	return ms, nil
}

// 是否包含
var inclusive = true

// findIndexFieldResult 获取表中符合条件字段
// indexName: 表名/索引名
// Required: 是否可以为空
func findIndexFieldResult(ctx context.Context, indexName string, Required int) (*bleve.SearchResult, error) {
	var queryBleve *query.ConjunctionQuery
	index := IndexMap[indexFieldIndexName]
	// 查询指定表
	queryIndexCode := bleve.NewTermQuery(indexName)
	// 查询指定字段
	queryIndexCode.SetField("IndexCode")
	if Required != 1 && Required != 0 {
		queryBleve = bleve.NewConjunctionQuery(queryIndexCode)
	} else {
		var f = float64(Required)
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
	return searchResult, err
}

// findIndexFieldStruct 获取表中符合条件字段,返回Struct对象
// indexName: 表名/索引名
// Required: 是否可以为空
func findIndexFieldStruct(ctx context.Context, indexName string, Required int) ([]IndexFieldStruct, error) {
	searchResult, err := findIndexFieldResult(ctx, indexName, Required)
	if err != nil {
		FuncLogError(err)
		return nil, err
	}
	maps, err := result2SliceMap(indexName, searchResult)
	if err != nil {
		FuncLogError(err)
		return nil, err
	}

	fields := make([]IndexFieldStruct, 0)
	jsonStr, err := json.Marshal(maps)
	if err != nil {
		FuncLogError(err)
		return nil, err
	}
	err = json.Unmarshal(jsonStr, &fields)
	if err != nil {
		FuncLogError(err)
		return nil, err
	}
	return fields, nil
}

// 保存新索引
func saveNewIndex(ctx context.Context, newIndex map[string]interface{}, tableName string) (ResponseData, error) {
	SearchResult, err := findIndexFieldResult(ctx, tableName, 1)

	responseData := ResponseData{StatusCode: 1}
	if err != nil {
		FuncLogError(err)
		responseData.StatusCode = 303
		responseData.Message = "查询异常"
		return responseData, err
	}
	id := FuncGenerateStringID()
	newIndex["ID"] = id
	result := SearchResult.Hits

	for _, v := range result {
		tmp := v.Fields["FieldCode"].(string) // 转为字符串
		_, ok := newIndex[tmp]
		if !ok {
			responseData.StatusCode = 401
			responseData.Message = tmp + "不能为空"
			return responseData, err
		}
	}
	err = IndexMap[tableName].Index(id, newIndex)

	if err != nil {
		FuncLogError(err)
		responseData.StatusCode = 304
		responseData.Message = "建立索引异常"
		return responseData, err
	}
	responseData.StatusCode = 200
	responseData.Message = "保存成功"
	return responseData, err
}

func updateIndex(ctx context.Context, tableName string, indexId string, newMap map[string]interface{}) error {
	// 查出原始数据
	index := IndexMap[tableName]                         // 拿到index
	queryIndex := bleve.NewDocIDQuery([]string{indexId}) // 查询索引
	// queryIndex := bleve.NewTermQuery(indexId)            //查询索引
	// queryIndex.SetField("ID")
	searchRequest := bleve.NewSearchRequestOptions(queryIndex, 1000, 0, false)
	searchRequest.Fields = []string{"*"} // 查询所有字段
	result, err := index.SearchInContext(ctx, searchRequest)
	if err != nil {
		FuncLogError(err)
		return err
	}
	// 如果没有查出来数据 证明数据错误
	if len(result.Hits) <= 0 {
		FuncLogError(err)
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
	index := IndexMap[tableName]
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
