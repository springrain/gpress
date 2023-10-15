package main

import "gitee.com/chunanyong/zorm"

// 设置
type Config struct {
	// 引入默认的struct,隔离IEntityStruct的方法改动
	zorm.EntityStruct

	// ID 主键 值为 TableName,也就是表名
	Id string `column:"id" json:"id,omitempty"`

	// Title
	BasePath string `column:"basePath" json:"basePath,omitempty"`

	// TableType index/module 表和模型,两种类型
	JwtSecret string `column:"jwtSecret" json:"jwtSecret,omitempty"`

	// TableType index/module 表和模型,两种类型
	JwttokenKey string `column:"jwttokenKey" json:"jwttokenKey,omitempty"`

	// TableType index/module 表和模型,两种类型
	ServerPort string `column:"serverPort" json:"serverPort,omitempty"`

	// TableType index/module 表和模型,两种类型
	Theme string `column:"theme" json:"theme,omitempty"`

	// TableType index/module 表和模型,两种类型
	Timeout int `column:"timeout" json:"timeout,omitempty"`

	// TableType index/module 表和模型,两种类型
	Proxy string `column:"proxy" json:"proxy,omitempty"`

	// TableType index/module 表和模型,两种类型
	MainServer string `column:"mainServer" json:"mainServer,omitempty"`

	// CreateTime 创建时间
	CreateTime string `column:"createTime" json:"createTime,omitempty"`

	// UpdateTime 更新时间
	UpdateTime string `column:"updateTime" json:"updateTime,omitempty"`

	// CreateUser  创建人,初始化 system
	CreateUser string `column:"createUser" json:"createUser,omitempty"`

	// SortNo 排序
	SortNo int `column:"sortNo" json:"sortNo,omitempty"`

	// 是否有效 无效(0),正常显示(1),界面不显示(3)
	Status int `column:"status" json:"status,omitempty"`
}

// GetTableName 获取表名称
// IEntityStruct 接口的方法,实体类需要实现!!!
func (entity *Config) GetTableName() string {
	return tableConfigName
}

// GetPKColumnName 获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称
// 不支持联合主键,变通认为无主键,业务控制实现(艰难取舍)
// 如果没有主键,也需要实现这个方法, return "" 即可
// IEntityStruct 接口的方法,实体类需要实现!!!
func (entity *Config) GetPKColumnName() string {
	return "id"
}

// 导航菜单
type Category struct {
	// 引入默认的struct,隔离IEntityStruct的方法改动
	zorm.EntityStruct

	// ID 主键 值为 TableName,也就是表名
	Id string `column:"id" json:"id,omitempty"`

	// Title
	Name string `column:"name" json:"name,omitempty"`

	// TableType index/module 表和模型,两种类型
	HrefURL string `column:"hrefURL" json:"hrefURL,omitempty"`

	// TableType index/module 表和模型,两种类型
	HrefTarget string `column:"hrefTarget" json:"hrefTarget,omitempty"`

	// TableType index/module 表和模型,两种类型
	Pid string `column:"pid" json:"pid,omitempty"`

	// TableType index/module 表和模型,两种类型
	ThemePC string `column:"themePC" json:"themePC,omitempty"`

	// TableType index/module 表和模型,两种类型
	ModuleID string `column:"moduleID" json:"moduleID,omitempty"`

	// TableType index/module 表和模型,两种类型
	ComCode string `column:"comCode" json:"comCode,omitempty"`

	// TableType index/module 表和模型,两种类型
	TemplateID string `column:"templateID" json:"templateID,omitempty"`

	// TableType index/module 表和模型,两种类型
	ChildTemplateID string `column:"childTemplateID" json:"childTemplateID,omitempty"`

	// CreateTime 创建时间
	CreateTime string `column:"createTime" json:"createTime,omitempty"`

	// UpdateTime 更新时间
	UpdateTime string `column:"updateTime" json:"updateTime,omitempty"`

	// CreateUser  创建人,初始化 system
	CreateUser string `column:"createUser" json:"createUser,omitempty"`

	// SortNo 排序
	SortNo int `column:"sortNo" json:"sortNo,omitempty"`

	// 是否有效 无效(0),正常显示(1),界面不显示(3)
	Status int `column:"status" json:"status,omitempty"`
}

