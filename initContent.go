package main

import (
	"context"
	"time"
)

func init() {
	if tableExist(tableContentName) {
		return
	}

	// 创建内容表的表
	createTableSQL := `CREATE TABLE IF NOT EXISTS content (
		id TEXT PRIMARY KEY     NOT NULL,
		moduleID         TEXT  ,
		title         TEXT   NOT NULL,
		keyword           TEXT,
		description           TEXT,
		hrefURL           TEXT,
		subtitle           TEXT,
		categoryID           TEXT,
		categoryName           TEXT,
		templateID           TEXT,
		author           TEXT,
		tag           TEXT,
		toc           TEXT,
		summary           TEXT,
		content           TEXT,
		markdown           TEXT,
		thumbnail           TEXT,
		createTime        TEXT,
		updateTime        TEXT,
		createUser        TEXT,
		sortNo            int NOT NULL,
		status            int NOT NULL
	 ) strict ;`
	ctx := context.Background()
	_, err := crateTable(ctx, createTableSQL)
	if err != nil {
		return
	}

	//创建fts虚拟表,用于全文检索,id只作为查询字段,并不索引
	//id UNINDEXED, 去掉id,直接使用rowid
	createFTSTableSQL := `CREATE VIRTUAL TABLE IF NOT EXISTS fts_content USING fts5(
		title, 
		keyword, 
		description,
		subtitle,
		categoryName,
		summary,
		toc,
		tag,
		author, 

	    tokenize = 'simple 0',
		content='content', 
		content_rowid='rowid'
	);`
	_, err = crateTable(ctx, createFTSTableSQL)
	if err != nil {
		return
	}
	//创建触发器
	triggerContentSQL := `CREATE TRIGGER IF NOT EXISTS trigger_content_insert AFTER INSERT ON content
		BEGIN
			INSERT INTO fts_content (rowid, title, keyword, description,subtitle,categoryName,summary,toc,tag,author)
			VALUES (new.rowid,  new.title, new.keyword, new.description,new.subtitle,new.categoryName,new.summary,new.toc,new.tag,new.author);
		END;
	
	CREATE TRIGGER IF NOT EXISTS trigger_content_delete AFTER DELETE ON content
		BEGIN
			INSERT INTO fts_content (fts_content,  title, keyword, description,subtitle,categoryName,summary,toc,tag,author)
			VALUES ('delete',  old.title, old.keyword, old.description,old.subtitle,old.categoryName,old.summary,old.toc,old.tag,old.author);
		END;
	
	CREATE TRIGGER IF NOT EXISTS trigger_content_update AFTER UPDATE ON content
		BEGIN
			INSERT INTO fts_content (fts_content, rowid, title, keyword, description,subtitle,categoryName,summary,toc,tag,author)
			VALUES ('delete', old.rowid,  old.title, old.keyword, old.description,old.subtitle,old.categoryName,old.summary,old.toc,old.tag,old.author);
			INSERT INTO fts_content (rowid, title, keyword, description,subtitle,categoryName,summary,toc,tag,author)
			VALUES (new.rowid, new.title, new.keyword, new.description,new.subtitle,new.categoryName,new.summary,new.toc,new.tag,new.author);
		END;`
	_, err = crateTable(ctx, triggerContentSQL)
	if err != nil {
		return
	}

	// 获取当前时间
	now := time.Now().Format("2006-01-02 15:04:05")
	sortNo := 1
	// 初始化各个字段
	contentId := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "id",
		FieldName:    "文章内容ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	// 放入文件中
	sortNo++
	saveTableField(ctx, contentId)

	contentModuleID := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "moduleID",
		FieldName:    "模型的Code",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, contentModuleID)

	contentTitle := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableContentName,
		FieldCode: "title",
		FieldName: "标题",
		FieldType: fieldType_文本框,
		// 文章标题使用中文分词
		AnalyzerName: gseAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, contentTitle)
	// title 字段使用 中文分词器的mapping gseAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("title", gseAnalyzerMapping)

	contentKeyword := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableContentName,
		FieldCode: "keyword",
		FieldName: "关键字",
		FieldType: fieldType_文本框,
		// 文章关键字使用逗号分词器
		AnalyzerName: commaAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, contentKeyword)
	// keyword 字段使用 逗号分词器的mapping commaAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("keyword", commaAnalyzerMapping)

	contentDescription := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableContentName,
		FieldCode: "description",
		FieldName: "站点描述",
		FieldType: fieldType_文本框,
		// 文章描述使用中文分词器
		AnalyzerName: gseAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, contentDescription)
	// description 字段使用 中文分词器的mapping gseAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("description", gseAnalyzerMapping)

	contentHrefURL := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "hrefURL",
		FieldName:    "自身页面路径",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, contentHrefURL)

	contentSubtitle := TableFieldStruct{
		ID:        FuncGenerateStringID(),
		TableCode: tableContentName,
		FieldCode: "subtitle",
		FieldName: "副标题",
		FieldType: fieldType_文本框,
		// 文章副标题使用中文分词器
		AnalyzerName: gseAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, contentSubtitle)
	// subtitle 字段使用 中文分词器的mapping gseAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("subtitle", gseAnalyzerMapping)

	contentCategoryId := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "categoryID",
		FieldName:    "导航ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: commaAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, contentCategoryId)
	// categoryId 字段使用 逗号分词器的mapping commaAnalyzerMapping
	// //mapping.DefaultMapping.AddFieldMappingsAt("categoryId", commaAnalyzerMapping)

	contentCategoryNames := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "categoryName",
		FieldName:    "导航名称,逗号(,)隔开",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, contentCategoryNames)
	// categoryNames 字段使用 中文分词器的mapping gseAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("categoryNames", gseAnalyzerMapping)

	contentTemplateID := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "templateID",
		FieldName:    "模板Id",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, contentTemplateID)

	contentAuthor := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "author",
		FieldName:    "作者",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, contentAuthor)

	contentTag := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "tag",
		FieldName:    "标签",
		FieldType:    fieldType_文本框,
		AnalyzerName: commaAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, contentTag)

	contentToc := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "toc",
		FieldName:    "文章目录",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, contentToc)

	contentSummary := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "summary",
		FieldName:    "文章摘要",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, contentSummary)

	contentContent := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "content",
		FieldName:    "文章内容",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, contentContent)
	// content 字段使用 中文分词器的mapping gseAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("content", gseAnalyzerMapping)

	contentMarkdown := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "markdown",
		FieldName:    "markdown原文",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, contentMarkdown)

	contentThumbnail := TableFieldStruct{
		ID:           FuncGenerateStringID(),
		TableCode:    tableContentName,
		FieldCode:    "thumbnail",
		FieldName:    "封面图",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		TableName:    "文章内容",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	saveTableField(ctx, contentThumbnail)

	// 添加公共字段
	indexCommonField(ctx, tableContentName, "文章内容", sortNo, now)

	//保存表信息
	//tableInfo, _, _ := openBleveTable(tableInfoName)
	saveTableInfo(ctx, TableInfoStruct{
		ID:         tableContentName,
		Name:       "文章内容",
		Code:       "content",
		TableType:  "table",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     1,
		Status:     1,
	})

	//保存默认模型表信息
	saveTableInfo(ctx, TableInfoStruct{
		ID:         tableModuleDefaultName,
		Name:       "默认模型",
		Code:       tableModuleDefaultName,
		TableType:  "module",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     4,
		Status:     1,
	})

}
