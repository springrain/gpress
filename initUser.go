package main

import (
	"time"
)

// initUser 初始化创建User表
func init() {
	if tableExist(tableUserName) {
		return
	}

	// 当前时间
	now := time.Now()
	sortNo := 1
	// 创建用户表的表
	createTableSQL := `CREATE TABLE user (
		id TEXT PRIMARY KEY     NOT NULL,
		account         TEXT  NOT NULL,
		password         TEXT   NOT NULL,
		userName         TEXT NOT NULL
	 ) strict ;`
	_, err := crateTable(createTableSQL)
	if err != nil {
		return
	}

	// 用户表的 ID 字段
	userId := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableUserName,
		FieldCode:    "id",
		FieldName:    "用户ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "用户信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(userId)

	// 用户表的 Account 字段
	userAccount := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableUserName,
		FieldCode:    "account",
		FieldName:    "账号",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "用户信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	saveTableField(userAccount)
	// 用户表的 PassWord 字段
	userPassWord := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableUserName,
		FieldCode:    "password",
		FieldName:    "密码",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "用户信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	saveTableField(userPassWord)
	// 用户表的 UserName 字段
	userName := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableUserName,
		FieldCode:    "userName",
		FieldName:    "用户名称",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "用户信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	saveTableField(userName)

	//保存表信息
	saveTableInfo(TableInfoStruct{
		ID:         tableUserName,
		Name:       "用户信息",
		TableType:  "table",
		Code:       "user",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     8,
		Status:     1,
	})
}
