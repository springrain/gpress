package main

const (
	// 默认名称
	appName = "gpress"

	// 基本目录
	datadir = "gpressdatadir/"
	// 数据目录,如果不存在认为是第一次安装启动,会创建默认的数据
	sqliteDBfile = datadir + "gpress.db"
	// 表信息的名称
	//tableInfoName = "tableInfo"
	// 表字段的名称
	//tableFieldName = "tableField"

	// config 配置的表名称
	tableConfigName = "config"

	// user 用户的表名称
	tableUserName = "user"
	// site  站点信息
	tableSiteName = "site"
	// 页面模板
	tablePageTemplateName = "pageTemplate"
	// 导航菜单
	tableCategoryName = "category"

	// 默认模型
	tableModuleDefaultName = "module_default"
	// 文章内容
	tableContentName = "content"
	//---------------------------//

	// 模板的路径
	templateDir = datadir + "template/"

	// 静态化文件目录,网站生成的静态html
	//statichtmlDir = datadir + "statichtml/"

	// 数据默认的创建用户
	createUser = "system"

	tokenUserId = "userId"

	letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	whereConditionKey = "whereCondition"
)
