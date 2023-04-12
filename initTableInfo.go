package main

import (
	"context"
	"time"

	"gitee.com/chunanyong/zorm"
)

// TableInfoStruct 记录所有的表信息(表名:tableInfo)
type TableInfoStruct struct {
	// 引入默认的struct,隔离IEntityStruct的方法改动
	zorm.EntityStruct
	// ID 主键 值为 TableName,也就是表名
	ID string `column:"id" json:"id"`
	// Name 表名称,类似表名中文说明
	Name string `column:"name" json:"name,omitempty"`
	// Code 表代码
	Code string `column:"code" json:"code,omitempty"`
	// TableType index/module 表和模型,两种类型
	TableType string `column:"tableType" json:"tableType,omitempty"`
	// CreateTime 创建时间
	CreateTime string `column:"createTime" json:"createTime,omitempty"`
	// UpdateTime 更新时间
	UpdateTime string `column:"updateTime" json:"updateTime,omitempty"`
	// CreateUser  创建人,初始化 system
	CreateUser string `column:"createUser" json:"createUser,omitempty"`
	// SortNo 排序
	SortNo int `column:"sortNo" json:"sortNo,omitempty"`
	// 是否有效 无效(0),正常显示(1),界面不显示(3)
	Status int `column:"status" json:"status,omitempty"`
}

// GetTableName 获取表名称
// IEntityStruct 接口的方法,实体类需要实现!!!
func (entity *TableInfoStruct) GetTableName() string {
	return tableInfoName
}

// GetPKColumnName 获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称
// 不支持联合主键,变通认为无主键,业务控制实现(艰难取舍)
// 如果没有主键,也需要实现这个方法, return "" 即可
// IEntityStruct 接口的方法,实体类需要实现!!!
func (entity *TableInfoStruct) GetPKColumnName() string {
	return "id"
}

// initTableInfo 初始化创建tableInfo表
func initTableInfo() (bool, error) {
	if tableExist(tableInfoName) {
		return true, nil
	}

	createTableSQL := `CREATE TABLE tableInfo (
		id TEXT PRIMARY KEY     NOT NULL,
		name         TEXT  NOT NULL,
		code         TEXT   NOT NULL,
		tableType         TEXT NOT NULL,
		createTime        TEXT,
		updateTime        TEXT,
		createUser        TEXT,
		sortNo            int NOT NULL,
		status            int NOT NULL
	 ) strict ;`
	ctx := context.Background()
	_, err := crateTable(ctx, createTableSQL)
	if err != nil {
		return false, err
	}
	// 获取当前时间
	now := time.Now().Format("2006-01-02 15:04:05")

	sortNo := 1
	// 初始化各个字段
	// 主键 值为 TableName,也就是表名
	infoId := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableInfoName,
		FieldCode:    "id",
		FieldName:    "表表ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "表信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, infoId)
	// 表名称,类似表名中文说明
	infoName := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableInfoName,
		FieldCode:    "name",
		FieldName:    "表名称",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		TableName:    "表信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, infoName)
	//表代码
	infoCode := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableInfoName,
		FieldCode:    "code",
		FieldName:    "表名",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "表信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, infoCode)
	// TableType index/module 表和模型,两种类型
	infoType := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableInfoName,
		FieldCode:    "tableType",
		FieldName:    "表类型",
		FieldType:    fieldType_下拉框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "表信息",
		DefaultValue: "table",
		SelectOption: `{"table":"表表","module":"模型"]`,
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, infoType)

	// 添加公共字段
	indexCommonField(ctx, tableInfoName, "表信息", sortNo, now)

	return true, nil
}
