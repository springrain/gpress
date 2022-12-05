package main

import (
	"time"

	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
)

func initModuleDefault() (bool, error) {
	indexField := IndexMap[indexFieldIndexName]

	// 获取当前时间
	now := time.Now()

	// 初始化各个字段
	moduleDefaultId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexModuleDefaultName,
		IndexName:    "模型数据",
		FieldCode:    "ID",
		FieldName:    "模型数据ID",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       1,
		Active:       3,
	}
	// 放入文件中
	indexField.Index(moduleDefaultId.ID, moduleDefaultId)

	moduleDefaultTitle := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexModuleDefaultName,
		IndexName: "模型数据",
		FieldCode: "Title",
		FieldName: "标题",
		FieldType: 3,
		// 文章标题使用中文分词
		AnalyzerName: gseAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       2,
		Active:       3,
	}
	indexField.Index(moduleDefaultTitle.ID, moduleDefaultTitle)

	moduleDefaultKeyWords := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexModuleDefaultName,
		IndexName: "模型数据",
		FieldCode: "KeyWords",
		FieldName: "关键字",
		FieldType: 3,
		// 文章关键字使用逗号分词器
		AnalyzerName: commaAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       3,
		Active:       3,
	}
	indexField.Index(moduleDefaultKeyWords.ID, moduleDefaultKeyWords)

	moduleDefaultDescription := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexModuleDefaultName,
		IndexName: "模型数据",
		FieldCode: "Description",
		FieldName: "站点描述",
		FieldType: 3,
		// 文章描述使用中文分词器
		AnalyzerName: gseAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       4,
		Active:       3,
	}
	indexField.Index(moduleDefaultDescription.ID, moduleDefaultDescription)

	moduleDefaultPageURL := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexModuleDefaultName,
		IndexName:    "模型数据",
		FieldCode:    "PageURL",
		FieldName:    "自身页面路径",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       5,
		Active:       3,
	}
	indexField.Index(moduleDefaultPageURL.ID, moduleDefaultPageURL)

	moduleDefaultSubtitle := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexModuleDefaultName,
		IndexName: "模型数据",
		FieldCode: "Subtitle",
		FieldName: "副标题",
		FieldType: 3,
		// 文章副标题使用中文分词器
		AnalyzerName: gseAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       6,
		Active:       3,
	}
	indexField.Index(moduleDefaultSubtitle.ID, moduleDefaultSubtitle)

	moduleDefaultContent := IndexFieldStruct{
		ID:        FuncGenerateStringID(),
		IndexCode: indexModuleDefaultName,
		IndexName: "模型数据",
		FieldCode: "Content",
		FieldName: "文章内容",
		FieldType: 3,
		// 文章内容使用中文分词器
		AnalyzerName: gseAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       7,
		Active:       3,
	}
	indexField.Index(moduleDefaultContent.ID, moduleDefaultContent)

	moduleDefaultCreateTime := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexModuleDefaultName,
		IndexName:    "模型数据",
		FieldCode:    "CreateTime",
		FieldName:    "创建时间",
		FieldType:    2,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       8,
		Active:       3,
	}
	indexField.Index(moduleDefaultCreateTime.ID, moduleDefaultCreateTime)

	moduleUpdateTime := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexModuleDefaultName,
		IndexName:    "模型数据",
		FieldCode:    "UpdateTime",
		FieldName:    "更新时间",
		FieldType:    2,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       9,
		Active:       3,
	}
	indexField.Index(moduleUpdateTime.ID, moduleUpdateTime)

	moduleCreateUser := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexModuleDefaultName,
		IndexName:    "模型数据",
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
		IndexCode:    indexModuleDefaultName,
		IndexName:    "模型数据",
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
		IndexCode:    indexModuleDefaultName,
		IndexName:    "模型数据",
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

	/*
		//在IndexField表里设置IndexCode='Module',记录所有的Module.
		//然后在IndexField中插入每个module的字段,每个module实例的ModuleCode都是不同的,使用Module_+后缀的方式命名,只是记录,并不创建index

		//创建用户表的索引
		mapping := bleve.NewIndexMapping()
		//指定默认的分词器
		mapping.DefaultMapping.DefaultAnalyzer = keyword.Name
		moduleDefaultIndex, err := bleve.New(moduleDefaultName, mapping)
		//放到IndexMap中
		IndexMap[moduleName] = moduleDefaultIndex

		if err != nil {
			FuncLogError(err)
			return false, err
		}
	*/

	return true, nil
}
