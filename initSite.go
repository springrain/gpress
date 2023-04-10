package main

import (
	"time"
)

// 初始化站点信息
func init() {

	if pathExist(bleveDataDir + indexSiteName) {
		return
	}

	// 创建用户表的索引
	mapping := bleveNewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultAnalyzer = gseAnalyzerName
	// //mapping.DefaultMapping.AddFieldMappingsAt("*", keywordMapping)
	// 获取当前时间
	now := time.Now()
	sortNo := 1
	// 初始化各个字段
	siteId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "id",
		FieldName:    "站点信息ID",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	// 放入文件中
	sortNo++
	addIndexField(mapping, siteId)

	siteTitle := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "title",
		FieldName:    "标题",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, siteTitle)
	// title 字段使用 中文分词器的mapping gseAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("title", gseAnalyzerMapping)

	siteName := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "name",
		FieldName:    "名称",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		IndexName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, siteName)

	domain := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "domain",
		FieldName:    "域名",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		IndexName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, domain)

	siteKeyword := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "keyword",
		FieldName:    "关键字",
		FieldType:    fieldType_文本框,
		AnalyzerName: commaAnalyzerName,
		IndexName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, siteKeyword)

	// keyword 字段使用 逗号分词器的mapping commaAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("keyword", commaAnalyzerMapping)

	siteDescription := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "description",
		FieldName:    "站点描述",
		FieldType:    fieldType_文本框,
		AnalyzerName: gseAnalyzerName,
		IndexName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, siteDescription)

	// description 字段使用 中文分词器的mapping gseAnalyzerMapping
	//mapping.DefaultMapping.AddFieldMappingsAt("description", gseAnalyzerMapping)

	siteTheme := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "theme",
		FieldName:    "默认主题",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, siteTheme)

	siteThemePC := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "themePC",
		FieldName:    "PC主题",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, siteThemePC)

	siteThemeWAP := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "themeWAP",
		FieldName:    "手机主题",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, siteThemeWAP)

	siteThemeWEIXIN := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "siteThemeWEIXIN",
		FieldName:    "微信主题",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, siteThemeWEIXIN)

	siteLogo := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "logo",
		FieldName:    "Logo",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, siteLogo)

	siteFavicon := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "favicon",
		FieldName:    "Favicon",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, siteFavicon)

	siteFooter := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexSiteName,
		FieldCode:    "footer",
		FieldName:    "底部信息",
		FieldType:    fieldType_文本框,
		AnalyzerName: keywordAnalyzerName,
		IndexName:    "站点信息",
		FieldComment: "",
		CreateTime:   now,
		UpdateTime:   now,
		CreateUser:   createUser,
		SortNo:       sortNo,
		Status:       3,
	}
	sortNo++
	addIndexField(mapping, siteFooter)

	// 添加公共字段
	indexCommonField(mapping, indexSiteName, "站点信息", sortNo, now)

	ok, err := bleveNew(indexSiteName, mapping)
	if err != nil || !ok {
		return
	}
	siteMap := make(map[string]interface{}, 0)
	siteMap["id"] = "gpress"
	siteMap["title"] = "gpress"
	siteMap["name"] = "gpress"
	siteMap["domain"] = "127.0.0.1"
	siteMap["keyword"] = "gpress"
	siteMap["description"] = "gpress"
	siteMap["theme"] = "default"
	siteMap["themePC"] = "default"
	siteMap["themeWAP"] = "default"
	siteMap["siteThemeWEIXIN"] = "default"
	siteMap["logo"] = "public/logo.png"
	siteMap["favicon"] = "public/favicon.png"
	siteMap["footer"] = `<div class="copyright">

	<span class="copyright-year">
	&copy; 
	2008 - 
	2023
	<span class="author">jiagou.com 版权所有 <a href='https://beian.miit.gov.cn' target='_blank'>豫ICP备2020026846号-1</a>   <a href='http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=41010302002680'  target='_blank'><img src='/public/gongan.png'>豫公网安备41010302002680号</a></span>
	</span>
	</div>`

	bleveSaveIndex(indexSiteName, "gpress", siteMap)
	//保存表信息
	bleveSaveIndex(indexInfoName, indexSiteName, IndexInfoStruct{
		ID:         indexSiteName,
		Name:       "站点信息",
		IndexType:  "index",
		Code:       "site",
		CreateTime: now,
		UpdateTime: now,
		CreateUser: createUser,
		SortNo:     7,
		Status:     1,
	})

}
