package main

import (
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
)

// initUser 初始化创建User索引
func initUser() (bool, error) {
	// 获取索引字段的表
	indexField := IndexMap[indexFieldIndexName]
	// 当前时间
	now := time.Now()

	// 用户表的 ID 字段
	userId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    userIndexName,
		IndexName:    "用户信息",
		FieldCode:    "ID",
		FieldName:    "用户ID",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       1,
		Active:       3,
	}
	indexField.Index(userId.ID, userId)

	// 用户表的 Account 字段
	userAccount := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    userIndexName,
		IndexName:    "用户信息",
		FieldCode:    "Account",
		FieldName:    "账号",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       2,
		Active:       1,
	}
	indexField.Index(userAccount.ID, userAccount)
	// 用户表的 PassWord 字段
	userPassWord := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    userIndexName,
		IndexName:    "用户信息",
		FieldCode:    "PassWord",
		FieldName:    "密码",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       3,
		Active:       1,
	}
	indexField.Index(userPassWord.ID, userPassWord)
	// 用户表的 UserName 字段
	userName := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    userIndexName,
		IndexName:    "用户信息",
		FieldCode:    "UserName",
		FieldName:    "用户名称",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       4,
		Active:       1,
	}
	indexField.Index(userName.ID, userName)

	// 创建用户表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name
	userIndex, err := bleve.New(userIndexName, mapping)
	if err != nil {
		FuncLogError(err)
		return false, err
	}
	IndexMap[userIndexName] = userIndex

	// 初始化数据
	user := make(map[string]string)
	id := FuncGenerateStringID()
	user["ID"] = id
	user["Account"] = "admin"
	user["PassWord"] = "21232f297a57a5a743894a0e4a801fc3"
	user["UserName"] = "管理员"
	// 初始化 admin用户
	userIndex.Index(id, user)
	return true, nil
}
