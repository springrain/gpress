package main

import (
	"time"
)

// initConfig 初始化创建Config索引
func initConfig() (bool, error) {
	if tableExist(indexConfigName) {
		return true, nil
	}
	// 获取索引字段的表
	//indexField ,_:= openBleveIndex(indexFieldName]
	// 当前时间
	now := time.Now()
	// 创建配置表的索引
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
	_, err := bleveNewIndexMapping(createTableSQL)
	if err != nil {
		return false, err
	}
	sortNo := 1
	// ID 字段
	configId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexConfigName,
		FieldCode:    "id",
		FieldName:    "配置ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "配置信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}

	sortNo++
	addIndexField(configId)

	// 配置 basePath 字段
	basePath := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexConfigName,
		FieldCode:    "basePath",
		FieldName:    "基础路径",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "配置信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	addIndexField(basePath)

	// 配置 jwtSecret 字段
	jwtSecret := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexConfigName,
		FieldCode:    "jwtSecret",
		FieldName:    "jwt加密字符串",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "配置信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	addIndexField(jwtSecret)

	// 配置 jwtSecret 字段
	jwttokenKey := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexConfigName,
		FieldCode:    "jwttokenKey",
		FieldName:    "jwt的key",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "配置信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	addIndexField(jwttokenKey)

	// 配置 serverPort 字段
	serverPort := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexConfigName,
		FieldCode:    "serverPort",
		FieldName:    "服务器ip:端口",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "配置信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	addIndexField(serverPort)

	// 配置 theme 字段
	theme := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexConfigName,
		FieldCode:    "theme",
		FieldName:    "主题",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "配置信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	addIndexField(theme)

	// 配置 timeout 字段
	timeout := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexConfigName,
		FieldCode:    "timeout",
		FieldName:    "超时时间",
		FieldType:    fieldType_数字,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "配置信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	addIndexField(timeout)

	// 添加公共字段
	indexCommonField(indexConfigName, "配置信息", sortNo, now)

	return true, nil
}

func init() {
	// 获取当前时间
	now := time.Now()

	//保存表信息
	//indexInfo, _, _ := openBleveIndex(indexInfoName)
	bleveSaveIndex(indexInfoName, indexConfigName, IndexInfoStruct{
		ID:         indexConfigName,
		Name:       "配置信息",
		Code:       "config",
		IndexType:  "index",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     2,
		Status:     1,
	})

}
