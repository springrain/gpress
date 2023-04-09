package main

import (
	"time"

	"github.com/blevesearch/bleve/v2"
)

// initUser 初始化创建User索引
func init() {
	ok, err := indexExist(indexUserName)
	if err != nil || ok {
		return
	}
	// 当前时间
	now := time.Now()
	sortNo := 1
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
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
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
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
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
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
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
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       1,
	}
	sortNo++
	addIndexField(mapping, userName)

	// 添加公共字段
	// indexCommonField(indexField, indexUserName, 4, now)

	// //mapping.DefaultMapping.AddFieldMappingsAt("*", keywordMapping)
	_, err = bleveNew(indexUserName, mapping)
	if err != nil {

		return
	}

	//保存表信息
	bleveSaveIndex(indexInfoName, indexUserName, IndexInfoStruct{
		ID:         indexUserName,
		Name:       "用户信息",
		IndexType:  "index",
		Code:       "user",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     8,
		Status:     1,
	})
}
