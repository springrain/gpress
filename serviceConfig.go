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
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"

	"gitee.com/chunanyong/zorm"
)

// loadInstallConfig 加载配置文件,只有初始化安装时需要读取配置文件,读取后,就写入表,通过后台管理,然后重命名为 install_config.json_配置已失效_请通过后台设置管理
func loadInstallConfig() (Config, Site) {
	var site = Site{Theme: "default"}
	defaultErr := errors.New(funcT("Failed to load install_config.json, using default configuration"))
	if installed { // 已经安装,需要表读取配置
		var err error
		site, err = funcSite()
		if err != nil {
			return defaultConfig, site
		}
		config, err := findConfig()
		if err != nil {
			return defaultConfig, site
		}
		return config, site
	}
	// 打开文件
	jsonFile, err := os.Open(datadir + "install_config.json")
	if err != nil {
		FuncLogError(nil, defaultErr)
		return defaultConfig, site
	}
	// 关闭文件
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		FuncLogError(nil, defaultErr)
		return defaultConfig, site
	}
	configJson := Config{}
	// Decode从输入流读取下一个json编码值并保存在v指向的值里
	err = json.Unmarshal([]byte(byteValue), &configJson)
	if err != nil {
		FuncLogError(nil, defaultErr)
		return defaultConfig, site
	}

	if configJson.JwtSecret == "" { // 如果没有配置jwtSecret,产生随机字符串
		configJson.JwtSecret = randStr(32)
	}
	if configJson.BasePath == "" {
		configJson.BasePath = "/"
	}
	if configJson.Locale == "" {
		configJson.Locale = "zh-CN"
	}
	configJson.Id = defaultConfig.Id

	return configJson, site
}

var defaultConfig = Config{
	Id:       "gpress_config",
	BasePath: "/",
	// 默认的加密Secret
	// JwtSecret:   "gpress+jwtSecret-2023",
	JwtSecret: randStr(32),
	//Theme:       "default",
	MaxRequestBodySize: 20 * 1024 * 1024,
	JwttokenKey:        "jwttoken", // jwt的key
	Timeout:            7200,       // 两小时超时
	ServerPort:         ":660",     // gpress: 103 + 112 + 114 + 101 + 115 + 115 = 660
	Locale:             "zh-CN",
}

// insertConfig 插入config
func insertConfig(ctx context.Context) error {
	//数据库存在config,不更新数据库,更新config变量
	finder := zorm.NewSelectFinder(tableConfigName).Append("WHERE id=?", "gpress_config")
	c := Config{}
	has, err := zorm.QueryRow(ctx, finder, &c)
	if has && err == nil && c.Id != "" {
		config = c
		return err
	}

	// 清空配置,重新创建
	deleteAll(ctx, tableConfigName)
	_, err = zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		return zorm.Insert(ctx, &config)
	})

	return err
}

// findConfig 查询配置
func findConfig() (Config, error) {

	finder := zorm.NewSelectFinder(tableConfigName, "*")

	m, err := zorm.QueryRowMap(context.Background(), finder)

	config := defaultConfig
	if err != nil {
		return config, err
	}
	b, err := json.Marshal(m)
	if err != nil {
		return config, err
	}
	json.Unmarshal(b, &config)

	if config.BasePath == "" {
		config.BasePath = "/"
	}
	if config.MaxRequestBodySize == 0 {
		config.MaxRequestBodySize = 20 * 1024 * 1024
	}
	if config.Locale == "" {
		config.Locale = "zh-CN"
	}
	return config, nil
}
