package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"
)

// FuncGenerateStringID 默认生成字符串ID的函数.方便自定义扩展
// FuncGenerateStringID Function to generate string ID by default. Convenient for custom extension
var FuncGenerateStringID func() string = generateStringID

// generateStringID 生成主键字符串
// generateStringID Generate primary key string
func generateStringID() string {
	// 使用 crypto/rand 真随机9位数
	randNum, randErr := rand.Int(rand.Reader, big.NewInt(1000000000))
	if randErr != nil {
		return ""
	}
	// 获取9位数,前置补0,确保9位数
	rand9 := fmt.Sprintf("%09d", randNum)

	// 获取纳秒 按照 年月日时分秒毫秒微秒纳秒 拼接为长度23位的字符串
	pk := time.Now().Format("2006.01.02.15.04.05.000000000")
	pk = strings.ReplaceAll(pk, ".", "")

	// 23位字符串+9位随机数=32位字符串,这样的好处就是可以使用ID进行排序
	pk = pk + rand9
	return pk
}

// PathExists 文件或者目录是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
