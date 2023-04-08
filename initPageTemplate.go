package main

import (
	"time"

	"github.com/blevesearch/bleve/v2"
)

func init() {
	_, ok, err := openBleveIndex(indexPageTemplateName)
	if err != nil || ok {
		return
	}

	// 获取当前时间
	now := time.Now()

	// 创建用户表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultAnalyzer = gseAnalyzerName
	sortNo := 1
	// 初始化各个字段
	pageId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexPageTemplateName,
		FieldCode:    "id",
		FieldName:    "页面模板id",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "页面模板",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	// 放入文件中
	sortNo++
	addIndexField(mapping, pageId)

	pageTemplateNameName := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexPageTemplateName,
		FieldCode:    "templateName",
		FieldName:    "模板名称",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "页面模板",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, pageTemplateNameName)

	pageTemplateNamePath := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexPageTemplateName,
		FieldCode:    "templatePath",
		FieldName:    "模板路径",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "页面模板",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, pageTemplateNamePath)

	// 添加公共字段
	indexCommonField(mapping, indexPageTemplateName, "页面模板", sortNo, now)

	// //mapping.DefaultMapping.AddFieldMappingsAt("*", keywordMapping)
	_, err = bleveNew(indexPageTemplateName, mapping)

	if err != nil {
		return
	}
	//保存表信息
	indexInfo, _, _ := openBleveIndex(indexInfoName)
	indexInfo.Index(indexPageTemplateName, IndexInfoStruct{
		ID:         indexPageTemplateName,
		Name:       "页面模板",
		IndexType:  "index",
		Code:       "pageTemplate",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     6,
		Status:     1,
	})
}
