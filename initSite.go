package main

import (
	"context"
	"time"

	"gitee.com/chunanyong/zorm"
)

// 初始化站点信息
func init() {
	if tableExist(tableSiteName) {
		return
	}

	// 创建用户表的表
	createTableSQL := `CREATE TABLE IF NOT EXISTS site (
		id TEXT PRIMARY KEY     NOT NULL,
		title         TEXT  NOT NULL,
		name         TEXT   NOT NULL,
		domain         TEXT,
		keyword         TEXT,
		description         TEXT,
		theme         TEXT,
		themePC         TEXT,
		themeWAP         TEXT,
		siteThemeWEIXIN         TEXT,
		logo         TEXT,
		favicon         TEXT,
		footer         TEXT,
		createTime        TEXT,
		updateTime        TEXT,
		createUser        TEXT,
		sortNo            int NOT NULL,
		status            int NOT NULL
	 ) strict ;`
	ctx := context.Background()
	_, err := execNativeSQL(ctx, createTableSQL)
	if err != nil {
		return
	}

	// //mapping.DefaultMapping.AddFieldMappingsAt("*", keywordMapping)
	// 获取当前时间
	now := time.Now().Format("2006-01-02 15:04:05")
	sortNo := 1
	// 初始化各个字段
	siteId := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableSiteName,
		FieldCode: "id",
		FieldName: "站点信息ID",
		FieldType: fieldType_文本框,
		// AnalyzerName: keyword// AnalyzerName,
		TableName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	// 放入文件中
	sortNo++
	saveTableField(ctx, siteId)

	siteTitle := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableSiteName,
		FieldCode: "title",
		FieldName: "标题",
		FieldType: fieldType_文本框,
		// AnalyzerName: keyword// AnalyzerName,
		TableName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, siteTitle)
	// title 字段使用 中文分词器的mapping gseAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("title", gseAnalyzerMapping)

	siteName := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableSiteName,
		FieldCode: "name",
		FieldName: "名称",
		FieldType: fieldType_文本框,
		// AnalyzerName: gse// AnalyzerName,
		TableName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, siteName)

	domain := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableSiteName,
		FieldCode: "domain",
		FieldName: "域名",
		FieldType: fieldType_文本框,
		// AnalyzerName: gse// AnalyzerName,
		TableName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, domain)

	siteKeyword := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableSiteName,
		FieldCode: "keyword",
		FieldName: "关键字",
		FieldType: fieldType_文本框,
		// AnalyzerName: comma// AnalyzerName,
		TableName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, siteKeyword)

	// keyword 字段使用 逗号分词器的mapping commaAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("keyword", commaAnalyzerMapping)

	siteDescription := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableSiteName,
		FieldCode: "description",
		FieldName: "站点描述",
		FieldType: fieldType_文本框,
		// AnalyzerName: gse// AnalyzerName,
		TableName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, siteDescription)

	// description 字段使用 中文分词器的mapping gseAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("description", gseAnalyzerMapping)

	siteTheme := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableSiteName,
		FieldCode: "theme",
		FieldName: "默认主题",
		FieldType: fieldType_文本框,
		// AnalyzerName: keyword// AnalyzerName,
		TableName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, siteTheme)

	siteThemePC := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableSiteName,
		FieldCode: "themePC",
		FieldName: "PC主题",
		FieldType: fieldType_文本框,
		// AnalyzerName: keyword// AnalyzerName,
		TableName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, siteThemePC)

	siteThemeWAP := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableSiteName,
		FieldCode: "themeWAP",
		FieldName: "手机主题",
		FieldType: fieldType_文本框,
		// AnalyzerName: keyword// AnalyzerName,
		TableName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, siteThemeWAP)

	siteThemeWEIXIN := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableSiteName,
		FieldCode: "siteThemeWEIXIN",
		FieldName: "微信主题",
		FieldType: fieldType_文本框,
		// AnalyzerName: keyword// AnalyzerName,
		TableName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, siteThemeWEIXIN)

	siteLogo := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableSiteName,
		FieldCode: "logo",
		FieldName: "Logo",
		FieldType: fieldType_文本框,
		// AnalyzerName: keyword// AnalyzerName,
		TableName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, siteLogo)

	siteFavicon := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableSiteName,
		FieldCode: "favicon",
		FieldName: "Favicon",
		FieldType: fieldType_文本框,
		// AnalyzerName: keyword// AnalyzerName,
		TableName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, siteFavicon)

	siteFooter := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableSiteName,
		FieldCode: "footer",
		FieldName: "底部信息",
		FieldType: fieldType_文本框,
		// AnalyzerName: keyword// AnalyzerName,
		TableName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, siteFooter)

	// 添加公共字段
	indexCommonField(ctx, tableSiteName, "站点信息", sortNo, now)

	siteMap := zorm.NewEntityMap(tableSiteName)
	siteMap.Set("id", "gpress")
	siteMap.Set("title", "gpress")
	siteMap.Set("name", "gpress")
	siteMap.Set("domain", "127.0.0.1")
	siteMap.Set("keyword", "gpress")
	siteMap.Set("description", "gpress")
	siteMap.Set("theme", "default")
	siteMap.Set("themePC", "default")
	siteMap.Set("themeWAP", "default")
	siteMap.Set("siteThemeWEIXIN", "default")
	siteMap.Set("logo", "public/logo.png")
	siteMap.Set("favicon", "public/favicon.png")
	siteMap.Set("footer", `<div class="copyright">

	<span class="copyright-year">
	&copy; 
	2008 - 
	2023
	<span class="author">jiagou.com 版权所有 <a href='https://beian.miit.gov.cn' target='_blank'>豫ICP备2020026846号-1</a>   <a href='http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=41010302002680'  target='_blank'><img src='/public/gongan.png'>豫公网安备41010302002680号</a></span>
	</span>
	</div>`)
	siteMap.Set("sortNo", 1)
	siteMap.Set("status", 1)
	//保存站点信息
	saveEntityMap(ctx, siteMap)
	//保存表信息
	saveTableInfo(ctx, TableInfoStruct{
		ID:         tableSiteName,
		Name:       "站点信息",
		TableType:  "table",
		Code:       "site",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     7,
		Status:     1,
	})

}
