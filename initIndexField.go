package main

import (
	"os"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/mapping"
)

// 全局存放 索引对象,启动之后,所有的索引都通过这个map获取,一个索引只能打开一次,类似数据库连接,用一个对象操作
var IndexMap map[string]bleve.Index = make(map[string]bleve.Index)

// 逗号分词器的mapping
var commaAnalyzerMapping *mapping.FieldMapping = bleve.NewTextFieldMapping()

// 中文分词器的mapping
var gseAnalyzerMapping *mapping.FieldMapping = bleve.NewTextFieldMapping()

// IndexFieldStruct 索引和字段(索引名:IndexField)
// 记录所有索引字段code和中文说明.
// 理论上所有的索引字段都可以放到这个表里,因为都是Map,就不需要再单独指定索引的字段了,可以动态创建Index(目前建议这样做)
// 这个可能是唯一的Struct......
type IndexFieldStruct struct {
	// ID 主键
	ID string
	// IndexCode 索引代码,类似表名 User,SiteInfo,PageTemplate,NavMenu,Module,Content
	IndexCode string
	// IndexCode 索引名称,类似表名中文说明
	IndexName string
	// BusinessID  业务ID,处理业务记录临时增加的字段,意外情况
	BusinessID string
	// FieldCode  字段代码
	FieldCode string
	// FieldName  字段中文名称
	FieldName string
	// FieldType  字段类型,数字(1),日期(2),文本框(3),文本域(4),富文本(5),下拉框(6),单选(7),多选(8),上传图片(9),上传附件(10),轮播图(11),音频(12),视频(13)
	FieldType int
	// FieldFormat 数据格式,用于日期或者数字
	FieldFormat string
	// Required 字段是否为空. 0可以为空,1必填
	Required int
	// DefaultValue  默认值
	DefaultValue string
	// AnalyzerName  分词器名称
	AnalyzerName string
	// CreateTime 创建时间
	CreateTime time.Time
	// CreateTime 更新时间
	UpdateTime time.Time
	// AnalyzerName  创建人,初始化 system
	CreateUser string
	// SortNo 排序
	SortNo int
	// 是否有效 无效(0),正常显示(1),界面不显示(3)
	Active int
}

// 初始化 bleve 索引
func checkBleveStatus() bool {
	// 初始化分词器
	commaAnalyzerMapping.DocValues = false
	commaAnalyzerMapping.Analyzer = commaAnalyzerName
	gseAnalyzerMapping.DocValues = false
	gseAnalyzerMapping.Analyzer = gseAnalyzerName

	status, err := checkBleveCreate()
	if err != nil {
		FuncLogError(err)
		return false
	}
	return status
}

// checkBleveCreate 检查是不是初始化安装,如果是就创建文件夹目录
func checkBleveCreate() (bool, error) {
	// 索引数据目录是否存在
	exists, errPathExists := pathExists(bleveDataDir)
	if errPathExists != nil {
		FuncLogError(errPathExists)
		return false, errPathExists
	}

	if exists { // 如果已经存在目录,遍历索引,放到全局map里
		fileInfo, _ := os.ReadDir(bleveDataDir)
		for _, dir := range fileInfo {
			if !dir.IsDir() {
				continue
			}

			// 打开所有的索引,放到map里,一个索引只能打开一次.
			index, err := bleve.Open(bleveDataDir + dir.Name())
			if err != nil {
				return false, err
			}
			IndexMap[bleveDataDir+dir.Name()] = index
		}
		return true, errPathExists
	}
	// 如果是初次安装,创建数据目录,默认的 ./gpressdatadir 必须存在,页面模板文件夹 ./gpressdatadir/template
	errMkdir := os.Mkdir(bleveDataDir, os.ModePerm)
	if errMkdir != nil {
		FuncLogError(errMkdir)
		return false, errMkdir
	}

	// 初始化IndexField
	initIndexField()

	// 初始化配置
	initConfig()
	// 初始化用户表
	initUser()

	// 初始化站点信息表
	initSitenInfo()

	// 初始化文章模型的类型表
	initModule()
	// 初始化文章默认模型的记录,往Module表里插入记录,不创建Index
	// 在IndexField表里设置IndexCode='Module',记录所有的Module.
	// 然后在IndexField中插入每个module的字段,每个module实例的ModuleCode都是不同的,使用Module_+后缀的方式命名,只是记录,并不创建index
	initModuleDefault()

	// 初始化文章内容
	initContent()

	// 初始化导航菜单
	initNavMenu()

	// 初始化页面模板
	initpageTemplateName()
	return true, nil
}

// initIndexField 初始化创建IndexField索引
func initIndexField() (bool, error) {
	mapping := bleve.NewIndexMapping()
	// 指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name
	index, err := bleve.New(indexFieldIndexName, mapping)
	if err != nil {
		FuncLogError(err)
		return false, err
	}
	IndexMap[indexFieldIndexName] = index
	return true, nil
}
