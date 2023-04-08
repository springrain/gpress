package main

import (
	"time"

	"github.com/blevesearch/bleve/v2"
)

// 导航菜单
func init() {
	_, ok, err := openBleveIndex(indexNavMenuName)
	if err != nil || ok {
		return
	}

	// 获取当前时间
	now := time.Now()
	// 创建用户表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultAnalyzer = gseAnalyzerName
	sortNo := 1
	// 初始化各个字段
	navMenuId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexNavMenuName,
		FieldCode:    "id",
		FieldName:    "导航菜单ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "导航菜单",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	// 放入文件中
	sortNo++
	addIndexField(mapping, navMenuId)

	navMenuMenuName := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexNavMenuName,
		FieldCode:    "menuName",
		FieldName:    "菜单名称",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "导航菜单",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	sortNo++
	addIndexField(mapping, navMenuMenuName)

	navMenuHrefURL := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexNavMenuName,
		FieldCode:    "hrefURL",
		FieldName:    "跳转路径",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "导航菜单",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     0,
	}
	sortNo++
	addIndexField(mapping, navMenuHrefURL)

	navMenuHrefTarget := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexNavMenuName,
		FieldCode:    "hrefTarget",
		FieldName:    "跳转方式",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "导航菜单",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     0,
	}
	sortNo++
	addIndexField(mapping, navMenuHrefTarget)

	navMenuPID := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexNavMenuName,
		FieldCode:    "pid",
		FieldName:    "父菜单ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "导航菜单",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     0,
	}
	sortNo++
	addIndexField(mapping, navMenuPID)

	navMenuThemePC := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexNavMenuName,
		FieldCode:    "themePC",
		FieldName:    "PC主题",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "导航菜单",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     0,
	}
	sortNo++
	addIndexField(mapping, navMenuThemePC)

	navMenuModuleID := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexNavMenuName,
		FieldCode:    "moduleID",
		FieldName:    "模型ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "导航菜单",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     0,
	}
	sortNo++
	addIndexField(mapping, navMenuModuleID)

	navMenuComCode := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexNavMenuName,
		FieldCode:    "comCode",
		FieldName:    "逗号隔开的全路径",
		FieldType:    fieldType_文本框,
		AnalyzerName: commaAnalyzerName,
		IndexName:    "导航菜单",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     0,
	}
	sortNo++
	addIndexField(mapping, navMenuComCode)

	navMenuTemplateID := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexNavMenuName,
		FieldCode:    "templateID",
		FieldName:    "模板Id",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "导航菜单",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     0,
	}
	sortNo++
	addIndexField(mapping, navMenuTemplateID)

	navMenuChildTemplateID := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexNavMenuName,
		FieldCode:    "childTemplateID",
		FieldName:    "子页面模板Id",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "导航菜单",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     0,
	}
	sortNo++
	addIndexField(mapping, navMenuChildTemplateID)

	// 添加公共字段
	indexCommonField(mapping, indexNavMenuName, "导航菜单", sortNo, now)

	// //mapping.DefaultMapping.AddFieldMappingsAt("*", keywordMapping)
	_, err = bleveNew(indexNavMenuName, mapping)

	if err != nil {
		return
	}
	//保存表信息
	indexInfo, _, _ := openBleveIndex(indexInfoName)
	indexInfo.Index(indexNavMenuName, IndexInfoStruct{
		ID:         indexNavMenuName,
		Name:       "导航菜单",
		Code:       "navMenu",
		IndexType:  "index",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     5,
		Status:     1,
	})
}
