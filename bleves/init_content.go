package bleves

import (
	"gitee.com/gpress/gpress/constant"
	"gitee.com/gpress/gpress/logger"
	"gitee.com/gpress/gpress/model"
	"gitee.com/gpress/gpress/util"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"time"
)

func initContent() (bool, error) {
	indexField := IndexMap[constant.INDEX_FIELD_INDEX_NAME]

	// 创建内容表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name // 这是要换成逗号分词吧

	// 获取当前时间
	now := time.Now()

	// 初始化各个字段
	contentId := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_CONTENT_NAME,
		IndexName:    "文章内容",
		FieldCode:    "ID",
		FieldName:    "文章内容ID",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       1,
		Active:       3,
	}
	// 放入文件中
	err := indexField.Index(contentId.ID, contentId)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	contentModuleIndexCode := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_CONTENT_NAME,
		IndexName:    "文章内容",
		FieldCode:    "ModuleIndexCode",
		FieldName:    "模型的Code",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       2,
		Active:       3,
	}
	err = indexField.Index(contentModuleIndexCode.ID, contentModuleIndexCode)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	contentTitle := model.IndexFieldStruct{
		ID:        util.FuncGenerateStringID(),
		IndexCode: constant.INDEX_CONTENT_NAME,
		IndexName: "文章内容",
		FieldCode: "Title",
		FieldName: "标题",
		FieldType: 3,
		// 文章标题使用中文分词
		AnalyzerName: constant.GSE_ANGLYZER_NAME,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       2,
		Active:       3,
	}
	err = indexField.Index(contentTitle.ID, contentTitle)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	// Title 字段使用 中文分词器的mapping gseAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("Title", gseAnalyzerMapping)

	contentKeyWords := model.IndexFieldStruct{
		ID:        util.FuncGenerateStringID(),
		IndexCode: constant.INDEX_CONTENT_NAME,
		IndexName: "文章内容",
		FieldCode: "KeyWords",
		FieldName: "关键字",
		FieldType: 3,
		// 文章关键字使用逗号分词器
		AnalyzerName: constant.COMMA_ANALYZER_NAME,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       3,
		Active:       3,
	}
	err = indexField.Index(contentKeyWords.ID, contentKeyWords)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	// KeyWords 字段使用 逗号分词器的mapping commaAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("KeyWords", commaAnalyzerMapping)

	contentDescription := model.IndexFieldStruct{
		ID:        util.FuncGenerateStringID(),
		IndexCode: constant.INDEX_CONTENT_NAME,
		IndexName: "文章内容",
		FieldCode: "Description",
		FieldName: "站点描述",
		FieldType: 3,
		// 文章描述使用中文分词器
		AnalyzerName: constant.GSE_ANGLYZER_NAME,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       4,
		Active:       3,
	}
	err = indexField.Index(contentDescription.ID, contentDescription)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	// Description 字段使用 中文分词器的mapping gseAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("Description", gseAnalyzerMapping)

	contentPageURL := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_CONTENT_NAME,
		IndexName:    "文章内容",
		FieldCode:    "PageURL",
		FieldName:    "自身页面路径",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       5,
		Active:       3,
	}
	err = indexField.Index(contentPageURL.ID, contentPageURL)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	contentSubtitle := model.IndexFieldStruct{
		ID:        util.FuncGenerateStringID(),
		IndexCode: constant.INDEX_CONTENT_NAME,
		IndexName: "文章内容",
		FieldCode: "Subtitle",
		FieldName: "副标题",
		FieldType: 3,
		// 文章副标题使用中文分词器
		AnalyzerName: constant.GSE_ANGLYZER_NAME,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       6,
		Active:       3,
	}
	err = indexField.Index(contentSubtitle.ID, contentSubtitle)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	// Subtitle 字段使用 中文分词器的mapping gseAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("Subtitle", gseAnalyzerMapping)

	contentNavMenuId := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_CONTENT_NAME,
		IndexName:    "文章内容",
		FieldCode:    "NavMenuId",
		FieldName:    "导航ID,逗号(,)隔开",
		FieldType:    3,
		AnalyzerName: constant.COMMA_ANALYZER_NAME,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       4,
		Active:       3,
	}
	err = indexField.Index(contentNavMenuId.ID, contentNavMenuId)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	// NavMenuId 字段使用 逗号分词器的mapping commaAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("NavMenuId", commaAnalyzerMapping)

	contentNavMenuName := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_CONTENT_NAME,
		IndexName:    "文章内容",
		FieldCode:    "NavMenuName",
		FieldName:    "导航名称",
		FieldType:    3,
		AnalyzerName: constant.GSE_ANGLYZER_NAME,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       5,
		Active:       3,
	}
	err = indexField.Index(contentNavMenuName.ID, contentNavMenuName)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	// NavMenuName 字段使用 中文分词器的mapping gseAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("NavMenuName", gseAnalyzerMapping)

	contentTemplateID := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_CONTENT_NAME,
		IndexName:    "文章内容",
		FieldCode:    "TemplateID",
		FieldName:    "模板Id",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       6,
		Active:       3,
	}
	err = indexField.Index(contentTemplateID.ID, contentTemplateID)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	contentContent := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_CONTENT_NAME,
		IndexName:    "文章内容",
		FieldCode:    "Content",
		FieldName:    "文章内容",
		FieldType:    3,
		AnalyzerName: constant.GSE_ANGLYZER_NAME,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       7,
		Active:       3,
	}
	err = indexField.Index(contentContent.ID, contentContent)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	// Content 字段使用 中文分词器的mapping gseAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("Content", gseAnalyzerMapping)

	contentCreateTime := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_CONTENT_NAME,
		IndexName:    "文章内容",
		FieldCode:    "CreateTime",
		FieldName:    "创建时间",
		FieldType:    2,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       8,
		Active:       3,
	}
	indexField.Index(contentCreateTime.ID, contentCreateTime)

	moduleUpdateTime := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_CONTENT_NAME,
		IndexName:    "文章内容",
		FieldCode:    "UpdateTime",
		FieldName:    "更新时间",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       9,
		Active:       3,
	}
	err = indexField.Index(moduleUpdateTime.ID, moduleUpdateTime)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	moduleCreateUser := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_CONTENT_NAME,
		IndexName:    "文章内容",
		FieldCode:    "CreateUser",
		FieldName:    "创建人",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       10,
		Active:       3,
	}
	err = indexField.Index(moduleCreateUser.ID, moduleCreateUser)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	moduleSortNo := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_CONTENT_NAME,
		IndexName:    "文章内容",
		FieldCode:    "SortNo",
		FieldName:    "排序",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       9,
		Active:       3,
	}
	err = indexField.Index(moduleSortNo.ID, moduleSortNo)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	moduleActive := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_CONTENT_NAME,
		IndexName:    "文章内容",
		FieldCode:    "Active",
		FieldName:    "是否有效",
		FieldType:    1,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       10,
		Active:       3,
	}
	err = indexField.Index(moduleActive.ID, moduleActive)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	contentIndex, err := bleve.New(constant.INDEX_CONTENT_NAME, mapping)
	// 放到IndexMap中
	IndexMap[constant.INDEX_CONTENT_NAME] = contentIndex

	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	return true, nil
}
