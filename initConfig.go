package main

import (
	"time"
)

// initConfig 初始化创建Config表
func initConfig() (bool, error) {
	if tableExist(tableConfigName) {
		return true, nil
	}
	// 获取表字段的表
	//tableField ,_:= openBleveTable(indexFieldName]
	// 当前时间
	now := time.Now()
	// 创建配置表的表
	createTableSQL := `CREATE TABLE config (
		id TEXT PRIMARY KEY     NOT NULL,
		basePath         TEXT  NOT NULL,
		jwtSecret        TEXT   NOT NULL,
		jwttokenKey      TEXT NOT NULL,
		serverPort       TEXT NOT NULL,
		theme            TEXT NOT NULL,
		timeout          INT NOT NULL,
		createTime       TEXT,
		updateTime       TEXT,
		createUser       TEXT,
		sortNo           int NOT NULL,
		status           int NOT NULL
	 ) strict ;`
	_, err := crateTable(createTableSQL)
	if err != nil {
		return false, err
	}
	sortNo := 1
	// ID 字段
	configId := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableConfigName,
		FieldCode:    "id",
		FieldName:    "配置ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "配置信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}

	sortNo++
	saveTableField(configId)

	// 配置 basePath 字段
	basePath := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableConfigName,
		FieldCode:    "basePath",
		FieldName:    "基础路径",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "配置信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	saveTableField(basePath)

	// 配置 jwtSecret 字段
	jwtSecret := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableConfigName,
		FieldCode:    "jwtSecret",
		FieldName:    "jwt加密字符串",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "配置信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	saveTableField(jwtSecret)

	// 配置 jwtSecret 字段
	jwttokenKey := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableConfigName,
		FieldCode:    "jwttokenKey",
		FieldName:    "jwt的key",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "配置信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	saveTableField(jwttokenKey)

	// 配置 serverPort 字段
	serverPort := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableConfigName,
		FieldCode:    "serverPort",
		FieldName:    "服务器ip:端口",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "配置信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	saveTableField(serverPort)

	// 配置 theme 字段
	theme := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableConfigName,
		FieldCode:    "theme",
		FieldName:    "主题",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "配置信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	saveTableField(theme)

	// 配置 timeout 字段
	timeout := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableConfigName,
		FieldCode:    "timeout",
		FieldName:    "超时时间",
		FieldType:    fieldType_数字,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "配置信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	saveTableField(timeout)

	// 添加公共字段
	indexCommonField(tableConfigName, "配置信息", sortNo, now)

	return true, nil
}

func init() {
	// 获取当前时间
	now := time.Now()

	//保存表信息
	//tableInfo, _, _ := openBleveTable(tableInfoName)
	saveTableInfo(TableInfoStruct{
		ID:         tableConfigName,
		Name:       "配置信息",
		Code:       "config",
		TableType:  "table",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     2,
		Status:     1,
	})

}
