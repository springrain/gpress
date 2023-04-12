package main

import (
	"context"
	"time"
)

func init() {
	if tableExist(tablePageTemplateName) {
		return
	}

	// 获取当前时间
	now := time.Now().Format("2006-01-02 15:04:05")

	// 创建用户表的表
	createTableSQL := `CREATE TABLE pageTemplate (
		id TEXT PRIMARY KEY     NOT NULL,
		templateName         TEXT  NOT NULL,
		templatePath         TEXT   NOT NULL,
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
	pageId := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tablePageTemplateName,
		FieldCode:    "id",
		FieldName:    "页面模板id",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "页面模板",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	// 放入文件中
	sortNo++
	saveTableField(ctx, pageId)

	pageTemplateNameName := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tablePageTemplateName,
		FieldCode:    "templateName",
		FieldName:    "模板名称",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "页面模板",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, pageTemplateNameName)

	pageTemplateNamePath := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tablePageTemplateName,
		FieldCode:    "templatePath",
		FieldName:    "模板路径",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "页面模板",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, pageTemplateNamePath)

	// 添加公共字段
	indexCommonField(ctx, tablePageTemplateName, "页面模板", sortNo, now)

	//保存表信息
	saveTableInfo(ctx, TableInfoStruct{
		ID:         tablePageTemplateName,
		Name:       "页面模板",
		TableType:  "table",
		Code:       "pageTemplate",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     6,
		Status:     1,
	})
}
