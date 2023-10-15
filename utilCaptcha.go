package main

import (
	"image/color"
	"sync"
	"sync/atomic"

	"github.com/mojocn/base64Captcha"
)

const errCount = 2

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
