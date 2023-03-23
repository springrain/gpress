package gbleve

import (
	"gitee.com/gpress/gpress/constant"
	"gitee.com/gpress/gpress/logger"
	"gitee.com/gpress/gpress/model"
	"gitee.com/gpress/gpress/util"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"time"
)

// 导航菜单
func initNavMenu() (bool, error) {
	indexField := IndexMap[constant.INDEX_FIELD_INDEX_NAME]

	// 获取当前时间
	now := time.Now()

	// 初始化各个字段
	navMenuId := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_NAV_MENU_NAME,
		IndexName:    "导航菜单",
		FieldCode:    "ID",
		FieldName:    "导航菜单ID",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       1,
		Active:       3,
		Required:     1,
	}
	// 放入文件中
	err := indexField.Index(navMenuId.ID, navMenuId)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	navMenuMenuName := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_NAV_MENU_NAME,
		IndexName:    "导航菜单",
		FieldCode:    "MenuName",
		FieldName:    "菜单名称",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       2,
		Active:       3,
		Required:     1,
	}
	err = indexField.Index(navMenuMenuName.ID, navMenuMenuName)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	navMenuHrefURL := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_NAV_MENU_NAME,
		IndexName:    "导航菜单",
		FieldCode:    "HrefURL",
		FieldName:    "跳转路径",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       3,
		Active:       3,
		Required:     0,
	}
	err = indexField.Index(navMenuHrefURL.ID, navMenuHrefURL)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	navMenuHrefTarget := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_NAV_MENU_NAME,
		IndexName:    "导航菜单",
		FieldCode:    "HrefTarget",
		FieldName:    "跳转方式",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       4,
		Active:       3,
		Required:     0,
	}
	err = indexField.Index(navMenuHrefTarget.ID, navMenuHrefTarget)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	navMenuPID := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_NAV_MENU_NAME,
		IndexName:    "导航菜单",
		FieldCode:    "PID",
		FieldName:    "父菜单ID",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       5,
		Active:       3,
		Required:     0,
	}
	err = indexField.Index(navMenuPID.ID, navMenuPID)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	navMenuThemePC := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_NAV_MENU_NAME,
		IndexName:    "导航菜单",
		FieldCode:    "themePC",
		FieldName:    "PC主题",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       6,
		Active:       3,
		Required:     0,
	}
	err = indexField.Index(navMenuThemePC.ID, navMenuThemePC)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	navMenuModuleIndexCode := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_NAV_MENU_NAME,
		IndexName:    "导航菜单",
		FieldCode:    "ModuleIndexCode",
		FieldName:    "Module的索引名称",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       7,
		Active:       3,
		Required:     0,
	}
	err = indexField.Index(navMenuModuleIndexCode.ID, navMenuModuleIndexCode)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	navMenuComCode := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_NAV_MENU_NAME,
		IndexName:    "导航菜单",
		FieldCode:    "ComCode",
		FieldName:    "逗号隔开的全路径",
		FieldType:    3,
		AnalyzerName: constant.COMMA_ANALYZER_NAME,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       8,
		Active:       3,
		Required:     0,
	}
	err = indexField.Index(navMenuComCode.ID, navMenuComCode)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	navMenuTemplateID := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_NAV_MENU_NAME,
		IndexName:    "导航菜单",
		FieldCode:    "TemplateID",
		FieldName:    "模板Id",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       9,
		Active:       3,
		Required:     0,
	}
	err = indexField.Index(navMenuTemplateID.ID, navMenuTemplateID)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	navMenuChildTemplateID := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_NAV_MENU_NAME,
		IndexName:    "导航菜单",
		FieldCode:    "ChildTemplateID",
		FieldName:    "子页面模板Id",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       10,
		Active:       3,
		Required:     0,
	}
	err = indexField.Index(navMenuChildTemplateID.ID, navMenuChildTemplateID)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	navMenuSortNo := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_NAV_MENU_NAME,
		IndexName:    "导航菜单",
		FieldCode:    "SortNo",
		FieldName:    "排序",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       11,
		Active:       3,
		Required:     0,
	}
	err = indexField.Index(navMenuSortNo.ID, navMenuSortNo)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	navMenuActive := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_NAV_MENU_NAME,
		IndexName:    "导航菜单",
		FieldCode:    "Active",
		FieldName:    "是否有效",
		FieldType:    1,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       12,
		Active:       3,
		Required:     0,
	}
	err = indexField.Index(navMenuActive.ID, navMenuActive)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	// 创建用户表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name
	navMenuIndex, err := bleve.New(constant.INDEX_NAV_MENU_NAME, mapping)

	// 放到IndexMap中
	IndexMap[constant.INDEX_NAV_MENU_NAME] = navMenuIndex

	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	return true, nil
}
