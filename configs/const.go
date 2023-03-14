package configs

import "github.com/blevesearch/bleve/v2"

// IndexMap 全局存放 索引对象,启动之后,所有的索引都通过这个map获取,一个索引只能打开一次,类似数据库连接,用一个对象操作
var IndexMap map[string]bleve.Index = make(map[string]bleve.Index)

// 索引名称
const (
	//NORM_TIME_FORMAT 标准时间格式
	NORM_TIME_FORMAT = "2006-01-02 15:04:05"

	// DATA_DIR 基本目录
	DATA_DIR = "gpressdatadir/"

	// BLEVE_DATA_DIR 数据目录,如果不存在认为是第一次安装启动,会创建默认的数据
	BLEVE_DATA_DIR = DATA_DIR + "bleve/"

	// INDEX_FIELD_INDEX_NAME 索引字段的名称
	INDEX_FIELD_INDEX_NAME = BLEVE_DATA_DIR + "indexField"

	// CONFIG_INDEX_NAME  配置的索引名称
	CONFIG_INDEX_NAME = BLEVE_DATA_DIR + "config"

	//USER_INDEX_NAME 用户的索引名称
	USER_INDEX_NAME = BLEVE_DATA_DIR + "user"

	// INDEX_SITE_INDEX_NAME  站点信息
	INDEX_SITE_INDEX_NAME = BLEVE_DATA_DIR + "siteInfo"

	//INDEX_PAGE_TEMPLATE_NAME 页面模板
	INDEX_PAGE_TEMPLATE_NAME = BLEVE_DATA_DIR + "pageTemplate"

	//INDEX_NAV_MENU_NAME 导航菜单
	INDEX_NAV_MENU_NAME = BLEVE_DATA_DIR + "navMenu"

	//INDEX_MODULE_NAME 模型
	INDEX_MODULE_NAME = BLEVE_DATA_DIR + "module"

	//INDEX_MODULE_DEFAULT_NAME 模型数据
	INDEX_MODULE_DEFAULT_NAME = "moduleDefault"

	//INDEX_CONTENT_NAME 文章内容
	INDEX_CONTENT_NAME = BLEVE_DATA_DIR + "content"

	//---------------------------//

	//TEMPLATE_DIR 模板的路径
	TEMPLATE_DIR = DATA_DIR + "template/"

	//STATIC_HTML_DIR 静态化文件目录,网站生成的静态html
	STATIC_HTML_DIR = DATA_DIR + "statichtml/"

	//DEFAULT_NAME 默认名称
	DEFAULT_NAME = "gpress"

	//CREATE_USER 数据默认的创建用户
	CREATE_USER = "system"

	TOKEN_USER_ID = "userId"

	LETTERS = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	GSE_ANGLYZER_NAME = "gse"

	COMMA_ANALYZER_NAME = "comma"
)
