package main

import (
	"time"
)

func init() {
	if pathExist(bleveDataDir + indexContentName) {
		return
	}

	// 创建内容表的索引
	mapping := bleveNewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultAnalyzer = gseAnalyzerName
	// //mapping.DefaultMapping.AddFieldMappingsAt("*", keywordMapping)

	// 获取当前时间
	now := time.Now()
	sortNo := 1
	// 初始化各个字段
	contentId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "id",
		FieldName:    "文章内容ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	// 放入文件中
	sortNo++
	addIndexField(mapping, contentId)

	contentModuleID := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "moduleID",
		FieldName:    "模型的Code",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, contentModuleID)

	contentTitle := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexContentName,
		FieldCode: "title",
		FieldName: "标题",
		FieldType: fieldType_文本框,
		// 文章标题使用中文分词
		AnalyzerName: gseAnalyzerName,
		IndexName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, contentTitle)
	// title 字段使用 中文分词器的mapping gseAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("title", gseAnalyzerMapping)

	contentKeyword := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexContentName,
		FieldCode: "keyword",
		FieldName: "关键字",
		FieldType: fieldType_文本框,
		// 文章关键字使用逗号分词器
		AnalyzerName: commaAnalyzerName,
		IndexName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, contentKeyword)
	// keyword 字段使用 逗号分词器的mapping commaAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("keyword", commaAnalyzerMapping)

	contentDescription := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexContentName,
		FieldCode: "description",
		FieldName: "站点描述",
		FieldType: fieldType_文本框,
		// 文章描述使用中文分词器
		AnalyzerName: gseAnalyzerName,
		IndexName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, contentDescription)
	// description 字段使用 中文分词器的mapping gseAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("description", gseAnalyzerMapping)

	contentPageURL := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "pageURL",
		FieldName:    "自身页面路径",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, contentPageURL)

	contentSubtitle := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexContentName,
		FieldCode: "subtitle",
		FieldName: "副标题",
		FieldType: fieldType_文本框,
		// 文章副标题使用中文分词器
		AnalyzerName: gseAnalyzerName,
		IndexName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, contentSubtitle)
	// subtitle 字段使用 中文分词器的mapping gseAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("subtitle", gseAnalyzerMapping)

	contentNavMenuId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "navMenuID",
		FieldName:    "导航ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: commaAnalyzerName,
		IndexName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, contentNavMenuId)
	// navMenuId 字段使用 逗号分词器的mapping commaAnalyzerMapping
	// //mapping.DefaultMapping.AddFieldMappingsAt("navMenuId", commaAnalyzerMapping)

	contentNavMenuNames := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "navMenuName",
		FieldName:    "导航名称,逗号(,)隔开",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		IndexName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, contentNavMenuNames)
	// navMenuNames 字段使用 中文分词器的mapping gseAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("navMenuNames", gseAnalyzerMapping)

	contentTemplateID := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "templateID",
		FieldName:    "模板Id",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, contentTemplateID)

	contentAuthor := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "author",
		FieldName:    "作者",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, contentAuthor)

	contentTag := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "tag",
		FieldName:    "标签",
		FieldType:    fieldType_文本框,
		AnalyzerName: commaAnalyzerName,
		IndexName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, contentTag)

	contentToc := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "toc",
		FieldName:    "文章目录",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, contentToc)

	contentSummary := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "summary",
		FieldName:    "文章摘要",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		IndexName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, contentSummary)

	contentContent := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "content",
		FieldName:    "文章内容",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		IndexName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, contentContent)
	// content 字段使用 中文分词器的mapping gseAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("content", gseAnalyzerMapping)

	contentMarkdown := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "markdown",
		FieldName:    "markdown原文",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		IndexName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, contentMarkdown)

	contentThumbnail := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "thumbnail",
		FieldName:    "封面图",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, contentThumbnail)

	// 添加公共字段
	indexCommonField(mapping, indexContentName, "文章内容", sortNo, now)

	ok, err := bleveNew(indexContentName, mapping)
	if err != nil || !ok {
		return
	}

	//保存表信息
	//indexInfo, _, _ := openBleveIndex(indexInfoName)
	bleveSaveIndex(indexInfoName, indexContentName, IndexInfoStruct{
		ID:         indexContentName,
		Name:       "文章内容",
		Code:       "content",
		IndexType:  "index",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     1,
		Status:     1,
	})

	//保存默认模型表信息
	bleveSaveIndex(indexInfoName, indexModuleDefaultName, IndexInfoStruct{
		ID:         indexModuleDefaultName,
		Name:       "默认模型",
		Code:       indexModuleDefaultName,
		IndexType:  "module",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     4,
		Status:     1,
	})
}
