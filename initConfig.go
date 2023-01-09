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

	// 配置表的 configKey 字段
	configKey := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    configIndexName,
		IndexName:    "配置的configKey",
		FieldCode:    "configKey",
		FieldName:    "编码",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       2,
		Active:       1,
	}
	indexField.Index(configKey.ID, configKey)
	// 用户表的 configValue 字段
	configValue := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    configIndexName,
		IndexName:    "用户信息",
		FieldCode:    "configValue",
		FieldName:    "用户名称",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       3,
		Active:       1,
	}
	indexField.Index(configValue.ID, configValue)

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
