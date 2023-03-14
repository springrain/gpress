package bleves

import (
	"gitee.com/gpress/gpress/configs"
	"gitee.com/gpress/gpress/logger"
	"gitee.com/gpress/gpress/util"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"os"
)

// IndexFieldStruct 索引和字段(索引名:IndexField)
// 记录所有索引字段code和中文说明.
// 理论上所有的索引字段都可以放到这个表里,因为都是Map,就不需要再单独指定索引的字段了,可以动态创建Index(目前建议这样做)
// 这个可能是唯一的Struct......

// 初始化 bleve 索引
func checkBleveStatus() bool {
	// 初始化分词器
	commaAnalyzerMapping.DocValues = false
	commaAnalyzerMapping.Analyzer = configs.COMMA_ANALYZER_NAME
	gseAnalyzerMapping.DocValues = false
	gseAnalyzerMapping.Analyzer = configs.GSE_ANGLYZER_NAME

	status, err := checkBleveCreate()
	if err != nil {
		logger.FuncLogError(err)
		return false
	}
	return status
}

// checkBleveCreate 检查是不是初始化安装,如果是就创建文件夹目录
func checkBleveCreate() (bool, error) {
	// 索引数据目录是否存在
	exists, errPathExists := util.PathExists(configs.BLEVE_DATA_DIR)
	if errPathExists != nil {
		logger.FuncLogError(errPathExists)
		return false, errPathExists
	}

	if exists { // 如果已经存在目录,遍历索引,放到全局map里
		fileInfo, _ := os.ReadDir(configs.BLEVE_DATA_DIR)
		for _, dir := range fileInfo {
			if !dir.IsDir() {
				continue
			}

			// 打开所有的索引,放到map里,一个索引只能打开一次.
			index, err := bleve.Open(configs.BLEVE_DATA_DIR + dir.Name())
			if err != nil {
				return false, err
			}
			configs.IndexMap[configs.BLEVE_DATA_DIR+dir.Name()] = index
		}
		return true, errPathExists
	}
	// 如果是初次安装,创建数据目录,默认的 ./gpressdatadir 必须存在,页面模板文件夹 ./gpressdatadir/template
	errMkdir := os.Mkdir(configs.BLEVE_DATA_DIR, os.ModePerm)
	if errMkdir != nil {
		logger.FuncLogError(errMkdir)
		return false, errMkdir
	}

	// 初始化IndexField
	_, err := initIndexField()
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	// 初始化配置
	_, err = initConfig()
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	// 初始化用户表
	_, err = initUser()
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	// 初始化站点信息表
	_, err = initSitenInfo()
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	// 初始化文章模型的类型表
	_, err = initModule()
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	// 初始化文章默认模型的记录,往Module表里插入记录,不创建Index
	// 在IndexField表里设置IndexCode='Module',记录所有的Module.
	// 然后在IndexField中插入每个module的字段,每个module实例的ModuleCode都是不同的,使用Module_+后缀的方式命名,只是记录,并不创建index
	_, err = initModuleDefault()
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	// 初始化文章内容
	_, err = initContent()
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	// 初始化导航菜单
	_, err = initNavMenu()
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}

	// 初始化页面模板
	_, err = initpageTemplateName()
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	return true, nil
}

// initIndexField 初始化创建IndexField索引
func initIndexField() (bool, error) {
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name
	index, err := bleve.New(configs.INDEX_FIELD_INDEX_NAME, mapping)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	configs.IndexMap[configs.INDEX_FIELD_INDEX_NAME] = index
	return true, nil
}
