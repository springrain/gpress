package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/cloudwego/hertz/pkg/app"

	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search/query"
)

// 全局存放 索引对象,启动之后,所有的索引都通过这个map获取,一个索引只能打开一次,类似数据库连接,用一个对象操作
var IndexMap map[string]bleve.Index = make(map[string]bleve.Index)

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

	// 注册bleve分词器,包括逗号,中文,小写keyword
	initRegisterAnalyzer()

	// 初始化分词器
	commaAnalyzerMapping.DocValues = false
	commaAnalyzerMapping.Analyzer = commaAnalyzerName
	gseAnalyzerMapping.DocValues = false
	gseAnalyzerMapping.Analyzer = gseAnalyzerName
	keywordAnalyzerMapping.DocValues = false
	keywordAnalyzerMapping.Analyzer = keywordAnalyzerName

	if !pathExists(bleveDataDir) { //目录如果不存在
		// 如果是初次安装,创建数据目录,默认的 ./gpressdatadir 必须存在,页面模板文件夹 ./gpressdatadir/template
		err := os.Mkdir(bleveDataDir, os.ModePerm)
		if err != nil {
			FuncLogError(err)
			return false
		}
	}

	//这三张表是系统表,使用变量初始化,优先级高于init,其他表使用 init函数初始化

	// 初始化indexField
	_, err := initIndexField()
	if err != nil {
		return false
	}
	// 初始化indexInfo
	_, err = initIndexInfo()
	if err != nil {
		return false
	}
	// 初始化 config
	ok, err := initConfig()
	if err != nil {
		return false
	}
	return ok
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
// required: 字段是否可以为空,0查询所有字段,1查询必填字段
func findIndexFieldResult(ctx context.Context, indexName string, required int) (*bleve.SearchResult, error) {
	var queryBleve query.Query
	index, _, _ := openBleveIndex(indexFieldName)
	// 查询指定表
	queryIndexCode := bleveNewTermQuery(indexName)
	// 查询指定字段,和json字段保持一致
	queryIndexCode.SetField("indexCode")
	if required == 0 { //可以为空
		queryBleve = queryIndexCode
	} else {
		var f = float64(required)
		queryIsRequired := bleve.NewNumericRangeInclusiveQuery(&f, &f, &inclusive, &inclusive)
		queryIsRequired.SetField("required") // 查询指定字段,和json字段保持一致
		queryBleve = bleve.NewConjunctionQuery(queryIndexCode, queryIsRequired)
	}

	// query: 条件  size:大小  from :起始
	searchRequest := bleve.NewSearchRequestOptions(queryBleve, 1000, 0, false)
	// 查询所有字段
	searchRequest.Fields = []string{"*"}

	// 按照 SortNo 升序排列.
	// 先将按"sortNo"字段对结果进行排序.如果两个文档在此字段中具有相同的值,则将按文档ID(_id)降序排序.
	searchRequest.SortBy([]string{"sortNo", "-_id"})
	//searchRequest.SortBy([]string{"sortNo"})

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
func saveNewIndex(ctx context.Context, tableName string, newIndex map[string]interface{}) (ResponseData, error) {
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
		tmp := v.Fields["fieldCode"].(string) // 转为字符串
		_, ok := newIndex[tmp]
		if !ok {
			responseData.StatusCode = 401
			responseData.Message = tmp + "不能为空"
			return responseData, err
		}
	}
	index, _, _ := openBleveIndex(tableName)
	err = index.Index(id, newIndex)

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
	index, _, _ := openBleveIndex(tableName)             // 拿到index
	queryIndex := bleve.NewDocIDQuery([]string{indexId}) // 查询索引
	// queryIndex := bleveNewTermQuery(indexId)            //查询索引
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
func deleteById(ctx context.Context, tableName string, id string) error {
	index, ok, _ := openBleveIndex(tableName)
	if !ok {
		return errors.New("数据不存在")
	}
	err := index.Delete(id)
	return err
}
func deleteAll(ctx context.Context, tableName string) error {
	index, _, _ := openBleveIndex(tableName)
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

func findIndexList(ctx context.Context, c *app.RequestContext, indexName string) (ResponseData, error) {
	searchIndex, ok, _ := openBleveIndex(indexName)
	if !ok { //索引不存在
		err := errors.New("索引不存在")
		return ResponseData{StatusCode: 0, ERR: err}, err
	}
	//获取页码
	pageNoStr := c.DefaultQuery("pageNo", "1")
	pageNo, _ := strconv.Atoi(pageNoStr)
	if pageNo == 0 {
		pageNo = 1
	}
	mapParams := make(map[string]interface{}, 0)
	//获取所有的参数
	c.Bind(&mapParams)
	//删除掉固定的两个
	delete(mapParams, "pageNo")
	delete(mapParams, "q")
	// 查询
	var searchQuery query.Query
	var queryKey query.Query
	q := strings.TrimSpace(c.Query("q"))
	if q == "" || q == "*" {
		queryKey = bleve.NewQueryStringQuery("*")
	} else {
		//gse分词之后,内容就被拆解用于索引了,也就是无法通过完整的字符串精确匹配了,所以需要把索引的默认分词器设置为gse,这样搜索词也被分词了
		// NewQueryStringQuery 在搜索所有字段时,使用的是索引的 mapping.DefaultAnalyzer 分词器,使用keyword的就需要明确声明了

		//不对q分词搜索,精确匹配
		//termQuery := bleve.NewTermQuery(q)
		//matchAllQuery := bleve.NewMatchAllQuery()
		//queryBoolean1 := bleve.NewBooleanQuery()
		//queryBoolean1.AddMust(termQuery, matchAllQuery)
		//不对q分词搜索,精确匹配,NewQueryStringQuery
		//queryBoolean1 := bleve.NewQueryStringQuery("\"" + q + "\"")
		//queryBoolean1.SetBoost(100)

		//对q分词搜索
		//queryBoolean2 := bleve.NewQueryStringQuery(q)

		//queryBoolean := bleve.NewBooleanQuery()
		//queryBoolean.AddShould(queryBoolean1, queryBoolean2)
		//queryKey = queryBoolean

		queryKey = bleve.NewQueryStringQuery(q)

	}

	if len(mapParams) < 1 { //没有其他参数了
		searchQuery = queryKey
	} else { //还有其他参数,认为是数据库字段,进行检索
		qs := make([]query.Query, 0)
		qs = append(qs, queryKey)
		for k := range mapParams {
			value := c.Query(k)
			term := bleveNewTermQuery(value)
			term.SetField(k)
			qs = append(qs, term)
		}
		searchQuery = bleve.NewConjunctionQuery(qs...)
	}
	page := NewPage()
	page.PageNo = pageNo
	from := (pageNo - 1) * page.PageSize
	if from < 0 {
		from = 0
	}
	searchRequest := bleve.NewSearchRequestOptions(searchQuery, page.PageSize, from, false)
	// 指定返回的字段
	searchRequest.Fields = []string{"*"}

	// 先将按"sortNo"字段对结果进行排序.如果两个文档在此字段中具有相同的值,则它们将按得分(_score)降序排序,如果文档具有相同的SortNo和得分,则将按文档ID(_id)降序排序.
	searchRequest.SortBy([]string{"sortNo", "-_score", "-_id"})

	searchResult, err := searchIndex.SearchInContext(ctx, searchRequest)
	if err != nil {
		return ResponseData{StatusCode: 0, ERR: err}, err
	}
	total, err := strconv.Atoi(strconv.FormatUint(searchResult.Total, 10))
	if err != nil {
		return ResponseData{StatusCode: 0, ERR: err}, err
	}
	page.setTotalCount(total)
	data, err := result2SliceMap(searchIndex.Name(), searchResult)
	if err != nil {
		return ResponseData{StatusCode: 0, ERR: err}, err
	}

	return ResponseData{StatusCode: 1, Data: data, Page: page}, err
}

func findIndexOne(ctx context.Context, c *app.RequestContext, indexName string, id string) (ResponseData, error) {
	searchIndex, ok, _ := openBleveIndex(indexName)
	if !ok { //索引不存在
		err := errors.New("索引不存在")
		return ResponseData{StatusCode: 0, ERR: err}, err
	}
	idQuery := bleveNewTermQuery(id)
	// 指定查询的字段
	idQuery.SetField("id")
	searchRequest := bleve.NewSearchRequest(idQuery)
	// 指定返回的字段
	searchRequest.Fields = []string{"*"}
	// 先将按"sortNo"字段对结果进行排序.如果两个文档在此字段中具有相同的值,则它们将按得分(_score)降序排序,如果文档具有相同的SortNo和得分,则将按文档ID(_id)降序排序.
	searchRequest.SortBy([]string{"sortNo", "-_score", "-_id"})
	searchResult, err := searchIndex.SearchInContext(ctx, searchRequest)
	if err != nil {
		return ResponseData{StatusCode: 0, ERR: err}, err
	}
	data, err := result2Map(indexName, searchResult)
	if err != nil {
		return ResponseData{StatusCode: 0, ERR: err}, err
	}

	return ResponseData{StatusCode: 1, Data: data}, err
}

func bleveNew(indexName string, mapping mapping.IndexMapping) (bleve.Index, error) {
	index, err := bleve.New(bleveDataDir+indexName, mapping)
	if err != nil {
		FuncLogError(err)
		return nil, err
	}
	IndexMap[indexName] = index
	return index, err
}

// openBleveIndex 打开索引目录
func openBleveIndex(indexName string) (bleve.Index, bool, error) {
	if !pathExists(bleveDataDir + indexName) { //如果索文件不存在
		return nil, false, nil
	}
	index, ok := IndexMap[indexName]
	if ok { //已经打开过
		return index, true, nil
	}
	// 打开所有的索引,放到map里,一个索引只能打开一次.
	index, err := bleve.Open(bleveDataDir + indexName)
	if err != nil {
		FuncLogError(err)
		return nil, false, err
	}
	IndexMap[indexName] = index
	return index, true, nil
}
func bleveNewTermQuery(term string) *query.TermQuery {
	term = strings.ToLower(strings.TrimSpace(term))
	return bleve.NewTermQuery(term)
}
