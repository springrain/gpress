package main

import (
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
)

func initpageTemplateName() (bool, error) {

	indexField := IndexMap[indexFieldIndexName]

	//获取当前时间
	now := time.Now()

	//初始化各个字段
	pageId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    pageTemplateName,
		IndexName:    "页面模板",
		FieldCode:    "ID",
		FieldName:    "页面模板id",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       1,
		Active:       3,
	}
	//放入文件中
	indexField.Index(pageId.ID, pageId)

	pageTemplateNameName := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    sitenIndexName,
		IndexName:    "页面模板",
		FieldCode:    "TemplateName",
		FieldName:    "模板名称",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       2,
		Active:       3,
	}
	indexField.Index(pageTemplateNameName.ID, pageTemplateNameName)

	pageTemplateNamePath := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    sitenIndexName,
		IndexName:    "页面模板",
		FieldCode:    "TemplatePath",
		FieldName:    "模板路径",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       3,
		Active:       3,
	}
	indexField.Index(pageTemplateNamePath.ID, pageTemplateNamePath)

	pageSortNo := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    sitenIndexName,
		IndexName:    "页面模板",
		FieldCode:    "SortNo",
		FieldName:    "排序",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       4,
		Active:       3,
	}
	indexField.Index(pageSortNo.ID, pageSortNo)

	pageActive := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    sitenIndexName,
		IndexName:    "页面模板",
		FieldCode:    "Active",
		FieldName:    "是否有效",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       5,
		Active:       3,
	}
	indexField.Index(pageActive.ID, pageActive)

	//创建用户表的索引
	mapping := bleve.NewIndexMapping()
	//指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name
	pageTemplateIndex, err := bleve.New(pageTemplateName, mapping)

	//放到IndexMap中
	IndexMap[pageTemplateName] = pageTemplateIndex

	if err != nil {
		FuncLogError(err)
		return false, err
	}

	return true, nil

}
