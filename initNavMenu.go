package main

import (
	"context"
	"time"
)

// 导航菜单
func init() {
	if tableExist(tableNavMenuName) {
		return
	}

	// 获取当前时间
	now := time.Now().Format("2006-01-02 15:04:05")
	// 创建用户表的表
	createTableSQL := `CREATE TABLE IF NOT EXISTS navMenu (
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
	ctx := context.Background()
	_, err := crateTable(ctx, createTableSQL)
	if err != nil {
		return
	}

	sortNo := 1
	// 初始化各个字段
	navMenuId := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableNavMenuName,
		FieldCode:    "id",
		FieldName:    "导航菜单ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "导航菜单",
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
	saveTableField(ctx, navMenuId)

	navMenuMenuName := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableNavMenuName,
		FieldCode:    "menuName",
		FieldName:    "菜单名称",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "导航菜单",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	sortNo++
	saveTableField(ctx, navMenuMenuName)

	navMenuHrefURL := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableNavMenuName,
		FieldCode:    "hrefURL",
		FieldName:    "跳转路径",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "导航菜单",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     0,
	}
	sortNo++
	saveTableField(ctx, navMenuHrefURL)

	navMenuHrefTarget := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableNavMenuName,
		FieldCode:    "hrefTarget",
		FieldName:    "跳转方式",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "导航菜单",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     0,
	}
	sortNo++
	saveTableField(ctx, navMenuHrefTarget)

	navMenuPID := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableNavMenuName,
		FieldCode:    "pid",
		FieldName:    "父菜单ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "导航菜单",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     0,
	}
	sortNo++
	saveTableField(ctx, navMenuPID)

	navMenuThemePC := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableNavMenuName,
		FieldCode:    "themePC",
		FieldName:    "PC主题",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "导航菜单",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     0,
	}
	sortNo++
	saveTableField(ctx, navMenuThemePC)

	navMenuModuleID := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableNavMenuName,
		FieldCode:    "moduleID",
		FieldName:    "模型ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "导航菜单",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     0,
	}
	sortNo++
	saveTableField(ctx, navMenuModuleID)

	navMenuComCode := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableNavMenuName,
		FieldCode:    "comCode",
		FieldName:    "逗号隔开的全路径",
		FieldType:    fieldType_文本框,
		AnalyzerName: commaAnalyzerName,
		TableName:    "导航菜单",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     0,
	}
	sortNo++
	saveTableField(ctx, navMenuComCode)

	navMenuTemplateID := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableNavMenuName,
		FieldCode:    "templateID",
		FieldName:    "模板Id",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "导航菜单",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     0,
	}
	sortNo++
	saveTableField(ctx, navMenuTemplateID)

	navMenuChildTemplateID := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableNavMenuName,
		FieldCode:    "childTemplateID",
		FieldName:    "子页面模板Id",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "导航菜单",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     0,
	}
	sortNo++
	saveTableField(ctx, navMenuChildTemplateID)

	// 添加公共字段
	indexCommonField(ctx, tableNavMenuName, "导航菜单", sortNo, now)

	//保存表信息
	saveTableInfo(ctx, TableInfoStruct{
		ID:         tableNavMenuName,
		Name:       "导航菜单",
		Code:       "navMenu",
		TableType:  "table",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     5,
		Status:     1,
	})
}
