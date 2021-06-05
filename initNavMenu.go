package main

import (
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
)

//导航菜单
func initNavMenu() (bool, error) {

	indexField := IndexMap[indexFieldIndexName]

	//获取当前时间
	now := time.Now()

	//初始化各个字段
	navMenuId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    navMenuName,
		IndexName:    "导航菜单",
		FieldCode:    "ID",
		FieldName:    "导航菜单ID",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       1,
		Active:       3,
	}
	//放入文件中
	indexField.Index(navMenuId.ID, navMenuId)

	navMenuMenuName := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    navMenuName,
		IndexName:    "导航菜单",
		FieldCode:    "MenuName",
		FieldName:    "菜单名称",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       2,
		Active:       3,
	}
	indexField.Index(navMenuMenuName.ID, navMenuMenuName)

	navMenuHrefURL := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    navMenuName,
		IndexName:    "导航菜单",
		FieldCode:    "HrefURL",
		FieldName:    "跳转路径",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       3,
		Active:       3,
	}
	indexField.Index(navMenuHrefURL.ID, navMenuHrefURL)

	navMenuHrefTarget := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    navMenuName,
		IndexName:    "导航菜单",
		FieldCode:    "HrefTarget",
		FieldName:    "跳转方式",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       4,
		Active:       3,
	}
	indexField.Index(navMenuHrefTarget.ID, navMenuHrefTarget)

	navMenuPID := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    navMenuName,
		IndexName:    "导航菜单",
		FieldCode:    "PID",
		FieldName:    "父菜单ID",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       5,
		Active:       3,
	}
	indexField.Index(navMenuPID.ID, navMenuPID)

	navMenuThemePC := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    navMenuName,
		IndexName:    "导航菜单",
		FieldCode:    "themePC",
		FieldName:    "PC主题",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       6,
		Active:       3,
	}
	indexField.Index(navMenuThemePC.ID, navMenuThemePC)

	navMenuModuleIndexCode := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    navMenuName,
		IndexName:    "导航菜单",
		FieldCode:    "ModuleIndexCode",
		FieldName:    "Module的索引名称",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       7,
		Active:       3,
	}
	indexField.Index(navMenuModuleIndexCode.ID, navMenuModuleIndexCode)

	navMenuComCode := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    navMenuName,
		IndexName:    "导航菜单",
		FieldCode:    "ComCode",
		FieldName:    "逗号隔开的全路径",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       8,
		Active:       3,
	}
	indexField.Index(navMenuComCode.ID, navMenuComCode)

	navMenuTemplateID := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    navMenuName,
		IndexName:    "导航菜单",
		FieldCode:    "TemplateID",
		FieldName:    "模板Id",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       9,
		Active:       3,
	}
	indexField.Index(navMenuTemplateID.ID, navMenuTemplateID)

	navMenuChildTemplateID := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    navMenuName,
		IndexName:    "导航菜单",
		FieldCode:    "ChildTemplateID",
		FieldName:    "子页面模板Id",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       10,
		Active:       3,
	}
	indexField.Index(navMenuChildTemplateID.ID, navMenuChildTemplateID)

	navMenuSortNo := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    navMenuName,
		IndexName:    "导航菜单",
		FieldCode:    "SortNo",
		FieldName:    "排序",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       11,
		Active:       3,
	}
	indexField.Index(navMenuSortNo.ID, navMenuSortNo)

	navMenuActive := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    navMenuName,
		IndexName:    "导航菜单",
		FieldCode:    "Active",
		FieldName:    "是否有效",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       12,
		Active:       3,
	}
	indexField.Index(navMenuActive.ID, navMenuActive)

	//创建用户表的索引
	mapping := bleve.NewIndexMapping()
	//指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name
	navMenuIndex, err := bleve.New(navMenuName, mapping)

	//放到IndexMap中
	IndexMap[navMenuName] = navMenuIndex

	if err != nil {
		FuncLogError(err)
		return false, err
	}

	return true, nil
}
