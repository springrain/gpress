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
	//"github.com/hertz-contrib/gzip"
)

// 变量的位置不要更改!!!!!,实际是做初始化使用的,优先级高于init函数!!!

// 是否已经安装过了
var installed = isInstalled()

// 加载配置文件
var config, site = loadInstallConfig()

// 服务器url路径
var httpServerPath = "http://"

// hertz对象,可以在其他地方使用
var h = server.Default(server.WithHostPorts(config.ServerPort), server.WithBasePath(config.BasePath), server.WithMaxRequestBodySize(config.MaxRequestBodySize))

func init() {

	// 设置随机种子
	//rand.Seed(time.Now().UnixNano())

	// gzip压缩文件,产生  xxx.html.gz 文件,https://www.cloudwego.io/zh/docs/hertz/tutorials/basic-feature/middleware/gzip/
	// h.Use(gzip.Gzip(gzip.DefaultCompression))
}

func main() {

	// 初始化语言包
	initLocale()

	//加载页面模板
	err := loadTemplate()
	if err != nil { // 初始化模板异常
		panic(funcT("Template initialization anomaly"))
	}

	message := funcT("Open the front-end in the browse") + ": "
	if strings.HasPrefix(config.ServerPort, ":") {
		httpServerPath += "127.0.0.1"
	}
	httpServerPath += config.ServerPort + config.BasePath
	message += httpServerPath
	message += "\n" + funcT("Open the back-end in the browser") + ": " + httpServerPath + "admin/login"
	fmt.Println(message)

	// 启动服务
	h.Spin()
}
