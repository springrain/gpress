package main

import (
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
)

// initConfig 初始化创建Config索引
func initConfig() (bool, error) {
	// 获取索引字段的表
	indexField := IndexMap[indexFieldIndexName]
	// 当前时间
	now := time.Now()

	// 用户表的 ID 字段
	configId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    configIndexName,
		IndexName:    "配置信息",
		FieldCode:    "id",
		FieldName:    "配置ID",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       1,
		Active:       3,
	}
	indexField.Index(configId.ID, configId)

	// 配置 basePath 字段
	basePath := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    configIndexName,
		IndexName:    "配置信息",
		FieldCode:    "basePath",
		FieldName:    "基础路径",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       2,
		Active:       1,
	}
	indexField.Index(basePath.ID, basePath)

	// 配置 jwtSecret 字段
	jwtSecret := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    configIndexName,
		IndexName:    "配置信息",
		FieldCode:    "jwtSecret",
		FieldName:    "jwt加密字符串",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       4,
		Active:       1,
	}
	indexField.Index(jwtSecret.ID, jwtSecret)

	// 配置 jwtSecret 字段
	jwttokenKey := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    configIndexName,
		IndexName:    "配置信息",
		FieldCode:    "jwttokenKey",
		FieldName:    "jwt的key",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       5,
		Active:       1,
	}
	indexField.Index(jwttokenKey.ID, jwttokenKey)

	// 配置 serverPort 字段
	serverPort := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    configIndexName,
		IndexName:    "配置信息",
		FieldCode:    "serverPort",
		FieldName:    "服务器ip:端口",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       6,
		Active:       1,
	}
	indexField.Index(serverPort.ID, serverPort)

	// 配置 theme 字段
	theme := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    configIndexName,
		IndexName:    "配置信息",
		FieldCode:    "theme",
		FieldName:    "主题",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       7,
		Active:       1,
	}
	indexField.Index(theme.ID, theme)

	// 配置 timeout 字段
	timeout := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    configIndexName,
		IndexName:    "配置信息",
		FieldCode:    "timeout",
		FieldName:    "超时时间",
		FieldType:    1,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       8,
		Active:       1,
	}
	indexField.Index(timeout.ID, timeout)

	// 创建配置表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name
	configIndex, err := bleve.New(configIndexName, mapping)
	if err != nil {
		FuncLogError(err)
		return false, err
	}
	IndexMap[configIndexName] = configIndex
	return true, nil
}
