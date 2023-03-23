package main

import (
	"gitee.com/gpress/gpress/gbleve"
	"gitee.com/gpress/gpress/route"
	"gitee.com/gpress/gpress/service"
)

func main() {

	//1.初始化bleve
	gbleve.InitCommaAnalyzer()
	gbleve.InitGesAnglyzer()
	gbleve.BleveStatus = gbleve.CheckBleveStatus()
	gbleve.Installed = gbleve.IsInstalled()

	//2.加载配置文件
	service.Config = service.LoadInstallConfig(gbleve.Installed)
	// 使用的主题
	route.ThemePath = "/theme/" + service.Config.Theme + "/"

	if !gbleve.BleveStatus { // 索引状态检查失败
		panic("索引检查失败")
	}
	route.InitRoute(service.Config.ServerPort)
	//初始化模版
	_ = route.InitTemplate()
	//启动路由
	route.RunServer()

}
