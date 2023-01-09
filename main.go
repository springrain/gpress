package main

import (
	"math/rand"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
)

// 变量的位置不要更改!!!!!,实际是做初始化使用的,优先级高于init

// 是否已经安装过了
var installed = isInstalled()

// 检查索引状态
var bleveStatus = checkBleveStatus()

// 加载配置文件
var config = loadInstallConfig()

// 使用的主题
var themePath = "/theme/" + config.Theme + "/"

// hertz对象,可以在其他地方使用
var h = server.Default(server.WithHostPorts(config.ServerPort))

func init() {
	if !bleveStatus { //索引状态检查失败
		panic("索引检查失败")
	}
	//设置随机种子
	rand.Seed(time.Now().UnixNano())

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
