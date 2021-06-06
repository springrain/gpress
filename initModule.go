package main

import (
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
)

func initModule() (bool, error) {

	indexField := IndexMap[indexFieldIndexName]

	//获取当前时间
	now := time.Now()

	//初始化各个字段
	moduleId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexModuleName,
		IndexName:    "模型",
		FieldCode:    "ID",
		FieldName:    "模型ID",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       1,
		Active:       3,
	}
	//放入文件中
	indexField.Index(moduleId.ID, moduleId)

	moduleModuleIndexCode := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexModuleName,
		IndexName:    "模型",
		FieldCode:    "ModuleIndexCode",
		FieldName:    "模型Code",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       2,
		Active:       3,
	}
	indexField.Index(moduleModuleIndexCode.ID, moduleModuleIndexCode)

	moduleModuleName := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexModuleName,
		IndexName:    "模型",
		FieldCode:    "ModuleName",
		FieldName:    "模型名称",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       3,
		Active:       3,
	}
	indexField.Index(moduleModuleName.ID, moduleModuleName)

	moduleCreateTime := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexModuleName,
		IndexName:    "模型",
		FieldCode:    "CreateTime",
		FieldName:    "创建时间",
		FieldType:    2,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       4,
		Active:       3,
	}
	indexField.Index(moduleCreateTime.ID, moduleCreateTime)

	moduleUpdateTime := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexModuleName,
		IndexName:    "模型",
		FieldCode:    "UpdateTime",
		FieldName:    "更新时间",
		FieldType:    2,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       5,
		Active:       3,
	}
	indexField.Index(moduleUpdateTime.ID, moduleUpdateTime)

	moduleCreateUser := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexModuleName,
		IndexName:    "模型",
		FieldCode:    "CreateUser",
		FieldName:    "创建人",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       6,
		Active:       3,
	}
	indexField.Index(moduleCreateUser.ID, moduleCreateUser)

	moduleSortNo := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexModuleName,
		IndexName:    "模型",
		FieldCode:    "SortNo",
		FieldName:    "排序",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       7,
		Active:       3,
	}
	indexField.Index(moduleSortNo.ID, moduleSortNo)

	moduleActive := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    indexModuleName,
		IndexName:    "模型",
		FieldCode:    "Active",
		FieldName:    "是否有效",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       8,
		Active:       3,
	}
	indexField.Index(moduleActive.ID, moduleActive)

	//创建用户表的索引
	mapping := bleve.NewIndexMapping()
	//指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name
	moduleIndex, err := bleve.New(indexModuleName, mapping)
	//放到IndexMap中
	IndexMap[indexModuleName] = moduleIndex

	//初始化数据
	module := make(map[string]interface{})
	id := FuncGenerateStringID()
	module["ID"] = id
	module["ModuleIndexCode"] = "Module_Default"
	module["ModuleName"] = "默认模型"
	module["CreateTime"] = now
	module["UpdateTime"] = now
	module["CreateUser"] = createUser
	module["SortNo"] = 0
	module["Active"] = 1

	//初始化
	moduleIndex.Index(id, module)

	if err != nil {
		FuncLogError(err)
		return false, err
	}

	return true, nil

}
