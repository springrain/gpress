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
	"log"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func init() {
	// 设置默认的日志显示信息,显示文件和行号
	// Set the default log display information, display file and line number.
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

// LogCallDepth 记录日志调用层级,用于定位到业务层代码
// Log Call Depth Record the log call level, used to locate the business layer code
var LogCallDepth = 4

// FuncLogError 记录error日志
// FuncLogError Record error log
var FuncLogError func(err error) = defaultLogError

// FuncLogPanic  记录panic日志,默认使用"ZormErrorLog"实现
// FuncLogPanic Record panic log, using "Zorm Error Log" by default
var FuncLogPanic func(err error) = defaultLogPanic

func defaultLogError(err error) {
	//log.Output(LogCallDepth, fmt.Sprintln(err))
	hlog.Error(err)
}

func defaultLogPanic(err error) {
	defaultLogError(err)
}
