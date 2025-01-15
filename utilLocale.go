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
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
)

// localeMap 用于记录翻译的Map
var localeMap = make(map[string]string)

// initLocale 要在config初始化之后,需要获取config中的语言配置
func initLocale() {
	defaultErr := errors.New(config.Locale + ".json加载失败,使用默认zh-CN.json")
	// 打开文件
	jsonFile, err := os.Open(datadir + "/locales/" + config.Locale + ".json")
	if err != nil {
		FuncLogError(nil, defaultErr)
		jsonFile, err = os.Open(datadir + "/locales/zh-CN.json")
		if err != nil {
			FuncLogError(nil, defaultErr)
			return
		}
	}
	// 关闭文件
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		FuncLogError(nil, defaultErr)
		return
	}
	localeMapTemp := make(map[string]string)
	// Decode从输入流读取下一个json编码值并保存在v指向的值里
	err = json.Unmarshal([]byte(byteValue), &localeMapTemp)
	if err != nil {
		FuncLogError(nil, defaultErr)
		return
	}
	//不区分大小写
	for key, value := range localeMapTemp {
		localeMap[strings.ToLower(key)] = value
	}
	localeMapTemp = nil
}
