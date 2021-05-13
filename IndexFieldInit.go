package main

import (
	"os"
	"time"

	"gitee.com/chunanyong/gouuid"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
)

//数据目录
const indexDataDir = "./zcmsdatadir/index/"
const createUser = "system"

// 存放 索引对象
var bleveIndexMap map[string]bleve.Index = make(map[string]bleve.Index)

//indexField的索引名称
const indexFieldIndexName = indexDataDir + "IndexField"
const userIndexName = indexDataDir + "User"

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

//checkInstall 检查是不是初始化安装,如果是就创建文件夹目录
func checkInstall() (bool, error) {
	exists, errPathExists := pathExists(indexDataDir)
	if errPathExists != nil {
		FuncLogError(errPathExists)
		return false, errPathExists
	}

	if exists {
		return true, errPathExists
	}
	//如果是初次安装,创建数据目录,默认的 ./zcmsdatadir 必须存在,和打包的二进制文件放到同一个路径下,里面有页面模板文件夹 ./zcmsdatadir/template
	errMkdir := os.Mkdir(indexDataDir, os.ModePerm)
	if errMkdir != nil {
		FuncLogError(errMkdir)
		return false, errMkdir
	}

	//初始化IndexField
	initIndexField()

	//初始化用户表
	initUser()

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
	bleveIndexMap[indexFieldIndexName] = index
	return true, nil
}

// initIndexField 初始化创建IndexField索引
func initUser() (bool, error) {

	//设置IndexField数据
	index := bleveIndexMap[indexFieldIndexName]

	now := time.Now()
	userId := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    userIndexName,
		IndexName:    "用户信息",
		FieldCode:    "ID",
		FieldName:    "用户ID",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       1,
		Active:       3,
	}
	index.Index(userId.ID, userId)

	userAccount := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    userIndexName,
		IndexName:    "用户信息",
		FieldCode:    "Account",
		FieldName:    "账号",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       2,
		Active:       1,
	}
	index.Index(userAccount.ID, userAccount)
	userPassWord := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    userIndexName,
		IndexName:    "用户信息",
		FieldCode:    "PassWord",
		FieldName:    "密码",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       3,
		Active:       1,
	}
	index.Index(userPassWord.ID, userPassWord)
	userName := IndexFieldStruct{
		ID:           FuncGenerateStringID(),
		IndexCode:    userIndexName,
		IndexName:    "用户信息",
		FieldCode:    "UserName",
		FieldName:    "用户名称",
		FieldType:    3,
		AnalyzerName: keyword.Name,
		CreateTime:   now,
		CreateUser:   createUser,
		SortNo:       4,
		Active:       1,
	}
	index.Index(userName.ID, userName)

	//创建用户表的索引
	mapping := bleve.NewIndexMapping()
	//指定默认的分词器
	mapping.DefaultMapping.DefaultAnalyzer = keyword.Name
	userIndex, err := bleve.New(userIndexName, mapping)
	if err != nil {
		FuncLogError(err)
		return false, err
	}
	bleveIndexMap[userIndexName] = userIndex

	//初始化数据
	user := make(map[string]string)
	id := FuncGenerateStringID()
	user["ID"] = id
	user["Account"] = "admin"
	user["PassWord"] = "21232f297a57a5a743894a0e4a801fc3"
	user["UserName"] = "管理员"
	//初始化 admin用户
	userIndex.Index(id, user)
	return true, nil
}

//初始化调用
func init() {
	checkInstall()
}

//FuncGenerateStringID 默认生成字符串ID的函数.方便自定义扩展
var FuncGenerateStringID func() string = generateStringID

//generateStringID 生成主键字符串
func generateStringID() string {
	pk, errUUID := gouuid.NewV4()
	if errUUID != nil {
		return ""
	}
	return pk.String()
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
