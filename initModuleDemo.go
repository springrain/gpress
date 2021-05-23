package main

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"time"
)

func initModuleDemo() (bool, error) {
	indexField := IndexMap[moduleDemoName]

	//获取当前时间
	now := time.Now()

	//初始化各个字段
	moduleDemoId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    moduleDemoName,
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
	//放入文件中
	indexField.Index(moduleDemoId.ID, moduleDemoId)

	moduleDemoTitle := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    moduleDemoName,
		IndexName:    "模型数据",
		FieldCode:    "Title",
		FieldName:    "标题",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       2,
		Active:       3,
	}
	indexField.Index(moduleDemoTitle.ID, moduleDemoTitle)

	moduleDemoKeyWords := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    moduleDemoName,
		IndexName:    "模型数据",
		FieldCode:    "KeyWords",
		FieldName:    "关键字",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       3,
		Active:       3,
	}
	indexField.Index(moduleDemoKeyWords.ID, moduleDemoKeyWords)

	moduleDemoDescription := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    moduleDemoName,
		IndexName:    "模型数据",
		FieldCode:    "Description",
		FieldName:    "站点描述",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       4,
		Active:       3,
	}
	indexField.Index(moduleDemoDescription.ID, moduleDemoDescription)

	moduleDemoPageURL := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    moduleDemoName,
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
	indexField.Index(moduleDemoPageURL.ID, moduleDemoPageURL)

	moduleDemoSubtitle := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    moduleDemoName,
		IndexName:    "模型数据",
		FieldCode:    "Subtitle",
		FieldName:    "副标题",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       6,
		Active:       3,
	}
	indexField.Index(moduleDemoSubtitle.ID, moduleDemoSubtitle)

	moduleDemoContent := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    moduleDemoName,
		IndexName:    "模型数据",
		FieldCode:    "Content",
		FieldName:    "文章内容",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       7,
		Active:       3,
	}
	indexField.Index(moduleDemoContent.ID, moduleDemoContent)

	moduleDemoCreateTime := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    moduleDemoName,
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
	indexField.Index(moduleDemoCreateTime.ID, moduleDemoCreateTime)

	moduleUpdateTime := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    moduleDemoName,
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
		IndexCode:    moduleDemoName,
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
		IndexCode:    moduleDemoName,
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
		IndexCode:    moduleDemoName,
		IndexName:    "模型数据",
		FieldCode:    "Active",
		FieldName:    "是否有效",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       10,
		Active:       3,
	}
	indexField.Index(moduleActive.ID, moduleActive)

	//创建用户表的索引
	mapping := bleve.NewIndexMapping()
	//指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name
	_, err := bleve.New(moduleDemoName, mapping)

	if err != nil {
		FuncLogError(err)
		return false, err
	}
	return true, nil
}
