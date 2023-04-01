package main

import (
	"time"

	"github.com/blevesearch/bleve/v2"
)

func init() {
	ok, err := openBleveIndex(indexPageTemplateName)
	if err != nil || ok {
		return
	}

	// 获取当前时间
	now := time.Now()

	// 创建用户表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultAnalyzer = gseAnalyzerName
	// 初始化各个字段
	pageId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexPageTemplateName,
		FieldCode:    "id",
		FieldName:    "页面模板id",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       1,
		Active:       3,
	}
	// 放入文件中
	addIndexField(mapping, pageId)

	pageTemplateNameName := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexPageTemplateName,
		FieldCode:    "templateName",
		FieldName:    "模板名称",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       2,
		Active:       3,
	}
	addIndexField(mapping, pageTemplateNameName)

	pageTemplateNamePath := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexPageTemplateName,
		FieldCode:    "templatePath",
		FieldName:    "模板路径",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       3,
		Active:       3,
	}
	addIndexField(mapping, pageTemplateNamePath)

	// 添加公共字段
	indexCommonField(mapping, indexInfoName, 3, now)

	// //mapping.DefaultMapping.AddFieldMappingsAt("*", keywordMapping)
	pageTemplateIndex, err := bleve.New(indexPageTemplateName, mapping)

	// 放到IndexMap中
	IndexMap[indexPageTemplateName] = pageTemplateIndex

	if err != nil {
		FuncLogError(err)
		return
	}
	//保存表信息
	indexInfo := IndexMap[indexInfoName]
	indexInfo.Index(indexPageTemplateName, IndexInfoStruct{
		ID:         indexPageTemplateName,
		Name:       "页面模板",
		IndexType:  "index",
		Code:       "pageTemplate",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     7,
		Active:     1,
	})
}
