package main

import (
	"time"

	"github.com/blevesearch/bleve/v2"
)

// IndexInfoStruct 记录所有的表信息(索引名:indexInfo)
type IndexInfoStruct struct {
	// ID 主键 值为 IndexName,也就是表名
	ID string `json:"id"`
	// Name 索引名称,类似表名中文说明
	Name string `json:"name,omitempty"`
	// Code 索引代码
	Code string `json:"code,omitempty"`
	// IndexType index/module 索引和模型,两种类型
	IndexType string `json:"indexType,omitempty"`
	// CreateTime 创建时间
	CreateTime time.Time `json:"createTime,omitempty"`
	// UpdateTime 更新时间
	UpdateTime time.Time `json:"updateTime,omitempty"`
	// CreateUser  创建人,初始化 system
	CreateUser string `json:"createUser,omitempty"`
	// SortNo 排序
	SortNo int `json:"sortNo,omitempty"`
	// 是否有效 无效(0),正常显示(1),界面不显示(3)
	Status int `json:"status,omitempty"`
}

// initIndexInfo 初始化创建indexInfo索引
func initIndexInfo() (bool, error) {
	if pathExist(bleveDataDir + indexInfoName) {
		return false, nil
	}
	// 创建内容表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultAnalyzer = gseAnalyzerName
	// //mapping.DefaultMapping.AddFieldMappingsAt("*", keywordMapping)

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
	addIndexField(mapping, infoId)
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
	addIndexField(mapping, infoName)
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
	addIndexField(mapping, infoCode)
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
	addIndexField(mapping, infoType)

	// 添加公共字段
	indexCommonField(mapping, indexInfoName, "表信息", sortNo, now)

	ok, err := bleveNew(indexInfoName, mapping)
	if err != nil || !ok {
		return false, err
	}

	return true, nil
}
