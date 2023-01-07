package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

// 模板路径,正常应该从siteInfo里获取,这里用于演示
var themePath = "theme/default/"

// hertz对象,可以在其他地方使用
var h *server.Hertz

func init() {
	h = server.Default(server.WithHostPorts(":8080"))
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
