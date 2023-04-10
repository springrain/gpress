package main

import (
	"time"

	"github.com/blevesearch/bleve/v2/mapping"
)

// IndexFieldStruct 索引和字段(索引名:indexField)
// 记录所有索引字段code和中文说明.
// 理论上所有的索引字段都可以放到这个表里,因为都是Map,就不需要再单独指定索引的字段了,可以动态创建Index(目前建议这样做)
type IndexFieldStruct struct {
	// ID 主键
	ID string `json:"id"`
	// IndexCode 索引代码,类似表名 User,Site,PageTemplate,NavMenu,Module,Content
	IndexCode string `json:"indexCode,omitempty"`
	// IndexName 索引表中文名
	IndexName string `json:"indexName,omitempty"`
	// BusinessID  业务ID,处理业务记录临时增加的字段,意外情况
	BusinessID string `json:"businessID,omitempty"`
	// FieldCode  字段代码
	FieldCode string `json:"fieldCode,omitempty"`
	// FieldName  字段中文名称
	FieldName string `json:"fieldName,omitempty"`
	// FieldComment 字段说明
	FieldComment string `json:"fieldComment,omitempty"`
	// FieldType  字段类型,数字(1),日期(2),文本框(3),文本域(4),富文本(5),下拉框(6),单选(7),多选(8),上传图片(9),上传附件(10),轮播图(11),音频(12),视频(13)
	FieldType int `json:"fieldType,omitempty"`
	// FieldFormat 数据格式,用于日期或者数字
	FieldFormat string `json:"fieldFormat,omitempty"`
	// Required 字段是否为空. 0可以为空,1必填
	Required int `json:"required,omitempty"`
	// DefaultValue  默认值
	DefaultValue string `json:"defaultValue,omitempty"`
	//SelectOption 下拉框的选项值
	SelectOption string `json:"selectOption,omitempty"`
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
	Status int `json:"status,omitempty"`
}

