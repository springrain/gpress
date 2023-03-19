package service

import (
	"gitee.com/gpress/gpress/model"
	"gitee.com/gpress/gpress/util"
)

// Inclusive 是否包含
var Inclusive = true

// Active 获取菜单树 pid 为0 为一级菜单
var Active float64 = 1

var defaultConfig = model.ConfigStruct{
	BasePath: "",
	// 默认的加密Secret
	// JwtSecret:   "gpress+jwtSecret-2023",
	JwtSecret:   util.RandStr(32),
	Theme:       "default",
	JwtTokenKey: "jwttoken", // jwt的key
	Timeout:     1800,       // 半个小时超时
	ServerPort:  ":660",     // gpress: 103 + 112 + 114 + 101 + 115 + 115 = 660
}
