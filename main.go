package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

// 变量的位置不要更改!!!!!,实际是做初始化使用的,优先级高于init函数!!!

// 检查表状态
var sqliteStatus = checkSQLiteStatus()

// 是否已经安装过了
var installed = isInstalled()

// 加载配置文件
var config = loadInstallConfig()

// 使用的主题
var themePath = "/theme/" + config.Theme + "/"

// hertz对象,可以在其他地方使用
var h = server.Default(server.WithHostPorts(config.ServerPort), server.WithBasePath(config.BasePath))

func init() {
	if !sqliteStatus { // 表状态检查失败
		panic("表检查失败")
	}

	// 设置随机种子
	//rand.Seed(time.Now().UnixNano())

	// h.Use(gzip.Gzip(gzip.DefaultCompression))
}

func main() {
	// 初始化admin路由,使用init实现
	//initAdminRoute()

	//注册category和context的href路由
	err := registerHrefRoute()
	if err != nil {
		FuncLogError(err)
	}

	// 启动服务
	h.Spin()
}
