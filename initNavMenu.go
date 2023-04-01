package main

import (
	"time"

	"github.com/blevesearch/bleve/v2"
)

// 导航菜单
func init() {
	ok, err := openBleveIndex(indexNavMenuName)
	if err != nil || ok {
		return
	}

	// 获取当前时间
	now := time.Now()
	// 创建用户表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultAnalyzer = gseAnalyzerName
	// 初始化各个字段
	navMenuId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexNavMenuName,
		FieldCode:    "id",
		FieldName:    "导航菜单ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       1,
		Active:       3,
		Required:     1,
	}
	// 放入文件中
	addIndexField(mapping, navMenuId)

	navMenuMenuName := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexNavMenuName,
		FieldCode:    "menuName",
		FieldName:    "菜单名称",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       2,
		Active:       3,
		Required:     1,
	}
	addIndexField(mapping, navMenuMenuName)

	navMenuHrefURL := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexNavMenuName,
		FieldCode:    "hrefURL",
		FieldName:    "跳转路径",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       3,
		Active:       3,
		Required:     0,
	}
	addIndexField(mapping, navMenuHrefURL)

	navMenuHrefTarget := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexNavMenuName,
		FieldCode:    "hrefTarget",
		FieldName:    "跳转方式",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       4,
		Active:       3,
		Required:     0,
	}
	addIndexField(mapping, navMenuHrefTarget)

	navMenuPID := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexNavMenuName,
		FieldCode:    "pid",
		FieldName:    "父菜单ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       5,
		Active:       3,
		Required:     0,
	}
	addIndexField(mapping, navMenuPID)

	navMenuThemePC := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexNavMenuName,
		FieldCode:    "themePC",
		FieldName:    "PC主题",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       6,
		Active:       3,
		Required:     0,
	}
	addIndexField(mapping, navMenuThemePC)

	navMenuModuleID := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexNavMenuName,
		FieldCode:    "moduleID",
		FieldName:    "模型ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       7,
		Active:       3,
		Required:     0,
	}
	addIndexField(mapping, navMenuModuleID)

	navMenuComCode := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexNavMenuName,
		FieldCode:    "comCode",
		FieldName:    "逗号隔开的全路径",
		FieldType:    fieldType_文本框,
		AnalyzerName: commaAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       8,
		Active:       3,
		Required:     0,
	}
	addIndexField(mapping, navMenuComCode)

	navMenuTemplateID := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexNavMenuName,
		FieldCode:    "templateID",
		FieldName:    "模板Id",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       9,
		Active:       3,
		Required:     0,
	}
	addIndexField(mapping, navMenuTemplateID)

	navMenuChildTemplateID := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexNavMenuName,
		FieldCode:    "childTemplateID",
		FieldName:    "子页面模板Id",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       10,
		Active:       3,
		Required:     0,
	}
	addIndexField(mapping, navMenuChildTemplateID)

	// 添加公共字段
	indexCommonField(mapping, indexNavMenuName, 10, now)

	// //mapping.DefaultMapping.AddFieldMappingsAt("*", keywordMapping)
	navMenuIndex, err := bleve.New(indexNavMenuName, mapping)

	// 放到IndexMap中
	IndexMap[indexNavMenuName] = navMenuIndex

	if err != nil {
		FuncLogError(err)
		return
	}
	//保存表信息
	indexInfo := IndexMap[indexInfoName]
	indexInfo.Index(indexNavMenuName, IndexInfoStruct{
		ID:         indexNavMenuName,
		Name:       "导航菜单",
		Code:       "navMenu",
		IndexType:  "index",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     6,
		Active:     1,
	})
}
