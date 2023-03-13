package bleves

import (
	"gitee.com/gpress/gpress/config"
	"gitee.com/gpress/gpress/logger"
	"gitee.com/gpress/gpress/util"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"time"
)

// 初始化站点信息
func initSitenInfo() (bool, error) {
	indexField := config.IndexMap[config.INDEX_FIELD_INDEX_NAME]

	// 创建用户表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name
	// 获取当前时间
	now := time.Now()

	// 初始化各个字段
	sitenInfoId := config.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    config.INDEX_SITE_INDEX_NAME,
		IndexName:    "站点信息",
		FieldCode:    "ID",
		FieldName:    "站点信息ID",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   config.CREATE_USER,
		SortNo:       1,
		Active:       3,
	}
	// 放入文件中
	err := indexField.Index(sitenInfoId.ID, sitenInfoId)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	sitenInfoTitle := config.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    config.INDEX_SITE_INDEX_NAME,
		IndexName:    "站点信息",
		FieldCode:    "Title",
		FieldName:    "标题",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   config.CREATE_USER,
		SortNo:       2,
		Active:       3,
	}
	err = indexField.Index(sitenInfoTitle.ID, sitenInfoTitle)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	sitenInfoKeyWords := config.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    config.INDEX_SITE_INDEX_NAME,
		IndexName:    "站点信息",
		FieldCode:    "KeyWords",
		FieldName:    "关键字",
		FieldType:    3,
		AnalyzerName: config.COMMA_ANALYZER_NAME,
		CreateTime:   now,
		CreateUser:   config.CREATE_USER,
		SortNo:       3,
		Active:       3,
	}
	err = indexField.Index(sitenInfoKeyWords.ID, sitenInfoKeyWords)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	// KeyWords 字段使用 逗号分词器的mapping commaAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("KeyWords", commaAnalyzerMapping)

	sitenInfoDescription := config.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    config.INDEX_SITE_INDEX_NAME,
		IndexName:    "站点信息",
		FieldCode:    "Description",
		FieldName:    "站点描述",
		FieldType:    3,
		AnalyzerName: config.GSE_ANGLYZER_NAME,
		CreateTime:   now,
		CreateUser:   config.CREATE_USER,
		SortNo:       4,
		Active:       3,
	}
	err = indexField.Index(sitenInfoDescription.ID, sitenInfoDescription)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	// Description 字段使用 中文分词器的mapping gseAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("Description", gseAnalyzerMapping)

	sitenInfoTheme := config.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    config.INDEX_SITE_INDEX_NAME,
		IndexName:    "站点信息",
		FieldCode:    "theme",
		FieldName:    "默认主题",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   config.CREATE_USER,
		SortNo:       5,
		Active:       3,
	}
	err = indexField.Index(sitenInfoTheme.ID, sitenInfoTheme)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	sitenInfoThemePC := config.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    config.INDEX_SITE_INDEX_NAME,
		IndexName:    "站点信息",
		FieldCode:    "themePC",
		FieldName:    "PC主题",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   config.CREATE_USER,
		SortNo:       6,
		Active:       3,
	}
	err = indexField.Index(sitenInfoThemePC.ID, sitenInfoThemePC)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	sitenInfoThemeWAP := config.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    config.INDEX_SITE_INDEX_NAME,
		IndexName:    "站点信息",
		FieldCode:    "themeWAP",
		FieldName:    "手机主题",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   config.CREATE_USER,
		SortNo:       7,
		Active:       3,
	}
	err = indexField.Index(sitenInfoThemeWAP.ID, sitenInfoThemeWAP)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	sitenInfoThemeWEIXIN := config.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    config.INDEX_SITE_INDEX_NAME,
		IndexName:    "站点信息",
		FieldCode:    "sitenInfoThemeWEIXIN",
		FieldName:    "微信主题",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   config.CREATE_USER,
		SortNo:       8,
		Active:       3,
	}
	err = indexField.Index(sitenInfoThemeWEIXIN.ID, sitenInfoThemeWEIXIN)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	sitenInfoLogo := config.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    config.INDEX_SITE_INDEX_NAME,
		IndexName:    "站点信息",
		FieldCode:    "Logo",
		FieldName:    "Logo",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   config.CREATE_USER,
		SortNo:       9,
		Active:       3,
	}
	err = indexField.Index(sitenInfoLogo.ID, sitenInfoLogo)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	sitenInfoFavicon := config.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    config.INDEX_SITE_INDEX_NAME,
		IndexName:    "站点信息",
		FieldCode:    "Favicon",
		FieldName:    "Favicon",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   config.CREATE_USER,
		SortNo:       10,
		Active:       3,
	}
	err = indexField.Index(sitenInfoFavicon.ID, sitenInfoFavicon)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	sitenIndexIndex, err := bleve.New(config.INDEX_SITE_INDEX_NAME, mapping)
	// 放到config.IndexMap中
	config.IndexMap[config.INDEX_SITE_INDEX_NAME] = sitenIndexIndex

	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	return true, nil
}
