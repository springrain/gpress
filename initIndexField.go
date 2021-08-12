package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/mapping"
)

// 数据默认的创建用户
const createUser = "system"

// 全局存放 索引对象,启动之后,所有的索引都通过这个map获取,一个索引只能打开一次,类似数据库连接,用一个对象操作
var IndexMap map[string]bleve.Index = make(map[string]bleve.Index)

//索引名称
const (
	//基本目录
	datadir = "./gpressdatadir/"
	//数据目录,如果不存在认为是第一次安装启动,会创建默认的数据
	indexDataDir = datadir + "index/"
	//索引字段的名称
	indexFieldIndexName = indexDataDir + "IndexField"
	//User 用户的索引名称
	userIndexName = indexDataDir + "User"
	//siteInfo  站点信息
	indexSitenIndexName = indexDataDir + "sitenInfo"
	//页面模板
	indexPageTemplateName = indexDataDir + "pageTemplate"
	//导航菜单
	indexNavMenuName = indexDataDir + "NavMenu"
	//模型
	indexModuleName = indexDataDir + "Module"
	//模型数据
	indexModuleDefaultName = "moduleDefault"
	//文章内容
	indexContentName = indexDataDir + "Content"
)

//逗号分词器的mapping
var commaAnalyzerMapping *mapping.FieldMapping = bleve.NewTextFieldMapping()

//中文分词器的mapping
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
	//BusinessID  业务ID,处理业务记录临时增加的字段,意外情况
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

//初始化调用
func init() {
	//初始化分词器
	commaAnalyzerMapping.DocValues = false
	commaAnalyzerMapping.Analyzer = commaAnalyzerName
	gseAnalyzerMapping.DocValues = false
	gseAnalyzerMapping.Analyzer = gseAnalyzerName

	checkInstall()
}

//checkInstall 检查是不是初始化安装,如果是就创建文件夹目录
func checkInstall() (bool, error) {
	//索引数据目录是否存在
	exists, errPathExists := pathExists(indexDataDir)
	if errPathExists != nil {
		FuncLogError(errPathExists)
		return false, errPathExists
	}

	if exists { //如果已经存在目录,遍历索引,放到全局map里
		fileInfo, _ := ioutil.ReadDir(indexDataDir)
		for _, dir := range fileInfo {
			if !dir.IsDir() {
				continue
			}

			//打开所有的索引,放到map里,一个索引只能打开一次.
			index, _ := bleve.Open(indexDataDir + dir.Name())
			IndexMap[indexDataDir+dir.Name()] = index
		}
		return true, errPathExists
	}
	//如果是初次安装,创建数据目录,默认的 ./gpressdatadir 必须存在,页面模板文件夹 ./gpressdatadir/template
	errMkdir := os.Mkdir(indexDataDir, os.ModePerm)
	if errMkdir != nil {
		FuncLogError(errMkdir)
		return false, errMkdir
	}

	//初始化IndexField
	initIndexField()

	//初始化用户表
	initUser()

	//初始化站点信息表
	initSitenInfo()

	//初始化文章模型的类型表
	initModule()
	//初始化文章默认模型的记录,往Module表里插入记录,不创建Index
	//在IndexField表里设置IndexCode='Module',记录所有的Module.
	//然后在IndexField中插入每个module的字段,每个module实例的ModuleCode都是不同的,使用Module_+后缀的方式命名,只是记录,并不创建index
	initModuleDefault()

	//初始化文章内容
	initContent()

	//初始化导航菜单
	initNavMenu()

	//初始化页面模板
	initpageTemplateName()
	return true, nil
}

// initIndexField 初始化创建IndexField索引
func initIndexField() (bool, error) {
	mapping := bleve.NewIndexMapping()
	//指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name
	index, err := bleve.New(indexFieldIndexName, mapping)

	if err != nil {
		FuncLogError(err)
		return false, err
	}
	IndexMap[indexFieldIndexName] = index
	return true, nil
}

//FuncGenerateStringID 默认生成字符串ID的函数.方便自定义扩展
//FuncGenerateStringID Function to generate string ID by default. Convenient for custom extension
var FuncGenerateStringID func() string = generateStringID

//generateStringID 生成主键字符串
//generateStringID Generate primary key string
func generateStringID() string {

	// 使用 crypto/rand 真随机9位数
	randNum, randErr := rand.Int(rand.Reader, big.NewInt(1000000000))
	if randErr != nil {
		return ""
	}
	//获取9位数,前置补0,确保9位数
	rand9 := fmt.Sprintf("%09d", randNum)

	//获取纳秒 按照 年月日时分秒毫秒微秒纳秒 拼接为长度23位的字符串
	pk := time.Now().Format("2006.01.02.15.04.05.000000000")
	pk = strings.ReplaceAll(pk, ".", "")

	//23位字符串+9位随机数=32位字符串,这样的好处就是可以使用ID进行排序
	pk = pk + rand9
	return pk
}

// pathExists 文件或者目录是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
