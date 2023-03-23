package gbleve

import (
	"gitee.com/gpress/gpress/constant"
	"gitee.com/gpress/gpress/logger"
	"gitee.com/gpress/gpress/model"
	"gitee.com/gpress/gpress/util"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"time"
)

// initConfig 初始化创建Config索引
func initConfig() (bool, error) {
	// 获取索引字段的表
	indexField := IndexMap[constant.INDEX_FIELD_INDEX_NAME]
	// 当前时间
	now := time.Now()

	// 用户表的 ID 字段
	configId := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.CONFIG_INDEX_NAME,
		IndexName:    "配置信息",
		FieldCode:    "id",
		FieldName:    "配置ID",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       1,
		Active:       3,
	}
	err := indexField.Index(configId.ID, configId)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	// 配置表的 configKey 字段
	configKey := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.CONFIG_INDEX_NAME,
		IndexName:    "配置的configKey",
		FieldCode:    "configKey",
		FieldName:    "编码",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       2,
		Active:       1,
	}
	err = indexField.Index(configKey.ID, configKey)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	// 用户表的 configValue 字段
	configValue := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.CONFIG_INDEX_NAME,
		IndexName:    "用户信息",
		FieldCode:    "configValue",
		FieldName:    "用户名称",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       3,
		Active:       1,
	}
	err = indexField.Index(configValue.ID, configValue)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	// 创建配置表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name
	configIndex, err := bleve.New(constant.CONFIG_INDEX_NAME, mapping)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	IndexMap[constant.CONFIG_INDEX_NAME] = configIndex
	return true, nil
}