// GetTableName 获取表名称
// IEntityStruct 接口的方法,实体类需要实现!!!
func (entity *Category) GetTableName() string {
	return tableCategoryName
}

// GetPKColumnName 获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称
// 不支持联合主键,变通认为无主键,业务控制实现(艰难取舍)
// 如果没有主键,也需要实现这个方法, return "" 即可
// IEntityStruct 接口的方法,实体类需要实现!!!
func (entity *Category) GetPKColumnName() string {
	return "id"
}

// 内容文章
type Content struct {

	// 引入默认的struct,隔离IEntityStruct的方法改动
	zorm.EntityStruct

	// ID 主键 值为 TableName,也就是表名
	Id string `column:"id" json:"id,omitempty"`

	// Title 文章标题
	Title string `column:"title" json:"title,omitempty"`

	// ModuleID 模型ID
	ModuleID string `column:"moduleID" json:"moduleID,omitempty"`

	// Keyword 关键字
	Keyword string `column:"keyword" json:"keyword,omitempty"`

	// Description 站点描述
	Description string `column:"description" json:"description,omitempty"`

	// HrefURL 自身页面路径
	HrefURL string `column:"hrefURL" json:"hrefURL,omitempty"`

	// Subtitle 副标题
	Subtitle string `column:"subtitle" json:"subtitle,omitempty"`

	// CategoryID 导航ID
	CategoryID string `column:"categoryID" json:"categoryID,omitempty"`

	// CategoryName 导航名称
	CategoryName string `column:"categoryName" json:"categoryName,omitempty"`

	// TableType index/module 表和模型,两种类型
	TemplateID string `column:"templateID" json:"templateID,omitempty"`

	// TableType index/module 表和模型,两种类型
	Author string `column:"author" json:"author,omitempty"`

	// TableType index/module 表和模型,两种类型
	Tag string `column:"tag" json:"tag,omitempty"`

	// TableType index/module 表和模型,两种类型
	Toc string `column:"toc" json:"toc,omitempty"`

	// TableType index/module 表和模型,两种类型
	Summary string `column:"summary" json:"summary,omitempty"`

	// TableType index/module 表和模型,两种类型
	Content string `column:"content" json:"content,omitempty"`

	// TableType index/module 表和模型,两种类型
	Markdown string `column:"markdown" json:"markdown,omitempty"`

	// TableType index/module 表和模型,两种类型
	Thumbnail string `column:"thumbnail" json:"thumbnail,omitempty"`

	// CreateTime 创建时间
	CreateTime string `column:"createTime" json:"createTime,omitempty"`

	// UpdateTime 更新时间
	UpdateTime string `column:"updateTime" json:"updateTime,omitempty"`

	// CreateUser  创建人,初始化 system
	CreateUser string `column:"createUser" json:"createUser,omitempty"`

	// SortNo 排序
	SortNo int `column:"sortNo" json:"sortNo,omitempty"`

	// 是否有效 无效(0),正常显示(1),界面不显示(3)
	Status int `column:"status" json:"status,omitempty"`
}

// GetTableName 获取表名称
// IEntityStruct 接口的方法,实体类需要实现!!!
func (entity *Content) GetTableName() string {
	return tableContentName
}

// GetPKColumnName 获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称
// 不支持联合主键,变通认为无主键,业务控制实现(艰难取舍)
// 如果没有主键,也需要实现这个方法, return "" 即可
// IEntityStruct 接口的方法,实体类需要实现!!!
func (entity *Content) GetPKColumnName() string {
	return "id"
}

// 模板
type PageTemplate struct {
	// 引入默认的struct,隔离IEntityStruct的方法改动
	zorm.EntityStruct

	// 主键
	Id string `column:"id" json:"id,omitempty"`

	// 名称
	Name string `column:"name" json:"name,omitempty"`

	// TemplatePath 模板路径
	TemplatePath string `column:"templatePath" json:"templatePath,omitempty"`

	// 创建时间
	CreateTime string `column:"createTime" json:"createTime,omitempty"`

	// 创建人
	CreateUser string `column:"createUser" json:"createUser,omitempty"`

	// 修改时间
	UpdateTime string `column:"updateTime" json:"updateTime,omitempty"`

	// 排序
	SortNo int `column:"sortNo" json:"sortNo,omitempty"`

	// 状态
	Status int `column:"status" json:"status,omitempty"`
}

