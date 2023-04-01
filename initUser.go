package main

import (
	"time"

	"github.com/blevesearch/bleve/v2"
)

// initUser 初始化创建User索引
func init() {
	ok, err := openBleveIndex(indexUserName)
	if err != nil || ok {
		return
	}
	// 当前时间
	now := time.Now()

	// 创建用户表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultAnalyzer = keywordAnalyzerName
	// 用户表的 ID 字段
	userId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexUserName,
		FieldCode:    "id",
		FieldName:    "用户ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "用户信息",
		FieldComment: "",
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       1,
		Active:       3,
	}
	addIndexField(mapping, userId)

	// 用户表的 Account 字段
	userAccount := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexUserName,
		FieldCode:    "account",
		FieldName:    "账号",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "用户信息",
		FieldComment: "",
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       2,
		Active:       1,
	}
	addIndexField(mapping, userAccount)
	// 用户表的 PassWord 字段
	userPassWord := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexUserName,
		FieldCode:    "password",
		FieldName:    "密码",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "用户信息",
		FieldComment: "",
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       3,
		Active:       1,
	}
	addIndexField(mapping, userPassWord)
	// 用户表的 UserName 字段
	userName := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexUserName,
		FieldCode:    "userName",
		FieldName:    "用户名称",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "用户信息",
		FieldComment: "",
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       4,
		Active:       1,
	}
	addIndexField(mapping, userName)

	// 添加公共字段
	// indexCommonField(indexField, indexUserName, 4, now)

	// //mapping.DefaultMapping.AddFieldMappingsAt("*", keywordMapping)
	userIndex, err := bleve.New(indexUserName, mapping)
	if err != nil {
		FuncLogError(err)
		return
	}
	IndexMap[indexUserName] = userIndex

	//保存表信息
	indexInfo := IndexMap[indexInfoName]
	indexInfo.Index(indexUserName, IndexInfoStruct{
		ID:         indexUserName,
		Name:       "用户信息",
		IndexType:  "index",
		Code:       "user",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     2,
		Active:     1,
	})
}
