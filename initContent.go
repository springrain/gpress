package main

import (
	"time"

	"github.com/blevesearch/bleve/v2"
)

func initContent() (bool, error) {
	indexField := IndexMap[indexFieldName]

	// 创建内容表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keywordAnalyzerName

	// 获取当前时间
	now := time.Now()

	// 初始化各个字段
	contentId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "id",
		FieldName:    "文章内容ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       1,
		Active:       3,
	}
	// 放入文件中
	indexField.Index(contentId.ID, contentId)

	contentModuleID := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "moduleID",
		FieldName:    "模型的Code",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       2,
		Active:       3,
	}
	indexField.Index(contentModuleID.ID, contentModuleID)
	contentTitle := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexContentName,
		FieldCode: "title",
		FieldName: "标题",
		FieldType: fieldType_文本框,
		// 文章标题使用中文分词
		AnalyzerName: gseAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       2,
		Active:       3,
	}
	indexField.Index(contentTitle.ID, contentTitle)
	// title 字段使用 中文分词器的mapping gseAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("title", gseAnalyzerMapping)

	contentKeyword := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexContentName,
		FieldCode: "keyword",
		FieldName: "关键字",
		FieldType: fieldType_文本框,
		// 文章关键字使用逗号分词器
		AnalyzerName: commaAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       3,
		Active:       3,
	}
	indexField.Index(contentKeyword.ID, contentKeyword)
	// keyword 字段使用 逗号分词器的mapping commaAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("keyword", commaAnalyzerMapping)

	contentDescription := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexContentName,
		FieldCode: "description",
		FieldName: "站点描述",
		FieldType: fieldType_文本框,
		// 文章描述使用中文分词器
		AnalyzerName: gseAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       4,
		Active:       3,
	}
	indexField.Index(contentDescription.ID, contentDescription)
	// description 字段使用 中文分词器的mapping gseAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("description", gseAnalyzerMapping)

	contentPageURL := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "pageURL",
		FieldName:    "自身页面路径",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       5,
		Active:       3,
	}
	indexField.Index(contentPageURL.ID, contentPageURL)

	contentSubtitle := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexContentName,
		FieldCode: "subtitle",
		FieldName: "副标题",
		FieldType: fieldType_文本框,
		// 文章副标题使用中文分词器
		AnalyzerName: gseAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       6,
		Active:       3,
	}
	indexField.Index(contentSubtitle.ID, contentSubtitle)
	// subtitle 字段使用 中文分词器的mapping gseAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("subtitle", gseAnalyzerMapping)

	contentNavMenuId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "navMenuID",
		FieldName:    "导航ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: commaAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       4,
		Active:       3,
	}
	indexField.Index(contentNavMenuId.ID, contentNavMenuId)
	// navMenuId 字段使用 逗号分词器的mapping commaAnalyzerMapping
	// mapping.DefaultMapping.AddFieldMappingsAt("navMenuId", commaAnalyzerMapping)

	contentNavMenuNames := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "navMenuNames",
		FieldName:    "导航名称,逗号(,)隔开",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       5,
		Active:       3,
	}
	indexField.Index(contentNavMenuNames.ID, contentNavMenuNames)
	// navMenuNames 字段使用 中文分词器的mapping gseAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("navMenuNames", gseAnalyzerMapping)

	contentTemplateID := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "templateID",
		FieldName:    "模板Id",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       6,
		Active:       3,
	}
	indexField.Index(contentTemplateID.ID, contentTemplateID)

	contentContent := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "content",
		FieldName:    "文章内容",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       7,
		Active:       3,
	}
	indexField.Index(contentContent.ID, contentContent)
	// content 字段使用 中文分词器的mapping gseAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("content", gseAnalyzerMapping)

	contentCreateTime := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "createTime",
		FieldName:    "创建时间",
		FieldType:    fieldType_日期,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       8,
		Active:       3,
	}
	indexField.Index(contentCreateTime.ID, contentCreateTime)

	moduleUpdateTime := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "updateTime",
		FieldName:    "更新时间",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       9,
		Active:       3,
	}
	indexField.Index(moduleUpdateTime.ID, moduleUpdateTime)

	moduleCreateUser := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "createUser",
		FieldName:    "创建人",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       10,
		Active:       3,
	}
	indexField.Index(moduleCreateUser.ID, moduleCreateUser)

	moduleSortNo := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "sortNo",
		FieldName:    "排序",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       9,
		Active:       3,
	}
	indexField.Index(moduleSortNo.ID, moduleSortNo)

	moduleActive := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		FieldCode:    "active",
		FieldName:    "是否有效",
		FieldType:    fieldType_数字,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       10,
		Active:       3,
	}
	indexField.Index(moduleActive.ID, moduleActive)

	contentIndex, err := bleve.New(indexContentName, mapping)
	// 放到IndexMap中
	IndexMap[indexContentName] = contentIndex

	if err != nil {
		FuncLogError(err)
		return false, err
	}

	//保存表信息
	indexInfo := IndexMap[indexInfoName]
	indexInfo.Index(indexContentName, IndexInfoStruct{
		ID:         indexContentName,
		Name:       "文章内容",
		Code:       "content",
		IndexType:  "index",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     5,
		Active:     1,
	})
	return true, nil
}
