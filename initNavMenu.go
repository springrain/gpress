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
	indexField := IndexMap[indexFieldName]

	// 获取当前时间
	now := time.Now()

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
	indexField.Index(navMenuId.ID, navMenuId)

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
	indexField.Index(navMenuMenuName.ID, navMenuMenuName)

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
	indexField.Index(navMenuHrefURL.ID, navMenuHrefURL)

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
	indexField.Index(navMenuHrefTarget.ID, navMenuHrefTarget)

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
	indexField.Index(navMenuPID.ID, navMenuPID)

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
	indexField.Index(navMenuThemePC.ID, navMenuThemePC)

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
	indexField.Index(navMenuModuleID.ID, navMenuModuleID)

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
	indexField.Index(navMenuComCode.ID, navMenuComCode)

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
	indexField.Index(navMenuTemplateID.ID, navMenuTemplateID)

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
	indexField.Index(navMenuChildTemplateID.ID, navMenuChildTemplateID)

	// 添加公共字段
	indexCommonField(indexField, indexNavMenuName, 10, now)

	// 创建用户表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器,存在问题:NewQueryStringQuery时不能正确匹配查询
	mapping.DefaultMapping.DefaultAnalyzer = keywordAnalyzerName
	// mapping.DefaultMapping.AddFieldMappingsAt("*", keywordMapping)
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
