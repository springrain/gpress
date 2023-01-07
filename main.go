package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

// hertz对象,可以在其他地方使用
var h = server.Default(server.WithHostPorts(":8080"))

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
