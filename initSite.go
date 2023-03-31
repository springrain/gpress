package main

import (
	"time"

	"github.com/blevesearch/bleve/v2"
)

// 初始化站点信息
func init() {
	ok, err := openBleveIndex(indexSiteName)
	if err != nil || ok {
		return
	}
	indexField := IndexMap[indexFieldName]

	// 创建用户表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keywordAnalyzerName
	// 获取当前时间
	now := time.Now()

	// 初始化各个字段
	siteId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "id",
		FieldName:    "站点信息ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       1,
		Active:       3,
	}
	// 放入文件中
	indexField.Index(siteId.ID, siteId)

	siteTitle := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "title",
		FieldName:    "标题",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       2,
		Active:       3,
	}
	indexField.Index(siteTitle.ID, siteTitle)
	// title 字段使用 中文分词器的mapping gseAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("title", gseAnalyzerMapping)

	siteKeyword := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "keyword",
		FieldName:    "关键字",
		FieldType:    fieldType_文本框,
		AnalyzerName: commaAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       3,
		Active:       3,
	}
	indexField.Index(siteKeyword.ID, siteKeyword)

	// keyword 字段使用 逗号分词器的mapping commaAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("keyword", commaAnalyzerMapping)

	siteDescription := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "description",
		FieldName:    "站点描述",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       4,
		Active:       3,
	}
	indexField.Index(siteDescription.ID, siteDescription)

	// description 字段使用 中文分词器的mapping gseAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("description", gseAnalyzerMapping)

	siteTheme := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "theme",
		FieldName:    "默认主题",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       5,
		Active:       3,
	}
	indexField.Index(siteTheme.ID, siteTheme)

	siteThemePC := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "themePC",
		FieldName:    "PC主题",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       6,
		Active:       3,
	}
	indexField.Index(siteThemePC.ID, siteThemePC)

	siteThemeWAP := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "themeWAP",
		FieldName:    "手机主题",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       7,
		Active:       3,
	}
	indexField.Index(siteThemeWAP.ID, siteThemeWAP)

	siteThemeWEIXIN := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "siteThemeWEIXIN",
		FieldName:    "微信主题",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       8,
		Active:       3,
	}
	indexField.Index(siteThemeWEIXIN.ID, siteThemeWEIXIN)

	siteLogo := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "logo",
		FieldName:    "Logo",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       9,
		Active:       3,
	}
	indexField.Index(siteLogo.ID, siteLogo)

	siteFavicon := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "favicon",
		FieldName:    "Favicon",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       10,
		Active:       3,
	}
	indexField.Index(siteFavicon.ID, siteFavicon)

	// 添加公共字段
	indexCommonField(indexField, indexSiteName, 10, now)

	siteIndexIndex, err := bleve.New(indexSiteName, mapping)
	// 放到IndexMap中
	IndexMap[indexSiteName] = siteIndexIndex

	if err != nil {
		FuncLogError(err)
		return
	}
	//保存表信息
	indexInfo := IndexMap[indexInfoName]
	indexInfo.Index(indexSiteName, IndexInfoStruct{
		ID:         indexSiteName,
		Name:       "站点信息",
		IndexType:  "index",
		Code:       "site",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     3,
		Active:     1,
	})
}
