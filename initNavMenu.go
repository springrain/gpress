package main

import (
	"time"
)

// 导航菜单
func init() {
	if tableExist(indexNavMenuName) {
		return
	}

	// 获取当前时间
	now := time.Now()
	// 创建用户表的索引
	createTableSQL := `CREATE TABLE navMenu (
		id TEXT PRIMARY KEY     NOT NULL,
		menuName          TEXT  NOT NULL,
		hrefURL           TEXT,
		hrefTarget        TEXT,
		pid        TEXT,
		themePC        TEXT,
		moduleID        TEXT,
		comCode        TEXT,
		templateID        TEXT,
		childTemplateID        TEXT,
		createTime        TEXT,
		updateTime        TEXT,
		createUser        TEXT,
		sortNo            int NOT NULL,
		status            int NOT NULL
	 ) strict ;`
	_, err := bleveNewIndexMapping(createTableSQL)
	if err != nil {
		return
	}

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
	addIndexField(navMenuId)

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
	addIndexField(navMenuMenuName)

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
	addIndexField(navMenuHrefURL)

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
	addIndexField(navMenuHrefTarget)

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
	addIndexField(navMenuPID)

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
	addIndexField(navMenuThemePC)

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
	addIndexField(navMenuModuleID)

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
	addIndexField(navMenuComCode)

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
	addIndexField(navMenuTemplateID)

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
	addIndexField(navMenuChildTemplateID)

	// 添加公共字段
	indexCommonField(indexNavMenuName, "导航菜单", sortNo, now)

	//保存表信息
	bleveSaveIndex(indexInfoName, indexNavMenuName, IndexInfoStruct{
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
