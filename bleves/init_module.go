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

func initModule() (bool, error) {
	indexField := IndexMap[constant.INDEX_FIELD_INDEX_NAME]

	// 获取当前时间
	now := time.Now()

	// 初始化各个字段
	moduleId := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_MODULE_NAME,
		IndexName:    "模型",
		FieldCode:    "ID",
		FieldName:    "模型ID",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       1,
		Active:       3,
	}
	// 放入文件中
	err := indexField.Index(moduleId.ID, moduleId)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	moduleModuleIndexCode := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_MODULE_NAME,
		IndexName:    "模型",
		FieldCode:    "ModuleIndexCode",
		FieldName:    "模型Code",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       2,
		Active:       3,
	}
	err = indexField.Index(moduleModuleIndexCode.ID, moduleModuleIndexCode)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	moduleModuleName := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_MODULE_NAME,
		IndexName:    "模型",
		FieldCode:    "ModuleName",
		FieldName:    "模型名称",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       3,
		Active:       3,
	}
	err = indexField.Index(moduleModuleName.ID, moduleModuleName)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	moduleCreateTime := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_MODULE_NAME,
		IndexName:    "模型",
		FieldCode:    "CreateTime",
		FieldName:    "创建时间",
		FieldType:    2,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       4,
		Active:       3,
	}
	err = indexField.Index(moduleCreateTime.ID, moduleCreateTime)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	moduleUpdateTime := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_MODULE_NAME,
		IndexName:    "模型",
		FieldCode:    "UpdateTime",
		FieldName:    "更新时间",
		FieldType:    2,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       5,
		Active:       3,
	}
	err = indexField.Index(moduleUpdateTime.ID, moduleUpdateTime)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	moduleCreateUser := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_MODULE_NAME,
		IndexName:    "模型",
		FieldCode:    "CreateUser",
		FieldName:    "创建人",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       6,
		Active:       3,
	}
	err = indexField.Index(moduleCreateUser.ID, moduleCreateUser)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	moduleSortNo := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_MODULE_NAME,
		IndexName:    "模型",
		FieldCode:    "SortNo",
		FieldName:    "排序",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       7,
		Active:       3,
	}
	err = indexField.Index(moduleSortNo.ID, moduleSortNo)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	moduleActive := model.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    constant.INDEX_MODULE_NAME,
		IndexName:    "模型",
		FieldCode:    "Active",
		FieldName:    "是否有效",
		FieldType:    1,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   constant.CREATE_USER,
		SortNo:       8,
		Active:       3,
	}
	err = indexField.Index(moduleActive.ID, moduleActive)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	// 创建用户表的索引
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name
	moduleIndex, err := bleve.New(constant.INDEX_MODULE_NAME, mapping)
	// 放到IndexMap中
	IndexMap[constant.INDEX_MODULE_NAME] = moduleIndex

	// 初始化数据
	module := make(map[string]interface{})
	id := util.FuncGenerateStringID()
	module["ID"] = id
	module["ModuleIndexCode"] = "Module_Default"
	module["ModuleName"] = "默认模型"
	module["CreateTime"] = now
	module["UpdateTime"] = now
	module["CreateUser"] = constant.CREATE_USER
	module["SortNo"] = 0
	module["Active"] = 1

	// 初始化
	err = moduleIndex.Index(id, module)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	return true, nil
}
