package main

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"time"
)

func initContent() (bool, error) {

	indexField := IndexMap[contentName]

	//获取当前时间
	now := time.Now()

	//初始化各个字段
	contentId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    contentName,
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
		IndexCode:    contentName,
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

	contentHrefURL := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    contentName,
		IndexName:    "文章内容",
		FieldCode:    "HrefURL",
		FieldName:    "页面路径",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       3,
		Active:       3,
	}
	indexField.Index(contentHrefURL.ID, contentHrefURL)

	contentNavMenuId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    contentName,
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

	contentNavMenuName := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    contentName,
		IndexName:    "文章内容",
		FieldCode:    "NavMenuName",
		FieldName:    "导航名称",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       5,
		Active:       3,
	}
	indexField.Index(contentNavMenuName.ID, contentNavMenuName)

	contentTemplateID := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    contentName,
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
		IndexCode:    contentName,
		IndexName:    "文章内容",
		FieldCode:    "Content",
		FieldName:    "文章内容",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       7,
		Active:       3,
	}
	indexField.Index(contentContent.ID, contentContent)

	contentCreateTime := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    contentName,
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
		IndexCode:    contentName,
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
		IndexCode:    contentName,
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
		IndexCode:    contentName,
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
		IndexCode:    contentName,
		IndexName:    "文章内容",
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
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name //这是要换成逗号分词吧
	_, err := bleve.New(contentName, mapping)

	if err != nil {
		FuncLogError(err)
		return false, err
	}
	return true, nil
}