// initIndexField 初始化创建IndexField索引
func initIndexField() (bool, error) {
	if pathExist(bleveDataDir + indexFieldName) {
		return true, nil
	}

	mapping := bleveNewIndexMapping()
	// 指定默认的分词器,为了检索字段名可以分词,默认分词器为gse,其他字段都要手动指定为keyword
	mapping.DefaultAnalyzer = gseAnalyzerName

	mapping.DefaultMapping.AddFieldMappingsAt("id", keywordAnalyzerMapping)
	mapping.DefaultMapping.AddFieldMappingsAt("indexCode", keywordAnalyzerMapping)
	mapping.DefaultMapping.AddFieldMappingsAt("indexName", gseAnalyzerMapping)
	mapping.DefaultMapping.AddFieldMappingsAt("businessID", keywordAnalyzerMapping)
	mapping.DefaultMapping.AddFieldMappingsAt("fieldCode", keywordAnalyzerMapping)
	mapping.DefaultMapping.AddFieldMappingsAt("fieldName", gseAnalyzerMapping)
	mapping.DefaultMapping.AddFieldMappingsAt("fieldComment", gseAnalyzerMapping)
	mapping.DefaultMapping.AddFieldMappingsAt("fieldType", numericAnalyzerMapping)
	mapping.DefaultMapping.AddFieldMappingsAt("fieldFormat", keywordAnalyzerMapping)
	mapping.DefaultMapping.AddFieldMappingsAt("required", numericAnalyzerMapping)
	mapping.DefaultMapping.AddFieldMappingsAt("defaultValue", keywordAnalyzerMapping)
	mapping.DefaultMapping.AddFieldMappingsAt("analyzerName", keywordAnalyzerMapping)
	mapping.DefaultMapping.AddFieldMappingsAt("createTime", datetimeAnalyzerMapping)
	mapping.DefaultMapping.AddFieldMappingsAt("updateTime", datetimeAnalyzerMapping)
	mapping.DefaultMapping.AddFieldMappingsAt("createUser", keywordAnalyzerMapping)
	mapping.DefaultMapping.AddFieldMappingsAt("sortNo", numericAnalyzerMapping)
	mapping.DefaultMapping.AddFieldMappingsAt("status", numericAnalyzerMapping)

	ok, err := bleveNew(indexFieldName, mapping)
	if err != nil || !ok {
		return false, err
	}

	sortNo := 1
	now := time.Now()
	id := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexFieldName,
		FieldCode:    "id",
		FieldName:    "字段ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	bleveSaveIndex(indexFieldName, id.ID, id)
	sortNo++

	indexCode := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexFieldName,
		FieldCode:    "indexCode",
		FieldName:    "表名",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	bleveSaveIndex(indexFieldName, indexCode.ID, indexCode)
	sortNo++

	indexName := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexFieldName,
		FieldCode:    "indexName",
		FieldName:    "表名",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		IndexName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	bleveSaveIndex(indexFieldName, indexName.ID, indexName)
	sortNo++

	businessID := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexFieldName,
		FieldCode:    "businessID",
		FieldName:    "业务ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	bleveSaveIndex(indexFieldName, businessID.ID, businessID)
	sortNo++

	fieldCode := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexFieldName,
		FieldCode:    "fieldCode",
		FieldName:    "字段",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	bleveSaveIndex(indexFieldName, fieldCode.ID, fieldCode)
	sortNo++

	fieldName := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexFieldName,
		FieldCode:    "fieldName",
		FieldName:    "字段名",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		IndexName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	bleveSaveIndex(indexFieldName, fieldName.ID, fieldName)
	sortNo++

	fieldComment := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexFieldName,
		FieldCode:    "fieldComment",
		FieldName:    "字段备注",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		IndexName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	bleveSaveIndex(indexFieldName, fieldComment.ID, fieldComment)
	sortNo++
	//ftSelect, _ := json.Marshal(fieldTypeMap)
	fieldType := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexFieldName,
		FieldCode:    "fieldType",
		FieldName:    "字段类型",
		FieldType:    fieldType_数字,
		AnalyzerName: numericAnalyzerName,
		//DefaultValue: "3",
		// FieldType  字段类型,数字(1),日期(2),文本框(3),文本域(4),富文本(5),下拉框(6),单选(7),多选(8),上传图片(9),上传附件(10),轮播图(11),音频(12),视频(13)
		//SelectOption: `[{"text":"数字","value":1},{"text":"日期","value":2},{"text":"文本框","value":3},{"text":"文本域","value":4},{"text":"富文本","value":5},{"text":"下拉框","value":6},{"text":"单选","value":7},{"text":"多选","value":8},{"text":"上传图片","value":9},{"text":"上传附件","value":10},{"text":"轮播图","value":11},{"text":"音频","value":12},{"text":"视频","value":13}]`,
		//SelectOption: string(ftSelect),
		IndexName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	bleveSaveIndex(indexFieldName, fieldType.ID, fieldType)
	sortNo++

	fieldFormat := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexFieldName,
		FieldCode:    "fieldFormat",
		FieldName:    "数据格式",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	bleveSaveIndex(indexFieldName, fieldFormat.ID, fieldFormat)
	sortNo++

	required := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexFieldName,
		FieldCode:    "required",
		FieldName:    "是否必填",
		FieldType:    fieldType_数字,
		AnalyzerName: numericAnalyzerName,
		IndexName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	bleveSaveIndex(indexFieldName, required.ID, required)
	sortNo++

	defaultValue := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexFieldName,
		FieldCode:    "defaultValue",
		FieldName:    "默认值",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	bleveSaveIndex(indexFieldName, defaultValue.ID, defaultValue)
	sortNo++

	analyzerName := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexFieldName,
		FieldCode:    "analyzerName",
		FieldName:    "分词器",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	bleveSaveIndex(indexFieldName, analyzerName.ID, analyzerName)
	sortNo++

	createTime := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexFieldName,
		FieldCode:    "createTime",
		FieldName:    "创建时间",
		FieldType:    fieldType_日期,
		AnalyzerName: datetimeAnalyzerName,
		IndexName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	bleveSaveIndex(indexFieldName, createTime.ID, createTime)
	sortNo++

	updateTime := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexFieldName,
		FieldCode:    "updateTime",
		FieldName:    "更新时间",
		FieldType:    fieldType_日期,
		AnalyzerName: datetimeAnalyzerName,
		IndexName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	bleveSaveIndex(indexFieldName, updateTime.ID, updateTime)
	sortNo++

	createUserField := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexFieldName,
		FieldCode:    "createUser",
		FieldName:    "创建人",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	bleveSaveIndex(indexFieldName, createUserField.ID, createUserField)
	sortNo++

	sortNoField := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexFieldName,
		FieldCode:    "sortNo",
		FieldName:    "排序",
		FieldType:    fieldType_数字,
		AnalyzerName: numericAnalyzerName,
		IndexName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	bleveSaveIndex(indexFieldName, sortNoField.ID, sortNoField)
	sortNo++

	status := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexFieldName,
		FieldCode:    "status",
		FieldName:    "是否有效",
		FieldType:    fieldType_数字,
		AnalyzerName: numericAnalyzerName,
		IndexName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	bleveSaveIndex(indexFieldName, status.ID, status)

	return true, nil
}

