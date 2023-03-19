package bleves

import (
	"errors"
	"gitee.com/gpress/gpress/constant"
	"gitee.com/gpress/gpress/logger"
	"gitee.com/gpress/gpress/model"
	"gitee.com/gpress/gpress/util"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
)

// 初始化站点信息
func initSiteInfo() (bool, error) {
	indexField := IndexMap[constant.INDEX_FIELD_INDEX_NAME]

	// 创建用户表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name
	// 获取当前时间
	now := time.Now()

	// 初始化各个字段
	siteInfoId := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_SITE_INDEX_NAME,
		IndexName:    "站点信息",
		FieldCode:    "ID",
		FieldName:    "站点信息ID",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       1,
		Active:       3,
	}
	// 放入文件中
	err := indexField.Index(siteInfoId.ID, siteInfoId)
	if err != nil {
		logger.FuncLogError(err)
		logger.FuncLogError(errors.New("初始化sitenInfoId失败"))
		return false, err
	}

	siteInfoTitle := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_SITE_INDEX_NAME,
		IndexName:    "站点信息",
		FieldCode:    "Title",
		FieldName:    "标题",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       2,
		Active:       3,
	}
	err = indexField.Index(siteInfoTitle.ID, siteInfoTitle)
	if err != nil {
		logger.FuncLogError(err)
		logger.FuncLogError(errors.New("初始化sitenInfoTitle失败"))
		return false, err
	}

	sitenInfoKeyWords := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_SITE_INDEX_NAME,
		IndexName:    "站点信息",
		FieldCode:    "KeyWords",
		FieldName:    "关键字",
		FieldType:    3,
		AnalyzerName: constant.COMMA_ANALYZER_NAME,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       3,
		Active:       3,
	}
	err = indexField.Index(sitenInfoKeyWords.ID, sitenInfoKeyWords)
	if err != nil {
		logger.FuncLogError(err)
		logger.FuncLogError(errors.New("初始化sitenInfoKeyWords失败"))
		return false, err
	}

	// KeyWords 字段使用 逗号分词器的mapping commaAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("KeyWords", commaAnalyzerMapping)

	sitenInfoDescription := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_SITE_INDEX_NAME,
		IndexName:    "站点信息",
		FieldCode:    "Description",
		FieldName:    "站点描述",
		FieldType:    3,
		AnalyzerName: constant.GSE_ANGLYZER_NAME,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       4,
		Active:       3,
	}
	err = indexField.Index(sitenInfoDescription.ID, sitenInfoDescription)
	if err != nil {
		logger.FuncLogError(err)
		logger.FuncLogError(errors.New("初始化siteInfoDescription失败"))
		return false, err
	}

	// Description 字段使用 中文分词器的mapping gseAnalyzerMapping
	mapping.DefaultMapping.AddFieldMappingsAt("Description", gseAnalyzerMapping)

	siteInfoTheme := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_SITE_INDEX_NAME,
		IndexName:    "站点信息",
		FieldCode:    "theme",
		FieldName:    "默认主题",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       5,
		Active:       3,
	}
	err = indexField.Index(siteInfoTheme.ID, siteInfoTheme)
	if err != nil {
		logger.FuncLogError(err)
		logger.FuncLogError(errors.New("初始化siteInfoTheme失败"))
		return false, err
	}

	siteInfoThemePC := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_SITE_INDEX_NAME,
		IndexName:    "站点信息",
		FieldCode:    "themePC",
		FieldName:    "PC主题",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       6,
		Active:       3,
	}
	err = indexField.Index(siteInfoThemePC.ID, siteInfoThemePC)
	if err != nil {
		logger.FuncLogError(err)
		logger.FuncLogError(errors.New("初始化sitenInfoThemePC失败"))
		return false, err
	}

	siteInfoThemeWAP := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_SITE_INDEX_NAME,
		IndexName:    "站点信息",
		FieldCode:    "themeWAP",
		FieldName:    "手机主题",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       7,
		Active:       3,
	}
	err = indexField.Index(siteInfoThemeWAP.ID, siteInfoThemeWAP)
	if err != nil {
		logger.FuncLogError(err)
		logger.FuncLogError(errors.New("初始化sitenInfoThemeWAP失败"))
		return false, err
	}

	siteInfoThemeWEIXIN := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_SITE_INDEX_NAME,
		IndexName:    "站点信息",
		FieldCode:    "sitenInfoThemeWEIXIN",
		FieldName:    "微信主题",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       8,
		Active:       3,
	}
	err = indexField.Index(siteInfoThemeWEIXIN.ID, siteInfoThemeWEIXIN)
	if err != nil {
		logger.FuncLogError(err)
		logger.FuncLogError(errors.New("初始化sitenInfoThemeWEIXIN失败"))
		return false, err
	}

	siteInfoLogo := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_SITE_INDEX_NAME,
		IndexName:    "站点信息",
		FieldCode:    "Logo",
		FieldName:    "Logo",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       9,
		Active:       3,
	}
	err = indexField.Index(siteInfoLogo.ID, siteInfoLogo)
	if err != nil {
		logger.FuncLogError(err)
		logger.FuncLogError(errors.New("初始化sitenInfoLogo失败"))
		return false, err
	}

	siteInfoFavicon := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_SITE_INDEX_NAME,
		IndexName:    "站点信息",
		FieldCode:    "Favicon",
		FieldName:    "Favicon",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       10,
		Active:       3,
	}
	err = indexField.Index(siteInfoFavicon.ID, siteInfoFavicon)
	if err != nil {
		logger.FuncLogError(err)
		logger.FuncLogError(errors.New("初始化siteInfoLogo失败"))
		return false, err
	}

	siteIndexIndex, err := bleve.New(constant.INDEX_SITE_INDEX_NAME, mapping)

	if err != nil {
		logger.FuncLogError(err)
		logger.FuncLogError(errors.New("初始化bleve.New (constant.INDEX_SITE_INDEX_NAME)失败"))
		return false, err
	}
	// 放到IndexMap中comma
	IndexMap[constant.INDEX_SITE_INDEX_NAME] = siteIndexIndex
	return true, nil
}
