package main

import (
	"time"

	"gitee.com/chunanyong/zorm"
)

// IndexInfoStruct 记录所有的表信息(索引名:indexInfo)
type IndexInfoStruct struct {
	// 引入默认的struct,隔离IEntityStruct的方法改动
	zorm.EntityStruct
	// ID 主键 值为 IndexName,也就是表名
	ID string `column:"id" json:"id"`
	// Name 索引名称,类似表名中文说明
	Name string `column:"name" json:"name,omitempty"`
	// Code 索引代码
	Code string `column:"code" json:"code,omitempty"`
	// IndexType index/module 索引和模型,两种类型
	IndexType string `column:"indexType" json:"indexType,omitempty"`
	// CreateTime 创建时间
	CreateTime time.Time `column:"createTime" json:"createTime,omitempty"`
	// UpdateTime 更新时间
	UpdateTime time.Time `column:"updateTime" json:"updateTime,omitempty"`
	// CreateUser  创建人,初始化 system
	CreateUser string `column:"createUser" json:"createUser,omitempty"`
	// SortNo 排序
	SortNo int `column:"sortNo" json:"sortNo,omitempty"`
	// 是否有效 无效(0),正常显示(1),界面不显示(3)
	Status int `column:"status" json:"status,omitempty"`
}

// GetTableName 获取表名称
// IEntityStruct 接口的方法,实体类需要实现!!!
func (entity *IndexInfoStruct) GetTableName() string {
	return indexInfoName
}

// GetPKColumnName 获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称
// 不支持联合主键,变通认为无主键,业务控制实现(艰难取舍)
// 如果没有主键,也需要实现这个方法, return "" 即可
// IEntityStruct 接口的方法,实体类需要实现!!!
func (entity *IndexInfoStruct) GetPKColumnName() string {
	return "id"
}

// initIndexInfo 初始化创建indexInfo索引
func initIndexInfo() (bool, error) {
	if tableExist(indexInfoName) {
		return true, nil
	}

	createTableSQL := `CREATE TABLE indexInfo (
		id TEXT PRIMARY KEY     NOT NULL,
		name         TEXT  NOT NULL,
		code         TEXT   NOT NULL,
		indexType         TEXT NOT NULL,
		createTime        TEXT,
		updateTime        TEXT,
		createUser        TEXT,
		sortNo            int NOT NULL,
		status            int NOT NULL
	 ) strict ;`
	_, err := bleveNewIndexMapping(createTableSQL)
	if err != nil {
		return false, err
	}
	// 获取当前时间
	now := time.Now()

	sortNo := 1
	// 初始化各个字段
	// 主键 值为 IndexName,也就是表名
	infoId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexInfoName,
		FieldCode:    "id",
		FieldName:    "索引表ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "表信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(infoId)
	// 索引名称,类似表名中文说明
	infoName := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexInfoName,
		FieldCode:    "name",
		FieldName:    "表名称",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		IndexName:    "表信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(infoName)
	//索引代码
	infoCode := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexInfoName,
		FieldCode:    "code",
		FieldName:    "表名",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "表信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(infoCode)
	// IndexType index/module 索引和模型,两种类型
	infoType := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexInfoName,
		FieldCode:    "indexType",
		FieldName:    "表类型",
		FieldType:    fieldType_下拉框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "表信息",
		DefaultValue: "index",
		SelectOption: `{"index":"索引表","module":"模型"]`,
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(infoType)

	// 添加公共字段
	indexCommonField(indexInfoName, "表信息", sortNo, now)

	return true, nil
}