// indexCommonField 插入公共字段
func indexCommonField(mapping *mapping.IndexMappingImpl, indexCode string, indexName string, sortNo int, now time.Time) {
	sortNo++
	commonCreateTime := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexCode,
		FieldCode:    "createTime",
		FieldName:    "创建时间",
		FieldType:    fieldType_日期,
		AnalyzerName: datetimeAnalyzerName,
		IndexName:    indexName,
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, commonCreateTime)

	sortNo++
	commonUpdateTime := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexCode,
		FieldCode:    "updateTime",
		FieldName:    "更新时间",
		FieldType:    fieldType_日期,
		AnalyzerName: datetimeAnalyzerName,
		IndexName:    indexName,
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, commonUpdateTime)

	sortNo++
	commonCreateUser := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexCode,
		FieldCode:    "createUser",
		FieldName:    "创建人",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    indexName,
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, commonCreateUser)

	sortNo++
	commonSortNo := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexCode,
		FieldCode:    "sortNo",
		FieldName:    "排序",
		FieldType:    fieldType_数字,
		AnalyzerName: numericAnalyzerName,
		IndexName:    indexName,
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, commonSortNo)

	sortNo++
	commonActive := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexCode,
		FieldCode:    "status",
		FieldName:    "是否有效",
		FieldType:    fieldType_数字,
		AnalyzerName: numericAnalyzerName,
		IndexName:    indexName,
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, commonActive)
}

func addIndexField(bleveMapping *mapping.IndexMappingImpl, indexFiledStruct IndexFieldStruct) {
	// 获取索引字段的表
	bleveSaveIndex(indexFieldName, indexFiledStruct.ID, indexFiledStruct)
	if bleveMapping == nil {
		return
	}
	var analyzerMapping *mapping.FieldMapping

	switch indexFiledStruct.AnalyzerName {
	case keywordAnalyzerName:
		analyzerMapping = keywordAnalyzerMapping
	case commaAnalyzerName:
		analyzerMapping = commaAnalyzerMapping
	case datetimeAnalyzerName:
		analyzerMapping = datetimeAnalyzerMapping
	case numericAnalyzerName:
		analyzerMapping = numericAnalyzerMapping
	default:
		analyzerMapping = gseAnalyzerMapping
	}
	bleveMapping.DefaultMapping.AddFieldMappingsAt(indexFiledStruct.FieldCode, analyzerMapping)
}
func init() {

	// 获取当前时间
	now := time.Now()

	bleveSaveIndex(indexInfoName, indexFieldName, IndexInfoStruct{
		ID:         indexFieldName,
		Name:       "表字段",
		Code:       "indexField",
		IndexType:  "index",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     3,
		Status:     1,
	})
}
