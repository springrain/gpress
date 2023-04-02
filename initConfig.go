package main

import (
	"time"

	"github.com/blevesearch/bleve/v2"
)

// initConfig 初始化创建Config索引
func initConfig() (bool, error) {

	_, ok, err := openBleveIndex(indexConfigName)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
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
		SortNo:       1,
		Active:       3,
	}
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
		SortNo:       2,
		Active:       1,
	}
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
		SortNo:       4,
		Active:       1,
	}
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
		SortNo:       5,
		Active:       1,
	}
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
		SortNo:       6,
		Active:       1,
	}
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
		SortNo:       7,
		Active:       1,
	}
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
		SortNo:       8,
		Active:       1,
	}
	addIndexField(mapping, timeout)

	_, err = bleveNew(indexConfigName, mapping)
	if err != nil {
		return false, err
	}

	//保存表信息
	indexInfo, _, _ := openBleveIndex(indexInfoName)
	indexInfo.Index(indexConfigName, IndexInfoStruct{
		ID:         indexConfigName,
		Name:       "配置信息",
		Code:       "config",
		IndexType:  "index",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     1,
		Active:     1,
	})

	return true, nil
}
