package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
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

// pathExists 文件或者目录是否存在
func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// result2Map 单个查询结果转map
func result2Map(indexName string, result *bleve.SearchResult) (map[string]interface{}, error) {
	if result == nil {
		return nil, errors.New("结果集为空")
	}
	if result.Total == 0 { //没有记录
		return nil, nil
	}
	if result.Total > 1 { // 大于1条记录
		return nil, errors.New("查询出多条记录")
	}
	//获取到查询的对象
	value := result.Hits[0]
	m := make(map[string]interface{}, 0)
	for k, v := range value.Fields {
		m[k] = v
	}
	return m, nil

}

// result2SliceMap 多条结果转map数组
func result2SliceMap(indexName string, result *bleve.SearchResult) ([]map[string]interface{}, error) {
	if result == nil {
		return nil, errors.New("结果集为空")
	}
	if result.Total == 0 { //没有记录
		return nil, nil
	}
	ms := make([]map[string]interface{}, 0)
	//获取到查询的对象
	for _, value := range result.Hits {
		m := make(map[string]interface{}, 0)
		for k, v := range value.Fields {
			m[k] = v
			ms = append(ms, m)
		}
	}
	return ms, nil
}
