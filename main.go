// Copyright (c) 2023 gpress Authors.
//
// This file is part of gpress.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"strings"

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
	/*
		err := registerHrefRoute()
		if err != nil {
			FuncLogError(err)
		}
	*/
	//生成json文件
	if !pathExist(searchDataJsonFile) {
		genSearchDataJson()
	}

	message := "浏览器打开前端: "
	httpStr := "http://"
	if strings.HasPrefix(config.ServerPort, ":") {
		httpStr += "127.0.0.1"
	}
	message += httpStr + config.ServerPort + config.BasePath
	message += "\n浏览器打开后台: " + httpStr + config.ServerPort + config.BasePath + "admin/login"
	fmt.Println(message)
	// 启动服务
	h.Spin()

	// 启用热更新

}
