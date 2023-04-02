package main

import (
	"time"
)

// 默认模型 indexInfo indexType="module". 只是记录,并不创建index,全部保存到context里,用于全局检索
func init() {
	// 创建内容表的索引
	//mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	//mapping.DefaultAnalyzer = keywordAnalyzerName
	// 获取当前时间
	now := time.Now()

	// 初始化各个字段
	moduleDefaultId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexModuleDefaultName,
		FieldCode:    "id",
		FieldName:    "模型数据ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "默认模型",
		FieldComment: "",
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       1,
		Active:       3,
	}
	// 放入文件中
	addIndexField(nil, moduleDefaultId)

	moduleDefaultTitle := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexModuleDefaultName,
		FieldCode: "title",
		FieldName: "标题",
		FieldType: fieldType_文本框,
		// 文章标题使用中文分词
		AnalyzerName: gseAnalyzerName,
		IndexName:    "默认模型",
		FieldComment: "",
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       2,
		Active:       3,
	}
	addIndexField(nil, moduleDefaultTitle)
	// title 字段使用 中文分词器的mapping gseAnalyzerMapping
	////mapping.DefaultMapping.AddFieldMappingsAt("title", gseAnalyzerMapping)

	moduleDefaultKeyword := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexModuleDefaultName,
		FieldCode: "keyword",
		FieldName: "关键字",
		FieldType: fieldType_文本框,
		// 文章关键字使用逗号分词器
		AnalyzerName: commaAnalyzerName,
		IndexName:    "默认模型",
		FieldComment: "",
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       3,
		Active:       3,
	}
	addIndexField(nil, moduleDefaultKeyword)
	// keyword 字段使用 逗号分词器的mapping commaAnalyzerMapping
	////mapping.DefaultMapping.AddFieldMappingsAt("keyword", commaAnalyzerMapping)

	moduleDefaultDescription := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexModuleDefaultName,
		FieldCode: "description",
		FieldName: "站点描述",
		FieldType: fieldType_文本框,
		// 文章描述使用中文分词器
		AnalyzerName: gseAnalyzerName,
		IndexName:    "默认模型",
		FieldComment: "",
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       4,
		Active:       3,
	}
	addIndexField(nil, moduleDefaultDescription)
	// description 字段使用 中文分词器的mapping gseAnalyzerMapping
	////mapping.DefaultMapping.AddFieldMappingsAt("description", gseAnalyzerMapping)

	moduleDefaultPageURL := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexModuleDefaultName,
		FieldCode:    "pageURL",
		FieldName:    "自身页面路径",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "默认模型",
		FieldComment: "",
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       5,
		Active:       3,
	}
	addIndexField(nil, moduleDefaultPageURL)

	moduleDefaultSubtitle := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexModuleDefaultName,
		FieldCode: "subtitle",
		FieldName: "副标题",
		FieldType: fieldType_文本框,
		// 文章副标题使用中文分词器
		AnalyzerName: gseAnalyzerName,
		IndexName:    "默认模型",
		FieldComment: "",
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       6,
		Active:       3,
	}
	addIndexField(nil, moduleDefaultSubtitle)
	// subtitle 字段使用 中文分词器的mapping gseAnalyzerMapping
	////mapping.DefaultMapping.AddFieldMappingsAt("subtitle", gseAnalyzerMapping)

	moduleDefaultContent := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexModuleDefaultName,
		FieldCode: "content",
		FieldName: "文章内容",
		FieldType: fieldType_文本框,
		// 文章内容使用中文分词器
		AnalyzerName: gseAnalyzerName,
		IndexName:    "默认模型",
		FieldComment: "",
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       7,
		Active:       3,
	}
	addIndexField(nil, moduleDefaultContent)
	// content 字段使用 中文分词器的mapping gseAnalyzerMapping
	////mapping.DefaultMapping.AddFieldMappingsAt("content", gseAnalyzerMapping)

	// 添加公共字段
	indexCommonField(nil, indexModuleDefaultName, "默认模型", 7, now)

	//保存表信息
	indexInfo, _, _ := openBleveIndex(indexInfoName)
	indexInfo.Index(indexModuleDefaultName, IndexInfoStruct{
		ID:         indexModuleDefaultName,
		Name:       "默认模型",
		Code:       indexModuleDefaultName,
		IndexType:  "module",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     4,
		Active:     1,
	})
}
