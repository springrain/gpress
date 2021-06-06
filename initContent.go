package main

import (
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
)

func initContent() (bool, error) {

	indexField := IndexMap[indexFieldIndexName]

	//创建内容表的索引
	mapping := bleve.NewIndexMapping()
	//指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name //这是要换成逗号分词吧

	//获取当前时间
	now := time.Now()

	//初始化各个字段
	contentId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		IndexName:    "文章内容",
		FieldCode:    "ID",
		FieldName:    "文章内容ID",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       1,
		Active:       3,
	}
	//放入文件中
	indexField.Index(contentId.ID, contentId)

	contentModuleIndexCode := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		IndexName:    "文章内容",
		FieldCode:    "ModuleIndexCode",
		FieldName:    "模型的Code",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       2,
		Active:       3,
	}
	indexField.Index(contentModuleIndexCode.ID, contentModuleIndexCode)
	contentTitle := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexContentName,
		IndexName: "文章内容",
		FieldCode: "Title",
		FieldName: "标题",
		FieldType: 3,
		//文章标题使用中文分词
		AnalyzerName: gseAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       2,
		Active:       3,
	}
	indexField.Index(contentTitle.ID, contentTitle)
	//Title 字段使用 中文分词器的mapping gseAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("Title", gseAnalyzerMapping)

	contentKeyWords := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexContentName,
		IndexName: "文章内容",
		FieldCode: "KeyWords",
		FieldName: "关键字",
		FieldType: 3,
		//文章关键字使用逗号分词器
		AnalyzerName: commaAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       3,
		Active:       3,
	}
	indexField.Index(contentKeyWords.ID, contentKeyWords)
	//KeyWords 字段使用 逗号分词器的mapping commaAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("KeyWords", commaAnalyzerMapping)

	contentDescription := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexContentName,
		IndexName: "文章内容",
		FieldCode: "Description",
		FieldName: "站点描述",
		FieldType: 3,
		//文章描述使用中文分词器
		AnalyzerName: gseAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       4,
		Active:       3,
	}
	indexField.Index(contentDescription.ID, contentDescription)
	//Description 字段使用 中文分词器的mapping gseAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("Description", gseAnalyzerMapping)

	contentPageURL := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		IndexName:    "文章内容",
		FieldCode:    "PageURL",
		FieldName:    "自身页面路径",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       5,
		Active:       3,
	}
	indexField.Index(contentPageURL.ID, contentPageURL)

	contentSubtitle := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexContentName,
		IndexName: "文章内容",
		FieldCode: "Subtitle",
		FieldName: "副标题",
		FieldType: 3,
		//文章副标题使用中文分词器
		AnalyzerName: gseAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       6,
		Active:       3,
	}
	indexField.Index(contentSubtitle.ID, contentSubtitle)
	//Subtitle 字段使用 中文分词器的mapping gseAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("Subtitle", gseAnalyzerMapping)

	contentNavMenuId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		IndexName:    "文章内容",
		FieldCode:    "NavMenuId",
		FieldName:    "导航ID,逗号(,)隔开",
		FieldType:    3,
		AnalyzerName: commaAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       4,
		Active:       3,
	}
	indexField.Index(contentNavMenuId.ID, contentNavMenuId)
	//NavMenuId 字段使用 逗号分词器的mapping commaAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("NavMenuId", commaAnalyzerMapping)

	contentNavMenuName := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		IndexName:    "文章内容",
		FieldCode:    "NavMenuName",
		FieldName:    "导航名称",
		FieldType:    3,
		AnalyzerName: gseAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       5,
		Active:       3,
	}
	indexField.Index(contentNavMenuName.ID, contentNavMenuName)
	//NavMenuName 字段使用 中文分词器的mapping gseAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("NavMenuName", gseAnalyzerMapping)

	contentTemplateID := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		IndexName:    "文章内容",
		FieldCode:    "TemplateID",
		FieldName:    "模板Id",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       6,
		Active:       3,
	}
	indexField.Index(contentTemplateID.ID, contentTemplateID)

	contentContent := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		IndexName:    "文章内容",
		FieldCode:    "Content",
		FieldName:    "文章内容",
		FieldType:    3,
		AnalyzerName: gseAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       7,
		Active:       3,
	}
	indexField.Index(contentContent.ID, contentContent)
	//Content 字段使用 中文分词器的mapping gseAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("Content", gseAnalyzerMapping)

	contentCreateTime := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		IndexName:    "文章内容",
		FieldCode:    "CreateTime",
		FieldName:    "创建时间",
		FieldType:    2,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       8,
		Active:       3,
	}
	indexField.Index(contentCreateTime.ID, contentCreateTime)

	moduleUpdateTime := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		IndexName:    "文章内容",
		FieldCode:    "UpdateTime",
		FieldName:    "更新时间",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       9,
		Active:       3,
	}
	indexField.Index(moduleUpdateTime.ID, moduleUpdateTime)

	moduleCreateUser := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		IndexName:    "文章内容",
		FieldCode:    "CreateUser",
		FieldName:    "创建人",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       10,
		Active:       3,
	}
	indexField.Index(moduleCreateUser.ID, moduleCreateUser)

	moduleSortNo := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		IndexName:    "文章内容",
		FieldCode:    "SortNo",
		FieldName:    "排序",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       9,
		Active:       3,
	}
	indexField.Index(moduleSortNo.ID, moduleSortNo)

	moduleActive := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexContentName,
		IndexName:    "文章内容",
		FieldCode:    "Active",
		FieldName:    "是否有效",
		FieldType:    1,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       10,
		Active:       3,
	}
	indexField.Index(moduleActive.ID, moduleActive)

	contentIndex, err := bleve.New(indexContentName, mapping)
	//放到IndexMap中
	IndexMap[indexContentName] = contentIndex

	if err != nil {
		FuncLogError(err)
		return false, err
	}
	return true, nil
}
