package main

import (
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
)

// 初始化站点信息
func initSitenInfo() (bool, error) {
	indexField := IndexMap[indexFieldIndexName]

	// 创建用户表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name
	// 获取当前时间
	now := time.Now()

	// 初始化各个字段
	sitenInfoId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSitenIndexName,
		IndexName:    "站点信息",
		FieldCode:    "ID",
		FieldName:    "站点信息ID",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       1,
		Active:       3,
	}
	// 放入文件中
	indexField.Index(sitenInfoId.ID, sitenInfoId)

	sitenInfoTitle := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSitenIndexName,
		IndexName:    "站点信息",
		FieldCode:    "Title",
		FieldName:    "标题",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       2,
		Active:       3,
	}
	indexField.Index(sitenInfoTitle.ID, sitenInfoTitle)

	sitenInfoKeyWords := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSitenIndexName,
		IndexName:    "站点信息",
		FieldCode:    "KeyWords",
		FieldName:    "关键字",
		FieldType:    3,
		AnalyzerName: commaAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       3,
		Active:       3,
	}
	indexField.Index(sitenInfoKeyWords.ID, sitenInfoKeyWords)

	// KeyWords 字段使用 逗号分词器的mapping commaAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("KeyWords", commaAnalyzerMapping)

	sitenInfoDescription := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSitenIndexName,
		IndexName:    "站点信息",
		FieldCode:    "Description",
		FieldName:    "站点描述",
		FieldType:    3,
		AnalyzerName: gseAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       4,
		Active:       3,
	}
	indexField.Index(sitenInfoDescription.ID, sitenInfoDescription)

	// Description 字段使用 中文分词器的mapping gseAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("Description", gseAnalyzerMapping)

	sitenInfoTheme := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSitenIndexName,
		IndexName:    "站点信息",
		FieldCode:    "theme",
		FieldName:    "默认主题",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       5,
		Active:       3,
	}
	indexField.Index(sitenInfoTheme.ID, sitenInfoTheme)

	sitenInfoThemePC := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSitenIndexName,
		IndexName:    "站点信息",
		FieldCode:    "themePC",
		FieldName:    "PC主题",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       6,
		Active:       3,
	}
	indexField.Index(sitenInfoThemePC.ID, sitenInfoThemePC)

	sitenInfoThemeWAP := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSitenIndexName,
		IndexName:    "站点信息",
		FieldCode:    "themeWAP",
		FieldName:    "手机主题",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       7,
		Active:       3,
	}
	indexField.Index(sitenInfoThemeWAP.ID, sitenInfoThemeWAP)

	sitenInfoThemeWEIXIN := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSitenIndexName,
		IndexName:    "站点信息",
		FieldCode:    "sitenInfoThemeWEIXIN",
		FieldName:    "微信主题",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       8,
		Active:       3,
	}
	indexField.Index(sitenInfoThemeWEIXIN.ID, sitenInfoThemeWEIXIN)

	sitenInfoLogo := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSitenIndexName,
		IndexName:    "站点信息",
		FieldCode:    "Logo",
		FieldName:    "Logo",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       9,
		Active:       3,
	}
	indexField.Index(sitenInfoLogo.ID, sitenInfoLogo)

	sitenInfoFavicon := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSitenIndexName,
		IndexName:    "站点信息",
		FieldCode:    "Favicon",
		FieldName:    "Favicon",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       10,
		Active:       3,
	}
	indexField.Index(sitenInfoFavicon.ID, sitenInfoFavicon)

	sitenIndexIndex, err := bleve.New(indexSitenIndexName, mapping)
	// 放到IndexMap中
	IndexMap[indexSitenIndexName] = sitenIndexIndex

	if err != nil {
		FuncLogError(err)
		return false, err
	}

	return true, nil
}
