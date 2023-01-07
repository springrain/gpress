package main

const (
	// 索引名称

	// 基本目录
	datadir = "gpressdatadir/"
	// 数据目录,如果不存在认为是第一次安装启动,会创建默认的数据
	bleveDataDir = datadir + "bleve/"
	// 索引字段的名称
	indexFieldIndexName = bleveDataDir + "IndexField"
	// User 用户的索引名称
	userIndexName = bleveDataDir + "User"
	// siteInfo  站点信息
	indexSitenIndexName = bleveDataDir + "sitenInfo"
	// 页面模板
	indexPageTemplateName = bleveDataDir + "pageTemplate"
	// 导航菜单
	indexNavMenuName = bleveDataDir + "NavMenu"
	// 模型
	indexModuleName = bleveDataDir + "Module"
	// 模型数据
	indexModuleDefaultName = "moduleDefault"
	// 文章内容
	indexContentName = bleveDataDir + "Content"
	//---------------------------//

	//模板的路径
	templateDir = datadir + "template/"

	//默认名称
	defaultName = "gpress"

	// 数据默认的创建用户
	createUser = "system"
)

// 变量
var (
	//
	basePath = ""
	//默认的加密Secret
	jwtSecret   = "gpressjwtSecret"
	timeout     = 1800       //半个小时超时
	jwttokenKey = "jwttoken" //jwt的key

)
