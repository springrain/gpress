package main

import "github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"

const (
	// 基本目录
	datadir = "gpressdatadir/"
	// 数据目录,如果不存在认为是第一次安装启动,会创建默认的数据
	bleveDataDir = datadir + "bleve/"
	// 索引字段的名称
	indexFieldIndexName = bleveDataDir + "indexField"

	// config 配置的索引名称
	configIndexName = bleveDataDir + "config"

	// user 用户的索引名称
	userIndexName = bleveDataDir + "user"
	// site  站点信息
	indexSiteIndexName = bleveDataDir + "site"
	// 页面模板
	indexPageTemplateName = bleveDataDir + "pageTemplate"
	// 导航菜单
	indexNavMenuName = bleveDataDir + "navMenu"
	// 模型
	indexModuleName = bleveDataDir + "module"
	// 模型数据
	indexModuleDefaultName = "moduleDefault"
	// 文章内容
	indexContentName = bleveDataDir + "content"
	//---------------------------//

	// 模板的路径
	templateDir = datadir + "template/"

	// 静态化文件目录,网站生成的静态html
	statichtmlDir = datadir + "statichtml/"

	// 默认名称
	defaultName = "gpress"

	// 数据默认的创建用户
	createUser = "system"

	tokenUserId = "userId"

	letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// keyword 分词器名称,避免引入错误的包
	keywordAnalyzerName = keyword.Name

	// 逗号分词器名称
	commaAnalyzerName = "comma"

	// gse分词器名称
	gseAnalyzerName = "gse"
)

const (
	fieldType_数字 = iota + 1
	fieldType_日期
	fieldType_文本框
	fieldType_文本域
	fieldType_富文本
	fieldType_下拉框
	fieldType_单选
	fieldType_多选
	fieldType_上传图片
	fieldType_上传附件
	fieldType_轮播图
	fieldType_音频
	fieldType_视频
)
