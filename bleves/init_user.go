package bleves

import (
	"gitee.com/gpress/gpress/constant"
	"gitee.com/gpress/gpress/logger"
	"gitee.com/gpress/gpress/model"
	"gitee.com/gpress/gpress/util"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
)

// initUser 初始化创建User索引
func initUser() (bool, error) {
	// 获取索引字段的表
	indexField := IndexMap[constant.INDEX_FIELD_INDEX_NAME]
	// 当前时间
	now := time.Now()

	// 用户表的 ID 字段
	userId := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.USER_INDEX_NAME,
		IndexName:    "用户信息",
		FieldCode:    "id",
		FieldName:    "用户ID",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       1,
		Active:       3,
	}
	err := indexField.Index(userId.ID, userId)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	// 用户表的 Account 字段
	userAccount := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.USER_INDEX_NAME,
		IndexName:    "用户信息",
		FieldCode:    "account",
		FieldName:    "账号",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       2,
		Active:       1,
	}
	err = indexField.Index(userAccount.ID, userAccount)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	// 用户表的 PassWord 字段
	userPassWord := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.USER_INDEX_NAME,
		IndexName:    "用户信息",
		FieldCode:    "password",
		FieldName:    "密码",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       3,
		Active:       1,
	}
	err = indexField.Index(userPassWord.ID, userPassWord)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	// 用户表的 UserName 字段
	userName := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.USER_INDEX_NAME,
		IndexName:    "用户信息",
		FieldCode:    "userName",
		FieldName:    "用户名称",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       4,
		Active:       1,
	}
	err = indexField.Index(userName.ID, userName)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	// 创建用户表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name
	userIndex, err := bleve.New(constant.USER_INDEX_NAME, mapping)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	IndexMap[constant.USER_INDEX_NAME] = userIndex
	return true, nil
}
