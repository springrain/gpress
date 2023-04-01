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
	Active int `json:"active,omitempty"`
}

// initIndexInfo 初始化创建indexInfo索引
func initIndexInfo() (bool, error) {

	ok, err := openBleveIndex(indexInfoName)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	indexField := IndexMap[indexFieldName]
	// 创建内容表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器,存在问题:NewQueryStringQuery时不能正确匹配查询
	//mapping.DefaultMapping.DefaultAnalyzer = keywordAnalyzerName
	mapping.DefaultMapping.AddFieldMappingsAt("*", keywordMapping)

	// 获取当前时间
	now := time.Now()

	// 初始化各个字段
	// 主键 值为 IndexName,也就是表名
	infoId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexInfoName,
		FieldCode:    "id",
		FieldName:    "索引表ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       1,
		Active:       3,
	}
	indexField.Index(infoId.ID, infoId)
	// 索引名称,类似表名中文说明
	infoName := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexInfoName,
		FieldCode:    "name",
		FieldName:    "索引表名称",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       2,
		Active:       3,
	}
	indexField.Index(infoName.ID, infoName)
	//索引代码
	infoCode := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexInfoName,
		FieldCode:    "code",
		FieldName:    "索引表英文名",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       3,
		Active:       3,
	}
	indexField.Index(infoCode.ID, infoCode)
	// IndexType index/module 索引和模型,两种类型
	infoType := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexInfoName,
		FieldCode:    "indexType",
		FieldName:    "索引表ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       4,
		Active:       3,
	}
	indexField.Index(infoType.ID, infoType)

	// 添加公共字段
	indexCommonField(indexField, indexInfoName, 4, now)

	indexInfo, err := bleve.New(indexInfoName, mapping)
	if err != nil {
		FuncLogError(err)
		return false, err
	}
	IndexMap[indexInfoName] = indexInfo
	return true, nil
}
