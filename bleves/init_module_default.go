package bleves

import (
	"gitee.com/gpress/gpress/configs"
	"gitee.com/gpress/gpress/logger"
	"gitee.com/gpress/gpress/util"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"time"
)

func initModuleDefault() (bool, error) {
	indexField := configs.IndexMap[configs.INDEX_FIELD_INDEX_NAME]

	// 获取当前时间
	now := time.Now()

	// 初始化各个字段
	moduleDefaultId := configs.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    configs.INDEX_MODULE_DEFAULT_NAME,
		IndexName:    "模型数据",
		FieldCode:    "ID",
		FieldName:    "模型数据ID",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   configs.CREATE_USER,
		SortNo:       1,
		Active:       3,
	}
	// 放入文件中
	err := indexField.Index(moduleDefaultId.ID, moduleDefaultId)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	moduleDefaultTitle := configs.IndexFieldStruct{
		ID:        util.FuncGenerateStringID(),
		IndexCode: configs.INDEX_MODULE_DEFAULT_NAME,
		IndexName: "模型数据",
		FieldCode: "Title",
		FieldName: "标题",
		FieldType: 3,
		// 文章标题使用中文分词
		AnalyzerName: configs.GSE_ANGLYZER_NAME,
		CreateTime:   now,
		CreateUser:   configs.CREATE_USER,
		SortNo:       2,
		Active:       3,
	}
	err = indexField.Index(moduleDefaultTitle.ID, moduleDefaultTitle)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	moduleDefaultKeyWords := configs.IndexFieldStruct{
		ID:        util.FuncGenerateStringID(),
		IndexCode: configs.INDEX_MODULE_DEFAULT_NAME,
		IndexName: "模型数据",
		FieldCode: "KeyWords",
		FieldName: "关键字",
		FieldType: 3,
		// 文章关键字使用逗号分词器
		AnalyzerName: configs.GSE_ANGLYZER_NAME,
		CreateTime:   now,
		CreateUser:   configs.CREATE_USER,
		SortNo:       3,
		Active:       3,
	}
	err = indexField.Index(moduleDefaultKeyWords.ID, moduleDefaultKeyWords)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	moduleDefaultDescription := configs.IndexFieldStruct{
		ID:        util.FuncGenerateStringID(),
		IndexCode: configs.INDEX_MODULE_DEFAULT_NAME,
		IndexName: "模型数据",
		FieldCode: "Description",
		FieldName: "站点描述",
		FieldType: 3,
		// 文章描述使用中文分词器
		AnalyzerName: configs.GSE_ANGLYZER_NAME,
		CreateTime:   now,
		CreateUser:   configs.CREATE_USER,
		SortNo:       4,
		Active:       3,
	}
	err = indexField.Index(moduleDefaultDescription.ID, moduleDefaultDescription)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	moduleDefaultPageURL := configs.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    configs.INDEX_MODULE_DEFAULT_NAME,
		IndexName:    "模型数据",
		FieldCode:    "PageURL",
		FieldName:    "自身页面路径",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   configs.CREATE_USER,
		SortNo:       5,
		Active:       3,
	}
	err = indexField.Index(moduleDefaultPageURL.ID, moduleDefaultPageURL)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	moduleDefaultSubtitle := configs.IndexFieldStruct{
		ID:        util.FuncGenerateStringID(),
		IndexCode: configs.INDEX_MODULE_DEFAULT_NAME,
		IndexName: "模型数据",
		FieldCode: "Subtitle",
		FieldName: "副标题",
		FieldType: 3,
		// 文章副标题使用中文分词器
		AnalyzerName: configs.GSE_ANGLYZER_NAME,
		CreateTime:   now,
		CreateUser:   configs.CREATE_USER,
		SortNo:       6,
		Active:       3,
	}
	err = indexField.Index(moduleDefaultSubtitle.ID, moduleDefaultSubtitle)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	moduleDefaultContent := configs.IndexFieldStruct{
		ID:        util.FuncGenerateStringID(),
		IndexCode: configs.INDEX_MODULE_DEFAULT_NAME,
		IndexName: "模型数据",
		FieldCode: "Content",
		FieldName: "文章内容",
		FieldType: 3,
		// 文章内容使用中文分词器
		AnalyzerName: configs.GSE_ANGLYZER_NAME,
		CreateTime:   now,
		CreateUser:   configs.CREATE_USER,
		SortNo:       7,
		Active:       3,
	}
	err = indexField.Index(moduleDefaultContent.ID, moduleDefaultContent)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	moduleDefaultCreateTime := configs.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    configs.INDEX_MODULE_DEFAULT_NAME,
		IndexName:    "模型数据",
		FieldCode:    "CreateTime",
		FieldName:    "创建时间",
		FieldType:    2,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   configs.CREATE_USER,
		SortNo:       8,
		Active:       3,
	}
	err = indexField.Index(moduleDefaultCreateTime.ID, moduleDefaultCreateTime)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	moduleUpdateTime := configs.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    configs.INDEX_MODULE_DEFAULT_NAME,
		IndexName:    "模型数据",
		FieldCode:    "UpdateTime",
		FieldName:    "更新时间",
		FieldType:    2,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   configs.CREATE_USER,
		SortNo:       9,
		Active:       3,
	}
	err = indexField.Index(moduleUpdateTime.ID, moduleUpdateTime)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	moduleCreateUser := configs.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    configs.INDEX_MODULE_DEFAULT_NAME,
		IndexName:    "模型数据",
		FieldCode:    "CreateUser",
		FieldName:    "创建人",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   configs.CREATE_USER,
		SortNo:       10,
		Active:       3,
	}
	err = indexField.Index(moduleCreateUser.ID, moduleCreateUser)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	moduleSortNo := configs.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    configs.INDEX_MODULE_DEFAULT_NAME,
		IndexName:    "模型数据",
		FieldCode:    "SortNo",
		FieldName:    "排序",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   configs.CREATE_USER,
		SortNo:       9,
		Active:       3,
	}
	err = indexField.Index(moduleSortNo.ID, moduleSortNo)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	moduleActive := configs.IndexFieldStruct{
		ID:           util.FuncGenerateStringID(),
		IndexCode:    configs.INDEX_MODULE_DEFAULT_NAME,
		IndexName:    "模型数据",
		FieldCode:    "Active",
		FieldName:    "是否有效",
		FieldType:    1,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   configs.CREATE_USER,
		SortNo:       10,
		Active:       3,
	}
	err = indexField.Index(moduleActive.ID, moduleActive)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	/*
		//在IndexField表里设置IndexCode='Module',记录所有的Module.
		//然后在IndexField中插入每个module的字段,每个module实例的ModuleCode都是不同的,使用Module_+后缀的方式命名,只是记录,并不创建index

		//创建用户表的索引
		mapping := bleve.NewIndexMapping()
		//指定默认的分词器
		mapping.DefaultMapping.DefaultAnalyzer = keyword.Name
		moduleDefaultIndex, err := bleve.New(moduleDefaultName, mapping)
		//放到IndexMap中
		IndexMap[moduleName] = moduleDefaultIndex

		if err != nil {
			FuncLogError(err)
			return false, err
		}
	*/

	return true, nil
}
