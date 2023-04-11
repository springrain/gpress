package main

import (
	"context"
	"time"

	"gitee.com/chunanyong/zorm"
)

// IndexFieldStruct 索引和字段(索引名:indexField)
// 记录所有索引字段code和中文说明.
// 理论上所有的索引字段都可以放到这个表里,因为都是Map,就不需要再单独指定索引的字段了,可以动态创建Index(目前建议这样做)
type IndexFieldStruct struct {
	// 引入默认的struct,隔离IEntityStruct的方法改动
	zorm.EntityStruct
	// ID 主键
	ID string `column:"id" json:"id"`
	// IndexCode 索引代码,类似表名 User,Site,PageTemplate,NavMenu,Module,Content
	IndexCode string `column:"indexCode" json:"indexCode,omitempty"`
	// IndexName 索引表中文名
	IndexName string `column:"indexName" json:"indexName,omitempty"`
	// BusinessID  业务ID,处理业务记录临时增加的字段,意外情况
	BusinessID string `column:"businessID" json:"businessID,omitempty"`
	// FieldCode  字段代码
	FieldCode string `column:"fieldCode" json:"fieldCode,omitempty"`
	// FieldName  字段中文名称
	FieldName string `column:"fieldName" json:"fieldName,omitempty"`
	// FieldComment 字段说明
	FieldComment string `column:"fieldComment" json:"fieldComment,omitempty"`
	// FieldType  字段类型,数字(1),日期(2),文本框(3),文本域(4),富文本(5),下拉框(6),单选(7),多选(8),上传图片(9),上传附件(10),轮播图(11),音频(12),视频(13)
	FieldType int `column:"fieldType" json:"fieldType,omitempty"`
	// FieldFormat 数据格式,用于日期或者数字
	FieldFormat string `column:"fieldFormat" json:"fieldFormat,omitempty"`
	// Required 字段是否为空. 0可以为空,1必填
	Required int `column:"required" json:"required,omitempty"`
	// DefaultValue  默认值
	DefaultValue string `column:"defaultValue" json:"defaultValue,omitempty"`
	//SelectOption 下拉框的选项值
	SelectOption string `column:"selectOption" json:"selectOption,omitempty"`
	// AnalyzerName  分词器名称
	AnalyzerName string `column:"analyzerName" json:"analyzerName,omitempty"`
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
func (entity *IndexFieldStruct) GetTableName() string {
	return indexFieldName
}

// GetPKColumnName 获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称
// 不支持联合主键,变通认为无主键,业务控制实现(艰难取舍)
// 如果没有主键,也需要实现这个方法, return "" 即可
// IEntityStruct 接口的方法,实体类需要实现!!!
func (entity *IndexFieldStruct) GetPKColumnName() string {
	return "id"
}

// initIndexField 初始化创建IndexField索引
func initIndexField() (bool, error) {
	if tableExist(indexFieldName) {
		return true, nil
	}

	createTableSQL := `CREATE TABLE indexField (
		id TEXT PRIMARY KEY     NOT NULL,
		indexCode         TEXT  NOT NULL,
		indexName         TEXT   NOT NULL,
		businessID        TEXT,
		fieldCode         TEXT NOT NULL,
		fieldName         TEXT NOT NULL,
		fieldComment      TEXT,
		fieldType         INT NOT NULL,
		fieldFormat       TEXT,
		required          INT ,
		defaultValue      TEXT,
		selectOption      TEXT,
		analyzerName      TEXT,
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
	addIndexField(id)
	//bleveSaveIndex
	//bleveSaveIndex(indexFieldName, id.ID, id)
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
	addIndexField(indexCode)
	//bleveSaveIndex
	//bleveSaveIndex(indexFieldName, indexCode.ID, indexCode)
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
	addIndexField(indexName)
	//bleveSaveIndex
	//bleveSaveIndex(indexFieldName, indexName.ID, indexName)
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
	addIndexField(businessID)
	//bleveSaveIndex
	//bleveSaveIndex(indexFieldName, businessID.ID, businessID)
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
	addIndexField(fieldCode)
	//bleveSaveIndex
	//bleveSaveIndex(indexFieldName, fieldCode.ID, fieldCode)
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
	addIndexField(fieldName)
	//bleveSaveIndex
	//bleveSaveIndex(indexFieldName, fieldName.ID, fieldName)
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
	addIndexField(fieldComment)
	//bleveSaveIndex
	//bleveSaveIndex(indexFieldName, fieldComment.ID, fieldComment)
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
	addIndexField(fieldType)
	//bleveSaveIndex
	//bleveSaveIndex(indexFieldName, fieldType.ID, fieldType)
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
	addIndexField(fieldFormat)
	//bleveSaveIndex
	//bleveSaveIndex(indexFieldName, fieldFormat.ID, fieldFormat)
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
	addIndexField(required)
	//bleveSaveIndex
	//bleveSaveIndex(indexFieldName, required.ID, required)
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
	addIndexField(defaultValue)
	//bleveSaveIndex
	//bleveSaveIndex(indexFieldName, defaultValue.ID, defaultValue)
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
	addIndexField(analyzerName)
	//bleveSaveIndex
	//bleveSaveIndex(indexFieldName, analyzerName.ID, analyzerName)
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
	addIndexField(createTime)
	//bleveSaveIndex
	//bleveSaveIndex(indexFieldName, createTime.ID, createTime)
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
	addIndexField(updateTime)
	//bleveSaveIndex
	//bleveSaveIndex(indexFieldName, updateTime.ID, updateTime)
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
	addIndexField(createUserField)
	//bleveSaveIndex
	//bleveSaveIndex(indexFieldName, createUserField.ID, createUserField)
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
	addIndexField(sortNoField)
	//bleveSaveIndex
	//bleveSaveIndex(indexFieldName, sortNoField.ID, sortNoField)
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
	addIndexField(status)
	//bleveSaveIndex
	//bleveSaveIndex(indexFieldName, status.ID, status)

	return true, nil
}

// indexCommonField 插入公共字段
func indexCommonField(indexCode string, indexName string, sortNo int, now time.Time) {
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
	addIndexField(commonCreateTime)

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
	addIndexField(commonUpdateTime)

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
	addIndexField(commonCreateUser)

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
	addIndexField(commonSortNo)

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
	addIndexField(commonActive)
}

func addIndexField(indexFiledStruct IndexFieldStruct) {
	zorm.Transaction(context.Background(), func(ctx context.Context) (interface{}, error) {
		_, err := zorm.Insert(ctx, &indexFiledStruct)
		return nil, err
	})
}
func init7() {

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