// GetTableName 获取表名称
// IEntityStruct 接口的方法,实体类需要实现!!!
func (entity *PageTemplate) GetTableName() string {
	return tablePageTemplateName
}

func (entity *PageTemplate) GetPKColumnName() string {
	return "id"
}

// 站点信息
type Site struct {
	// 引入默认的struct,隔离IEntityStruct的方法改动
	zorm.EntityStruct

	// ID 主键 值为 TableName,也就是表名
	Id string `column:"id" json:"id"`

	// Title 头
	Title string `column:"title" json:"title,omitempty"`

	// Name 名字
	Name string `column:"name" json:"name,omitempty"`

	// Domain
	Domain string `column:"domain" json:"domain,omitempty"`

	// Keyword
	Keyword string `column:"keyword" json:"keyword,omitempty"`

	// Description
	Description string `column:"description" json:"description,omitempty"`

	// Theme
	Theme string `column:"theme" json:"theme,omitempty"`

	// ThemePC
	ThemePC string `column:"themePC" json:"themePC,omitempty"`

	// ThemeWAP
	ThemeWAP string `column:"themeWAP" json:"themeWAP,omitempty"`

	// SiteThemeWEIXIN
	SiteThemeWEIXIN string `column:"siteThemeWEIXIN" json:"siteThemeWEIXIN,omitempty"`

	// Logo
	Logo string `column:"logo" json:"logo,omitempty"`

	// Favicon
	Favicon string `column:"favicon" json:"favicon,omitempty"`

	// Footer
	Footer string `column:"footer" json:"footer,omitempty"`

	// CreateTime 创建时间
	CreateTime string `column:"createTime" json:"createTime,omitempty"`

	// UpdateTime 更新时间
	UpdateTime string `column:"updateTime" json:"updateTime,omitempty"`

	// CreateUser  创建人,初始化 system
	CreateUser string `column:"createUser" json:"createUser,omitempty"`

	// SortNo 排序
	SortNo int `column:"sortNo" json:"sortNo,omitempty"`

	// 是否有效 无效(0),正常显示(1),界面不显示(3)
	Status int `column:"status" json:"status,omitempty"`
}

// GetTableName 获取表名称
// IEntityStruct 接口的方法,实体类需要实现!!!
func (entity *Site) GetTableName() string {
	return tableSiteName
}

// GetPKColumnName 获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称
// 不支持联合主键,变通认为无主键,业务控制实现(艰难取舍)
// 如果没有主键,也需要实现这个方法, return "" 即可
// IEntityStruct 接口的方法,实体类需要实现!!!
func (entity *Site) GetPKColumnName() string {
	return "id"
}

// 用户信息
type User struct {
	// 引入默认的struct,隔离IEntityStruct的方法改动
	zorm.EntityStruct

	// Id 主键
	Id string `column:"id" json:"id,omitempty"`

	// 账号
	Account string `column:"account" json:"account,omitempty"`

	// 密码
	Password string `column:"password" json:"password,omitempty"`

	// 用户名
	UserName string `column:"userName" json:"userName,omitempty"`

	// 链类型
	ChainType string `column:"chainType" json:"chainType,omitempty"`

	// 链address
	ChainAddress string `column:"chainAddress" json:"chainAddress,omitempty"`

	// rsa公钥
	RsaPublicKey string `column:"rsaPublicKey" json:"rsaPublicKey,omitempty"`

	// CreateTime
	CreateTime string `column:"createTime" json:"createTime,omitempty"`

	// UpdateTime
	UpdateTime string `column:"updateTime" json:"updateTime,omitempty"`

	// 创建人
	CreateUser string `column:"createUser" json:"createUser,omitempty"`

	//排序
	SortNo int `column:"sortNo" json:"sortNo,omitempty"`

	//状态
	Status int `column:"status" json:"status,omitempty"`
}

// GetTableName 获取表名称
// IEntityStruct 接口的方法,实体类需要实现!!!
func (entity *User) GetTableName() string {
	return tableUserName
}

func (entity *User) GetPKColumnName() string {
	return "id"
}
