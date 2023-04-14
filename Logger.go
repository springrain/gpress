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
