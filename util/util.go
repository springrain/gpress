package util

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"gitee.com/gpress/gpress/constant"
	"math/big"
	rand2 "math/rand"
	"os"
	"strings"
	"time"
)

/**
工具包
*/

// RandStr 生成随机字符串
func RandStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = constant.LETTERS[rand2.Intn(len(constant.LETTERS))]
	}
	return string(b)
}

// HashSha256 使用sha256计算hash值
func HashSha256(str string) string {
	hashByte := sha256.Sum256([]byte(str))
	hashStr := hex.EncodeToString(hashByte[:])
	return hashStr
}

// PathExists 文件或者目录是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	return !os.IsNotExist(err), err
}

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

// 测试自定义函数
func FuncMD5(in string) ([]string, error) {
	list := make([]string, 2)

	hash := md5.Sum([]byte(in))
	list[0] = in
	list[1] = hex.EncodeToString(hash[:])
	return list, nil
}

// FuncT 多语言i18n适配,例如 {{ T "nextPage" }}
func FuncT(key string) (string, error) {
	return key, nil
}
