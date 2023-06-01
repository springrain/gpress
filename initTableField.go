package main

import (
	"context"
	"time"

	"gitee.com/chunanyong/zorm"
)

// TableFieldStruct 表和字段(表名:tableField)
// 记录所有表字段code和中文说明.
// 理论上所有的表字段都可以放到这个表里,因为都是Map,就不需要再单独指定表的字段了,可以动态创建Table(目前建议这样做)
type TableFieldStruct struct {
	// 引入默认的struct,隔离IEntityStruct的方法改动
	zorm.EntityStruct
	// ID 主键
	ID string `column:"id" json:"id"`
	// TableCode 表代码,类似表名 User,Site,PageTemplate,Category,Module,Content
	TableCode string `column:"tableCode" json:"tableCode,omitempty"`
	// TableName 表中文名
	TableName string `column:"tableName" json:"tableName,omitempty"`
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
	// AnalyzerName string `column:"analyzerName" json:"analyzerName,omitempty"`
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
func (entity *TableFieldStruct) GetTableName() string {
	return tableFieldName
}

// GetPKColumnName 获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称
// 不支持联合主键,变通认为无主键,业务控制实现(艰难取舍)
// 如果没有主键,也需要实现这个方法, return "" 即可
// IEntityStruct 接口的方法,实体类需要实现!!!
func (entity *TableFieldStruct) GetPKColumnName() string {
	return "id"
}

// initTableField 初始化创建TableField表
func initTableField() (bool, error) {
	if tableExist(tableFieldName) {
		return true, nil
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS tableField (
		id TEXT PRIMARY KEY     NOT NULL,
		tableCode         TEXT  NOT NULL,
		tableName         TEXT   NOT NULL,
		businessID        TEXT,
		fieldCode         TEXT NOT NULL,
		fieldName         TEXT NOT NULL,
		fieldComment      TEXT,
		fieldType         INT NOT NULL,
		fieldFormat       TEXT,
		required          INT ,
		defaultValue      TEXT,
		selectOption      TEXT,
		createTime        TEXT,
		updateTime        TEXT,
		createUser        TEXT,
		sortNo            int NOT NULL,
		status            int NOT NULL
	 ) strict ;`
	ctx := context.Background()
	_, err := execNativeSQL(ctx, createTableSQL)
	if err != nil {
		return false, err
	}

	sortNo := 1
	now := time.Now().Format("2006-01-02 15:04:05")
	id := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableFieldName,
		FieldCode: "id",
		FieldName: "字段ID",
		FieldType: fieldType_文本框,
		// AnalyzerName: keyword// AnalyzerName,
		TableName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	saveTableField(ctx, id)
	//bleveSaveTable
	//bleveSaveTable(indexFieldName, id.ID, id)
	sortNo++

	tableCode := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableFieldName,
		FieldCode: "tableCode",
		FieldName: "表名",
		FieldType: fieldType_文本框,
		// AnalyzerName: keyword// AnalyzerName,
		TableName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	saveTableField(ctx, tableCode)
	//bleveSaveTable
	//bleveSaveTable(indexFieldName, tableCode.ID, tableCode)
	sortNo++

	tableName := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableFieldName,
		FieldCode: "tableName",
		FieldName: "表名",
		FieldType: fieldType_文本框,
		// AnalyzerName: gse// AnalyzerName,
		TableName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	saveTableField(ctx, tableName)
	//bleveSaveTable
	//bleveSaveTable(indexFieldName, tableName.ID, tableName)
	sortNo++

	businessID := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableFieldName,
		FieldCode: "businessID",
		FieldName: "业务ID",
		FieldType: fieldType_文本框,
		// AnalyzerName: keyword// AnalyzerName,
		TableName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	saveTableField(ctx, businessID)
	//bleveSaveTable
	//bleveSaveTable(indexFieldName, businessID.ID, businessID)
	sortNo++

	fieldCode := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableFieldName,
		FieldCode: "fieldCode",
		FieldName: "字段",
		FieldType: fieldType_文本框,
		// AnalyzerName: keyword// AnalyzerName,
		TableName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	saveTableField(ctx, fieldCode)
	//bleveSaveTable
	//bleveSaveTable(indexFieldName, fieldCode.ID, fieldCode)
	sortNo++

	fieldName := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableFieldName,
		FieldCode: "fieldName",
		FieldName: "字段名",
		FieldType: fieldType_文本框,
		// AnalyzerName: gse// AnalyzerName,
		TableName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	saveTableField(ctx, fieldName)
	//bleveSaveTable
	//bleveSaveTable(indexFieldName, fieldName.ID, fieldName)
	sortNo++

	fieldComment := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableFieldName,
		FieldCode: "fieldComment",
		FieldName: "字段备注",
		FieldType: fieldType_文本框,
		// AnalyzerName: gse// AnalyzerName,
		TableName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	saveTableField(ctx, fieldComment)
	//bleveSaveTable
	//bleveSaveTable(indexFieldName, fieldComment.ID, fieldComment)
	sortNo++
	//ftSelect, _ := json.Marshal(fieldTypeMap)
	fieldType := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableFieldName,
		FieldCode: "fieldType",
		FieldName: "字段类型",
		FieldType: fieldType_数字,
		// AnalyzerName: numeric// AnalyzerName,
		//DefaultValue: "3",
		// FieldType  字段类型,数字(1),日期(2),文本框(3),文本域(4),富文本(5),下拉框(6),单选(7),多选(8),上传图片(9),上传附件(10),轮播图(11),音频(12),视频(13)
		//SelectOption: `[{"text":"数字","value":1},{"text":"日期","value":2},{"text":"文本框","value":3},{"text":"文本域","value":4},{"text":"富文本","value":5},{"text":"下拉框","value":6},{"text":"单选","value":7},{"text":"多选","value":8},{"text":"上传图片","value":9},{"text":"上传附件","value":10},{"text":"轮播图","value":11},{"text":"音频","value":12},{"text":"视频","value":13}]`,
		//SelectOption: string(ftSelect),
		TableName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	saveTableField(ctx, fieldType)
	//bleveSaveTable
	//bleveSaveTable(indexFieldName, fieldType.ID, fieldType)
	sortNo++

	fieldFormat := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableFieldName,
		FieldCode: "fieldFormat",
		FieldName: "数据格式",
		FieldType: fieldType_文本框,
		// AnalyzerName: keyword// AnalyzerName,
		TableName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	saveTableField(ctx, fieldFormat)
	//bleveSaveTable
	//bleveSaveTable(indexFieldName, fieldFormat.ID, fieldFormat)
	sortNo++

	required := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableFieldName,
		FieldCode: "required",
		FieldName: "是否必填",
		FieldType: fieldType_数字,
		// AnalyzerName: numeric// AnalyzerName,
		TableName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	saveTableField(ctx, required)
	//bleveSaveTable
	//bleveSaveTable(indexFieldName, required.ID, required)
	sortNo++

	defaultValue := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableFieldName,
		FieldCode: "defaultValue",
		FieldName: "默认值",
		FieldType: fieldType_文本框,
		// AnalyzerName: keyword// AnalyzerName,
		TableName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	saveTableField(ctx, defaultValue)
	//bleveSaveTable
	//bleveSaveTable(indexFieldName, defaultValue.ID, defaultValue)
	sortNo++

	/*
		analyzerName := TableFieldStruct{
			ID:        FuncGenerateStringID(),
			TableCode: tableFieldName,
			FieldCode: "analyzerName",
			FieldName: "分词器",
			FieldType: fieldType_文本框,
			// AnalyzerName: keyword// AnalyzerName,
			TableName:    "表字段",
			FieldComment: "",
			CreateTime:   now,
			UpdateTime:   now,
			CreateUser:   createUser,
			SortNo:       sortNo,
			Status:       3,
			Required:     1,
		}
		saveTableField(ctx, analyzerName)
		//bleveSaveTable
		//bleveSaveTable(indexFieldName, analyzerName.ID, analyzerName)
		sortNo++
	*/
	createTime := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableFieldName,
		FieldCode: "createTime",
		FieldName: "创建时间",
		FieldType: fieldType_日期,
		// AnalyzerName: datetime// AnalyzerName,
		TableName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	saveTableField(ctx, createTime)
	//bleveSaveTable
	//bleveSaveTable(indexFieldName, createTime.ID, createTime)
	sortNo++

	updateTime := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableFieldName,
		FieldCode: "updateTime",
		FieldName: "更新时间",
		FieldType: fieldType_日期,
		// AnalyzerName: datetime// AnalyzerName,
		TableName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	saveTableField(ctx, updateTime)
	//bleveSaveTable
	//bleveSaveTable(indexFieldName, updateTime.ID, updateTime)
	sortNo++

	createUserField := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableFieldName,
		FieldCode: "createUser",
		FieldName: "创建人",
		FieldType: fieldType_文本框,
		// AnalyzerName: keyword// AnalyzerName,
		TableName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	saveTableField(ctx, createUserField)
	//bleveSaveTable
	//bleveSaveTable(indexFieldName, createUserField.ID, createUserField)
	sortNo++

	sortNoField := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableFieldName,
		FieldCode: "sortNo",
		FieldName: "排序",
		FieldType: fieldType_数字,
		// AnalyzerName: numeric// AnalyzerName,
		TableName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	saveTableField(ctx, sortNoField)
	//bleveSaveTable
	//bleveSaveTable(indexFieldName, sortNoField.ID, sortNoField)
	sortNo++

	status := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableFieldName,
		FieldCode: "status",
		FieldName: "是否有效",
		FieldType: fieldType_数字,
		// AnalyzerName: numeric// AnalyzerName,
		TableName:    "表字段",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
		Required:     1,
	}
	saveTableField(ctx, status)
	//bleveSaveTable
	//bleveSaveTable(indexFieldName, status.ID, status)

	return true, nil
}

// indexCommonField 插入公共字段
func indexCommonField(ctx context.Context, tableCode string, tableName string, sortNo int, now string) {
	sortNo++
	commonCreateTime := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableCode,
		FieldCode: "createTime",
		FieldName: "创建时间",
		FieldType: fieldType_日期,
		// AnalyzerName: datetime// AnalyzerName,
		TableName:    tableName,
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, commonCreateTime)

	sortNo++
	commonUpdateTime := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableCode,
		FieldCode: "updateTime",
		FieldName: "更新时间",
		FieldType: fieldType_日期,
		// AnalyzerName: datetime// AnalyzerName,
		TableName:    tableName,
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, commonUpdateTime)

	sortNo++
	commonCreateUser := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableCode,
		FieldCode: "createUser",
		FieldName: "创建人",
		FieldType: fieldType_文本框,
		// AnalyzerName: keyword// AnalyzerName,
		TableName:    tableName,
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, commonCreateUser)

	sortNo++
	commonSortNo := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableCode,
		FieldCode: "sortNo",
		FieldName: "排序",
		FieldType: fieldType_数字,
		// AnalyzerName: numeric// AnalyzerName,
		TableName:    tableName,
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, commonSortNo)

	sortNo++
	commonActive := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableCode,
		FieldCode: "status",
		FieldName: "是否有效",
		FieldType: fieldType_数字,
		// AnalyzerName: numeric// AnalyzerName,
		TableName:    tableName,
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, commonActive)
}

func init() {
	if tableExist(tableFieldName) {
		return
	}
	// 获取当前时间
	now := time.Now().Format("2006-01-02 15:04:05")
	saveTableInfo(context.Background(), TableInfoStruct{
		ID:         tableFieldName,
		Name:       "表字段",
		Code:       "tableField",
		TableType:  "index",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     3,
		Status:     1,
	})
}
