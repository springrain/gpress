package main

import (
	"time"
)

func init() {
	if tableExist(tableContentName) {
		return
	}

	// 创建内容表的表
	createTableSQL := `CREATE TABLE content (
		id TEXT PRIMARY KEY     NOT NULL,
		moduleID         TEXT  ,
		title         TEXT   NOT NULL,
		keyword           TEXT,
		description           TEXT,
		pageURL           TEXT,
		subtitle           TEXT,
		navMenuID           TEXT,
		navMenuName           TEXT,
		templateID           TEXT,
		author           TEXT,
		tag           TEXT,
		toc           TEXT,
		summary           TEXT,
		content           TEXT,
		markdown           TEXT,
		thumbnail           TEXT,
		createTime        TEXT,
		updateTime        TEXT,
		createUser        TEXT,
		sortNo            int NOT NULL,
		status            int NOT NULL
	 ) strict ;`
	_, err := crateTable(createTableSQL)
	if err != nil {
		return
	}

	// 获取当前时间
	now := time.Now()
	sortNo := 1
	// 初始化各个字段
	contentId := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "id",
		FieldName:    "文章内容ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	// 放入文件中
	sortNo++
	saveTableField(contentId)

	contentModuleID := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "moduleID",
		FieldName:    "模型的Code",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(contentModuleID)

	contentTitle := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableContentName,
		FieldCode: "title",
		FieldName: "标题",
		FieldType: fieldType_文本框,
		// 文章标题使用中文分词
		AnalyzerName: gseAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(contentTitle)
	// title 字段使用 中文分词器的mapping gseAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("title", gseAnalyzerMapping)

	contentKeyword := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableContentName,
		FieldCode: "keyword",
		FieldName: "关键字",
		FieldType: fieldType_文本框,
		// 文章关键字使用逗号分词器
		AnalyzerName: commaAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(contentKeyword)
	// keyword 字段使用 逗号分词器的mapping commaAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("keyword", commaAnalyzerMapping)

	contentDescription := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableContentName,
		FieldCode: "description",
		FieldName: "站点描述",
		FieldType: fieldType_文本框,
		// 文章描述使用中文分词器
		AnalyzerName: gseAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(contentDescription)
	// description 字段使用 中文分词器的mapping gseAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("description", gseAnalyzerMapping)

	contentPageURL := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "pageURL",
		FieldName:    "自身页面路径",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(contentPageURL)

	contentSubtitle := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableContentName,
		FieldCode: "subtitle",
		FieldName: "副标题",
		FieldType: fieldType_文本框,
		// 文章副标题使用中文分词器
		AnalyzerName: gseAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(contentSubtitle)
	// subtitle 字段使用 中文分词器的mapping gseAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("subtitle", gseAnalyzerMapping)

	contentNavMenuId := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "navMenuID",
		FieldName:    "导航ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: commaAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(contentNavMenuId)
	// navMenuId 字段使用 逗号分词器的mapping commaAnalyzerMapping
	// //mapping.DefaultMapping.AddFieldMappingsAt("navMenuId", commaAnalyzerMapping)

	contentNavMenuNames := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "navMenuName",
		FieldName:    "导航名称,逗号(,)隔开",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(contentNavMenuNames)
	// navMenuNames 字段使用 中文分词器的mapping gseAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("navMenuNames", gseAnalyzerMapping)

	contentTemplateID := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "templateID",
		FieldName:    "模板Id",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(contentTemplateID)

	contentAuthor := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "author",
		FieldName:    "作者",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(contentAuthor)

	contentTag := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "tag",
		FieldName:    "标签",
		FieldType:    fieldType_文本框,
		AnalyzerName: commaAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(contentTag)

	contentToc := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "toc",
		FieldName:    "文章目录",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(contentToc)

	contentSummary := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "summary",
		FieldName:    "文章摘要",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(contentSummary)

	contentContent := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "content",
		FieldName:    "文章内容",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(contentContent)
	// content 字段使用 中文分词器的mapping gseAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("content", gseAnalyzerMapping)

	contentMarkdown := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "markdown",
		FieldName:    "markdown原文",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(contentMarkdown)

	contentThumbnail := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "thumbnail",
		FieldName:    "封面图",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(contentThumbnail)

	// 添加公共字段
	indexCommonField(tableContentName, "文章内容", sortNo, now)

	//保存表信息
	//tableInfo, _, _ := openBleveTable(tableInfoName)
	saveTableInfo(TableInfoStruct{
		ID:         tableContentName,
		Name:       "文章内容",
		Code:       "content",
		TableType:  "table",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     1,
		Status:     1,
	})

	//保存默认模型表信息
	saveTableInfo(TableInfoStruct{
		ID:         tableModuleDefaultName,
		Name:       "默认模型",
		Code:       tableModuleDefaultName,
		TableType:  "module",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     4,
		Status:     1,
	})
}
