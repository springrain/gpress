package main

import (
	"gitee.com/gpress/gpress/bleves"
	"gitee.com/gpress/gpress/route"
	"gitee.com/gpress/gpress/service"
)

func main() {

	//1.初始化bleve
	bleves.InitCommaAnalyzer()
	bleves.InitGesAnglyzer()
	bleves.BleveStatus = bleves.CheckBleveStatus()
	bleves.Installed = bleves.IsInstalled()

	//2.加载配置文件
	service.Config = service.LoadInstallConfig(bleves.Installed)
	// 使用的主题
	route.ThemePath = "/theme/" + service.Config.Theme + "/"

	if !bleves.BleveStatus { // 索引状态检查失败
		panic("索引检查失败")
	}
	route.InitRoute(service.Config.ServerPort)
	//初始化模版
	_ = route.InitTemplate()
	//启动路由
	route.RunServer()

}
