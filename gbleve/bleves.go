package gbleve

import (
	"errors"
	"fmt"
	"gitee.com/gpress/gpress/constant"
	"gitee.com/gpress/gpress/logger"
	"gitee.com/gpress/gpress/util"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"os"
)

var BleveStatus bool
var Installed bool

// IndexMap 全局存放 索引对象,启动之后,所有的索引都通过这个map获取,一个索引只能打开一次,类似数据库连接,用一个对象操作
var IndexMap = make(map[string]bleve.Index)

// 中文分词器的mapping
var gseAnalyzerMapping = bleve.NewTextFieldMapping()

// 逗号分词器的mapping
var commaAnalyzerMapping = bleve.NewTextFieldMapping()

// CheckBleveStatus 初始化 bleve 索引
func CheckBleveStatus() bool {
	// 初始化分词器
	commaAnalyzerMapping.DocValues = false
	commaAnalyzerMapping.Analyzer = constant.COMMA_ANALYZER_NAME
	gseAnalyzerMapping.DocValues = false
	gseAnalyzerMapping.Analyzer = constant.GSE_ANGLYZER_NAME

	status, err := checkBleveCreate()
	if err != nil {
		logger.FuncLogError(err)
		return false
	}
	return status
}
func IsInstalled() bool {
	// 依赖bleveStatus变量,确保bleve在isInstalled之前初始化
	if !BleveStatus {
		logger.FuncLogError(errors.New("bleveStatus状态为false"))
	}
	_, err := os.Lstat(constant.TEMPLATE_DIR + "admin/install.html")
	return os.IsNotExist(err)
}

// checkBleveCreate 检查是不是初始化安装,如果是就创建文件夹目录
func checkBleveCreate() (bool, error) {
	// 索引数据目录是否存在
	exists, errPathExists := util.PathExists(constant.BLEVE_DATA_DIR)
	fmt.Println(errPathExists)
	if errPathExists != nil {
		logger.FuncLogError(errPathExists)
		return false, errPathExists
	}

	if exists { // 如果已经存在目录,遍历索引,放到全局map里
		fileInfo, _ := os.ReadDir(constant.BLEVE_DATA_DIR)
		for _, dir := range fileInfo {
			if !dir.IsDir() {
				continue
			}

			// 打开所有的索引,放到map里,一个索引只能打开一次.
			index, err := bleve.Open(constant.BLEVE_DATA_DIR + dir.Name())
			if err != nil {
				return false, err
			}
			IndexMap[constant.BLEVE_DATA_DIR+dir.Name()] = index
		}
		return true, errPathExists
	}
	// 如果是初次安装,创建数据目录,默认的 ./gpressdatadir 必须存在,页面模板文件夹 ./gpressdatadir/template
	errMkdir := os.Mkdir(constant.BLEVE_DATA_DIR, os.ModePerm)
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
		logger.FuncLogError(errors.New("初始化Config失败"))
		return false, err
	}

	// 初始化用户表
	_, err = initUser()
	if err != nil {
		logger.FuncLogError(errors.New("初始化User失败"))
		return false, err
	}

	// 初始化站点信息表
	_, err = initSiteInfo()
	if err != nil {
		logger.FuncLogError(errors.New("初始化siteInfo失败"))
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
		logger.FuncLogError(errors.New("初始化Module失败"))
		return false, err
	}

	// 初始化文章内容
	_, err = initContent()

	// 初始化导航菜单
	_, err = initNavMenu()
	if err != nil {
		logger.FuncLogError(errors.New("初始化NavMenu失败"))
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
	index, err := bleve.New(constant.INDEX_FIELD_INDEX_NAME, mapping)
	if err != nil {
		logger.FuncLogError(err)
		return false, err
	}
	IndexMap[constant.INDEX_FIELD_INDEX_NAME] = index
	return true, nil
}
