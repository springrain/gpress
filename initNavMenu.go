package main

import (
	"context"
	"time"
)

// 导航菜单
func init() {
	if tableExist(tableCategoryName) {
		return
	}

	// 获取当前时间
	now := time.Now().Format("2006-01-02 15:04:05")
	// 创建用户表的表
	createTableSQL := `CREATE TABLE IF NOT EXISTS category (
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
	categoryId := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableCategoryName,
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
	saveTableField(ctx, categoryId)

	categoryMenuName := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableCategoryName,
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
	saveTableField(ctx, categoryMenuName)

	categoryHrefURL := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableCategoryName,
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
	saveTableField(ctx, categoryHrefURL)

	categoryHrefTarget := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableCategoryName,
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
	saveTableField(ctx, categoryHrefTarget)

	categoryPID := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableCategoryName,
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
	saveTableField(ctx, categoryPID)

	categoryThemePC := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableCategoryName,
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
	saveTableField(ctx, categoryThemePC)

	categoryModuleID := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableCategoryName,
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
	saveTableField(ctx, categoryModuleID)

	categoryComCode := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableCategoryName,
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
	saveTableField(ctx, categoryComCode)

	categoryTemplateID := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableCategoryName,
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
	saveTableField(ctx, categoryTemplateID)

	categoryChildTemplateID := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableCategoryName,
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
	saveTableField(ctx, categoryChildTemplateID)

	// 添加公共字段
	indexCommonField(ctx, tableCategoryName, "导航菜单", sortNo, now)

	//保存表信息
	saveTableInfo(ctx, TableInfoStruct{
		ID:         tableCategoryName,
		Name:       "导航菜单",
		Code:       "category",
		TableType:  "table",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     5,
		Status:     1,
	})
}
