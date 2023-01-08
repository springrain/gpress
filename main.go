package main

import (
	"strconv"

	"github.com/cloudwego/hertz/pkg/app/server"
)

// 是否已经安装过了
// var installed = isInstalled()
var installed = true //方便开发测试

// 加载配置文件
var config = loadInstallConfig()

// hertz对象,可以在其他地方使用
var h = server.Default(server.WithHostPorts(":" + strconv.Itoa(config.Port)))

func init() {
	//h.Use(gzip.Gzip(gzip.DefaultCompression))

}

func main() {
	//初始化模板
	initTemplate()

	//初始化admin路由
	initAdminRoute()

	// 启动服务
	h.Spin()
}
