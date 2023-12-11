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
	"image/color"
	"sync"
	"sync/atomic"

	"github.com/mojocn/base64Captcha"
)

const errCount = 3

// 配置验证码的参数
var driverString = base64Captcha.DriverString{
	Height:          60,
	Width:           200,
	NoiseCount:      0,
	ShowLineOptions: 2 | 4,
	Length:          4,
	Source:          "34678acdefghkmnprtuvwxy",
	BgColor:         &color.RGBA{R: 3, G: 102, B: 214, A: 125},
	Fonts:           []string{"wqy-microhei.ttc"},
}

// var driver = base64Captcha.NewDriverMath(60, 200, 0, base64Captcha.OptionShowHollowLine, nil, nil, []string{"wqy-microhei.ttc"})
var driver = driverString.ConvertFonts()

var captchaQuestion, captchaAnswer, captchaBase64 string
var captchaLock = &sync.Mutex{}

var chainRandStr string

var errorLoginCount atomic.Uint32

func generateCaptcha() {
	captchaLock.Lock()
	defer captchaLock.Unlock()
	_, captchaQuestion, captchaAnswer = driver.GenerateIdQuestionAnswer()
	item, err := driver.DrawCaptcha(captchaQuestion)
	if err != nil {
		captchaBase64 = ""
		return
	}
	captchaBase64 = item.EncodeB64string()
}

func generateChainRandStr() {
	captchaLock.Lock()
	defer captchaLock.Unlock()
	//先记录到全局变量
	chainRandStr = randStr(30)
}
