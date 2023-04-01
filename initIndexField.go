package main

import (
	"time"

	"github.com/blevesearch/bleve/v2"
)

// IndexFieldStruct 索引和字段(索引名:indexField)
// 记录所有索引字段code和中文说明.
// 理论上所有的索引字段都可以放到这个表里,因为都是Map,就不需要再单独指定索引的字段了,可以动态创建Index(目前建议这样做)
type IndexFieldStruct struct {
	// ID 主键
	ID string `json:"id"`
	// IndexCode 索引代码,类似表名 User,Site,PageTemplate,NavMenu,Module,Content
	IndexCode string `json:"indexCode,omitempty"`
	// BusinessID  业务ID,处理业务记录临时增加的字段,意外情况
	BusinessID string `json:"businessID,omitempty"`
	// FieldCode  字段代码
	FieldCode string `json:"fFieldCode,omitempty"`
	// FieldName  字段中文名称
	FieldName string `json:"fieldName,omitempty"`
	// FieldType  字段类型,数字(1),日期(2),文本框(3),文本域(4),富文本(5),下拉框(6),单选(7),多选(8),上传图片(9),上传附件(10),轮播图(11),音频(12),视频(13)
	FieldType int `json:"fieldType,omitempty"`
	// FieldFormat 数据格式,用于日期或者数字
	FieldFormat string `json:"fieldFormat,omitempty"`
	// Required 字段是否为空. 0可以为空,1必填
	Required int `json:"required,omitempty"`
	// DefaultValue  默认值
	DefaultValue string `json:"defaultValue,omitempty"`
	// AnalyzerName  分词器名称
	AnalyzerName string `json:"analyzerName,omitempty"`
	// CreateTime 创建时间
	CreateTime time.Time `json:"createTime,omitempty"`
	// UpdateTime 更新时间
	UpdateTime time.Time `json:"updateTime,omitempty"`
	// CreateUser  创建人,初始化 system
	CreateUser string `json:"createUser,omitempty"`
	// SortNo 排序
	SortNo int `json:"sortNo,omitempty"`
	// 是否有效 无效(0),正常显示(1),界面不显示(3)
	Active int `json:"active,omitempty"`
}

// initIndexField 初始化创建IndexField索引
func initIndexField() (bool, error) {
	ok, err := openBleveIndex(indexFieldName)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keywordAnalyzerName
	// mapping.DefaultMapping.AddFieldMappingsAt("*", keywordMapping)
	index, err := bleve.New(indexFieldName, mapping)
	if err != nil {
		FuncLogError(err)
		return false, err
	}
	IndexMap[indexFieldName] = index
	return true, nil
}

// indexCommonField 插入公共字段
func indexCommonField(indexField bleve.Index, indexName string, sortNo int, now time.Time) {
	sortNo++
	commonCreateTime := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexName,
		FieldCode:    "createTime",
		FieldName:    "创建时间",
		FieldType:    fieldType_日期,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Active:       3,
	}
	indexField.Index(commonCreateTime.ID, commonCreateTime)

	sortNo++
	commonUpdateTime := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexName,
		FieldCode:    "updateTime",
		FieldName:    "更新时间",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       9,
		Active:       3,
	}
	indexField.Index(commonUpdateTime.ID, commonUpdateTime)

	sortNo++
	commonCreateUser := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexName,
		FieldCode:    "createUser",
		FieldName:    "创建人",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       10,
		Active:       3,
	}
	indexField.Index(commonCreateUser.ID, commonCreateUser)

	sortNo++
	commonSortNo := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexName,
		FieldCode:    "sortNo",
		FieldName:    "排序",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       9,
		Active:       3,
	}
	indexField.Index(commonSortNo.ID, commonSortNo)

	sortNo++
	commonActive := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexName,
		FieldCode:    "active",
		FieldName:    "是否有效",
		FieldType:    fieldType_数字,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       10,
		Active:       3,
	}
	indexField.Index(commonActive.ID, commonActive)
}
