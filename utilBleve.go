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

	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search/query"
)

// 全局存放 索引对象,启动之后,所有的索引都通过这个map获取,一个索引只能打开一次,类似数据库连接,用一个对象操作
var IndexMap map[string]bleve.Index = make(map[string]bleve.Index)

// 逗号分词器的mapping
var commaAnalyzerMapping *mapping.FieldMapping = bleve.NewTextFieldMapping()

// 中文分词器的mapping
var gseAnalyzerMapping *mapping.FieldMapping = bleve.NewTextFieldMapping()

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

// 初始化 bleve 索引
func checkBleveStatus() bool {

	// 注册逗号分词器
	initRegisterCommaAnalyzer()
	// 注册gse中文分词器
	initRegistergseAnalyzer()

	// 初始化分词器
	commaAnalyzerMapping.DocValues = false
	commaAnalyzerMapping.Analyzer = commaAnalyzerName
	gseAnalyzerMapping.DocValues = false
	gseAnalyzerMapping.Analyzer = gseAnalyzerName

	status, err := checkBleveCreate()
	if err != nil {
		FuncLogError(err)
		return false
	}
	return status
}

// checkBleveCreate 检查是不是初始化安装,如果是就创建文件夹目录
func checkBleveCreate() (bool, error) {
	if pathExists(bleveDataDir) { // 如果已经存在目录,遍历索引,放到全局map里
		fileInfo, _ := os.ReadDir(bleveDataDir)
		for _, dir := range fileInfo {
			if !dir.IsDir() {
				continue
			}

			// 打开所有的索引,放到map里,一个索引只能打开一次.
			index, err := bleve.Open(bleveDataDir + dir.Name())
			if err != nil {
				return false, err
			}
			IndexMap[bleveDataDir+dir.Name()] = index
		}
		return true, nil
	}
	// 如果是初次安装,创建数据目录,默认的 ./gpressdatadir 必须存在,页面模板文件夹 ./gpressdatadir/template
	errMkdir := os.Mkdir(bleveDataDir, os.ModePerm)
	if errMkdir != nil {
		FuncLogError(errMkdir)
		return false, errMkdir
	}
	// 初始化IndexInfo
	initIndexInfo()
	// 初始化IndexField
	initIndexField()

	// 初始化配置
	initConfig()
	// 初始化用户表
	initUser()

	// 初始化站点信息表
	initSite()

	// 初始化文章默认模型的记录,indexInfo indexType="module". 只是记录,并不创建index,全部保存到context里,用于全局检索
	initModuleDefault()

	// 初始化文章内容
	initContent()

	// 初始化导航菜单
	initNavMenu()

	// 初始化页面模板
	initpageTemplateName()
	return true, nil
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
// Required: 字段是否可以为空,0查询所有字段,1查询必填字段
func findIndexFieldResult(ctx context.Context, indexName string, Required int) (*bleve.SearchResult, error) {
	var queryBleve *query.ConjunctionQuery
	index := IndexMap[indexFieldName]
	// 查询指定表
	queryIndexCode := bleve.NewTermQuery(indexName)
	// 查询指定字段
	queryIndexCode.SetField("IndexCode")
	if Required == 0 { //可以为空
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
	// 先将按"sortNo"字段对结果进行排序.如果两个文档在此字段中具有相同的值,则它们将按得分(_score)降序排序,如果文档具有相同的SortNo和得分,则将按文档ID(_id)升序排序.
	searchRequest.SortBy([]string{"sortNo", "-_score", "_id"})

	searchResult, err := index.SearchInContext(ctx, searchRequest)
	return searchResult, err
}

// findIndexFieldStruct 获取表中符合条件字段,返回Struct对象
// indexName: 表名/索引名
func findIndexFieldStruct(ctx context.Context, indexName string) ([]IndexFieldStruct, error) {
	searchResult, err := findIndexFieldResult(ctx, indexName, 0)
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
	searchResult, err := findIndexFieldResult(ctx, tableName, 1)

	responseData := ResponseData{StatusCode: 1}
	if err != nil {
		FuncLogError(err)
		responseData.StatusCode = 303
		responseData.Message = "查询异常"
		return responseData, err
	}
	id := FuncGenerateStringID()
	newIndex["id"] = id
	result := searchResult.Hits

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
	responseData.Data = id
	return responseData, err
}

func updateIndex(ctx context.Context, tableName string, indexId string, newMap map[string]interface{}) error {
	// 查出原始数据
	index := IndexMap[tableName]                         // 拿到index
	queryIndex := bleve.NewDocIDQuery([]string{indexId}) // 查询索引
	// queryIndex := bleve.NewTermQuery(indexId)            //查询索引
	// queryIndex.SetField("id")
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
		return errors.New("此数据不存在 ,请检查数据")
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
