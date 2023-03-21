package main

import (
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
)

// 初始化站点信息
func initSite() (bool, error) {
	indexField := IndexMap[indexFieldIndexName]

	// 创建用户表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name
	// 获取当前时间
	now := time.Now()

	// 初始化各个字段
	siteId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteIndexName,
		IndexName:    "站点信息",
		FieldCode:    "ID",
		FieldName:    "站点信息ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       1,
		Active:       3,
	}
	// 放入文件中
	indexField.Index(siteId.ID, siteId)

	siteTitle := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteIndexName,
		IndexName:    "站点信息",
		FieldCode:    "Title",
		FieldName:    "标题",
		FieldType:    fieldType_文本框,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       2,
		Active:       3,
	}
	indexField.Index(siteTitle.ID, siteTitle)

	siteKeyWords := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteIndexName,
		IndexName:    "站点信息",
		FieldCode:    "KeyWords",
		FieldName:    "关键字",
		FieldType:    fieldType_文本框,
		AnalyzerName: commaAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       3,
		Active:       3,
	}
	indexField.Index(siteKeyWords.ID, siteKeyWords)

	// KeyWords 字段使用 逗号分词器的mapping commaAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("KeyWords", commaAnalyzerMapping)

	siteDescription := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteIndexName,
		IndexName:    "站点信息",
		FieldCode:    "Description",
		FieldName:    "站点描述",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       4,
		Active:       3,
	}
	indexField.Index(siteDescription.ID, siteDescription)

	// Description 字段使用 中文分词器的mapping gseAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("Description", gseAnalyzerMapping)

	siteTheme := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteIndexName,
		IndexName:    "站点信息",
		FieldCode:    "theme",
		FieldName:    "默认主题",
		FieldType:    fieldType_文本框,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       5,
		Active:       3,
	}
	indexField.Index(siteTheme.ID, siteTheme)

	siteThemePC := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteIndexName,
		IndexName:    "站点信息",
		FieldCode:    "themePC",
		FieldName:    "PC主题",
		FieldType:    fieldType_文本框,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       6,
		Active:       3,
	}
	indexField.Index(siteThemePC.ID, siteThemePC)

	siteThemeWAP := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteIndexName,
		IndexName:    "站点信息",
		FieldCode:    "themeWAP",
		FieldName:    "手机主题",
		FieldType:    fieldType_文本框,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       7,
		Active:       3,
	}
	indexField.Index(siteThemeWAP.ID, siteThemeWAP)

	siteThemeWEIXIN := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteIndexName,
		IndexName:    "站点信息",
		FieldCode:    "siteThemeWEIXIN",
		FieldName:    "微信主题",
		FieldType:    fieldType_文本框,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       8,
		Active:       3,
	}
	indexField.Index(siteThemeWEIXIN.ID, siteThemeWEIXIN)

	siteLogo := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteIndexName,
		IndexName:    "站点信息",
		FieldCode:    "Logo",
		FieldName:    "Logo",
		FieldType:    fieldType_文本框,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       9,
		Active:       3,
	}
	indexField.Index(siteLogo.ID, siteLogo)

	siteFavicon := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteIndexName,
		IndexName:    "站点信息",
		FieldCode:    "Favicon",
		FieldName:    "Favicon",
		FieldType:    fieldType_文本框,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       10,
		Active:       3,
	}
	indexField.Index(siteFavicon.ID, siteFavicon)

	siteIndexIndex, err := bleve.New(indexSiteIndexName, mapping)
	// 放到IndexMap中
	IndexMap[indexSiteIndexName] = siteIndexIndex

	if err != nil {
		FuncLogError(err)
		return false, err
	}

	return true, nil
}
