package main

import (
	"time"

	"github.com/blevesearch/bleve/v2"
)

// initConfig 初始化创建Config索引
func initConfig() (bool, error) {
	if pathExist(bleveDataDir + indexConfigName) {
		return false, nil
	}
	// 获取索引字段的表
	//indexField ,_:= openBleveIndex(indexFieldName]
	// 当前时间
	now := time.Now()
	// 创建配置表的索引
	mapping := bleve.NewIndexMapping()

	// 指定默认的分词器
	mapping.DefaultAnalyzer = keywordAnalyzerName
	// //mapping.DefaultMapping.AddFieldMappingsAt("*", keywordMapping)
	sortNo := 1
	// ID 字段
	configId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexConfigName,
		FieldCode:    "id",
		FieldName:    "配置ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "配置信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}

	sortNo++
	addIndexField(mapping, configId)

	// 配置 basePath 字段
	basePath := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexConfigName,
		FieldCode:    "basePath",
		FieldName:    "基础路径",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "配置信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	addIndexField(mapping, basePath)

	// 配置 jwtSecret 字段
	jwtSecret := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexConfigName,
		FieldCode:    "jwtSecret",
		FieldName:    "jwt加密字符串",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "配置信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	addIndexField(mapping, jwtSecret)

	// 配置 jwtSecret 字段
	jwttokenKey := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexConfigName,
		FieldCode:    "jwttokenKey",
		FieldName:    "jwt的key",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "配置信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	addIndexField(mapping, jwttokenKey)

	// 配置 serverPort 字段
	serverPort := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexConfigName,
		FieldCode:    "serverPort",
		FieldName:    "服务器ip:端口",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "配置信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	addIndexField(mapping, serverPort)

	// 配置 theme 字段
	theme := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexConfigName,
		FieldCode:    "theme",
		FieldName:    "主题",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "配置信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	addIndexField(mapping, theme)

	// 配置 timeout 字段
	timeout := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexConfigName,
		FieldCode:    "timeout",
		FieldName:    "超时时间",
		FieldType:    fieldType_数字,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "配置信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	addIndexField(mapping, timeout)

	// 添加公共字段
	indexCommonField(mapping, indexConfigName, "配置信息", sortNo, now)

	ok, err := bleveNew(indexConfigName, mapping)
	if err != nil || !ok {
		return false, err
	}

	return true, nil
}

func init() {

	// 获取当前时间
	now := time.Now()

	//保存表信息
	//indexInfo, _, _ := openBleveIndex(indexInfoName)
	bleveSaveIndex(indexInfoName, indexConfigName, IndexInfoStruct{
		ID:         indexConfigName,
		Name:       "配置信息",
		Code:       "config",
		IndexType:  "index",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     2,
		Status:     1,
	})

}
