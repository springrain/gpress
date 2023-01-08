package main

const (
	// 索引名称

	// 基本目录
	datadir = "gpressdatadir/"
	// 数据目录,如果不存在认为是第一次安装启动,会创建默认的数据
	bleveDataDir = datadir + "bleve/"
	// 索引字段的名称
	indexFieldIndexName = bleveDataDir + "indexField"
	// User 用户的索引名称
	userIndexName = bleveDataDir + "user"
	// siteInfo  站点信息
	indexSitenIndexName = bleveDataDir + "sitenInfo"
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

	//模板的路径
	templateDir = datadir + "template/"

	//静态化文件目录,网站生成的静态html
	statichtmlDir = datadir + "statichtml/"

	//默认名称
	defaultName = "gpress"

	// 数据默认的创建用户
	createUser = "system"

	tokenUserId = "userId"

	letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)
