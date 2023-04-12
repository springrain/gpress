package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"

	"gitee.com/chunanyong/zorm"
	"github.com/blevesearch/bleve/v2"
)

// 加载配置文件,只有初始化安装时需要读取配置文件,读取后,就写入索引,通过后台管理,然后重命名为 install_config.json_配置已失效_请通过后台设置管理
func loadInstallConfig() configStruct {
	defaultErr := errors.New("install_config.json加载失败,使用默认配置")
	if installed { // 如果已经安装,需要从索引读取配置,这里暂时返回defaultConfig
		config, err := findConfig()
		if err != nil {
			return defaultConfig
		}
		return config
	}
	// 打开文件
	jsonFile, err := os.Open(datadir + "install_config.json")
	if err != nil {
		FuncLogError(defaultErr)
		return defaultConfig
	}
	// 关闭文件
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		FuncLogError(defaultErr)
		return defaultConfig
	}
	configJson := configStruct{}
	// Decode从输入流读取下一个json编码值并保存在v指向的值里
	err = json.Unmarshal([]byte(byteValue), &configJson)
	if err != nil {
		FuncLogError(defaultErr)
		return defaultConfig
	}

	if configJson.JwtSecret == "" { // 如果没有配置jwtSecret,产生随机字符串
		configJson.JwtSecret = randStr(32)
	}
	if configJson.BasePath == "" {
		configJson.BasePath = "/"
	}

	return configJson
}

var defaultConfig = configStruct{
	BasePath: "/",
	// 默认的加密Secret
	// JwtSecret:   "gpress+jwtSecret-2023",
	JwtSecret:   randStr(32),
	Theme:       "default",
	JwttokenKey: "jwttoken", // jwt的key
	Timeout:     1800,       // 半个小时超时
	ServerPort:  ":660",     // gpress: 103 + 112 + 114 + 101 + 115 + 115 = 660
}

type configStruct struct {
	BasePath    string `json:"basePath"`
	JwtSecret   string `json:"jwtSecret"`
	JwttokenKey string `json:"jwttokenKey"`
	Timeout     int    `json:"timeout"`
	ServerPort  string `json:"serverPort"`
	Theme       string `json:"theme"`
}

// insertConfig 插入config
func insertConfig(ctx context.Context, config configStruct) error {
	// 清空配置,重新创建
	deleteAll(ctx, indexConfigName)

	ID := FuncGenerateStringID()

	m := make(map[string]interface{})
	m["id"] = ID
	b, _ := json.Marshal(config)
	json.Unmarshal(b, &m)

	entityMap := zorm.NewEntityMap(indexConfigName)
	for k, v := range m {
		entityMap.Set(k, v)
	}
	saveEntityMap(entityMap)
	return nil
}

func findConfig() (configStruct, error) {
	query := bleve.NewQueryStringQuery("*")
	searchRequest := bleve.NewSearchRequestOptions(query, 100, 0, false)
	searchRequest.Fields = []string{"*"}
	config := defaultConfig
	result, err := bleveSearchInContext(context.Background(), indexConfigName, searchRequest)
	if err != nil {
		return config, err
	}
	m, err := result2Map(indexConfigName, result)
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

	return config, nil
}
