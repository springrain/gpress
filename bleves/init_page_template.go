package bleves

import (
	"gitee.com/gpress/gpress/constant"
	"gitee.com/gpress/gpress/logger"
	"gitee.com/gpress/gpress/model"
	"gitee.com/gpress/gpress/util"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
)

func initpageTemplateName() (bool, error) {
	indexField := IndexMap[constant.INDEX_FIELD_INDEX_NAME]

	// 获取当前时间
	now := time.Now()

	// 初始化各个字段
	pageId := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_PAGE_TEMPLATE_NAME,
		IndexName:    "页面模板",
		FieldCode:    "ID",
		FieldName:    "页面模板id",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       1,
		Active:       3,
	}
	// 放入文件中
	err := indexField.Index(pageId.ID, pageId)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	pageTemplateNameName := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_PAGE_TEMPLATE_NAME,
		IndexName:    "页面模板",
		FieldCode:    "TemplateName",
		FieldName:    "模板名称",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       2,
		Active:       3,
	}
	err = indexField.Index(pageTemplateNameName.ID, pageTemplateNameName)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	pageTemplateNamePath := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_PAGE_TEMPLATE_NAME,
		IndexName:    "页面模板",
		FieldCode:    "TemplatePath",
		FieldName:    "模板路径",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       3,
		Active:       3,
	}
	err = indexField.Index(pageTemplateNamePath.ID, pageTemplateNamePath)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	pageSortNo := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_PAGE_TEMPLATE_NAME,
		IndexName:    "页面模板",
		FieldCode:    "SortNo",
		FieldName:    "排序",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       4,
		Active:       3,
	}
	err = indexField.Index(pageSortNo.ID, pageSortNo)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	pageActive := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_PAGE_TEMPLATE_NAME,
		IndexName:    "页面模板",
		FieldCode:    "Active",
		FieldName:    "是否有效",
		FieldType:    1,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       5,
		Active:       3,
	}
	err = indexField.Index(pageActive.ID, pageActive)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	// 创建用户表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name
	pageTemplateIndex, err := bleve.New(constant.INDEX_PAGE_TEMPLATE_NAME, mapping)

	// 放到IndexMap中
	IndexMap[constant.INDEX_PAGE_TEMPLATE_NAME] = pageTemplateIndex

	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	return true, nil
}
