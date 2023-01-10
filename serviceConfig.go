package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"

	"github.com/blevesearch/bleve/v2"
)

// 加载配置文件,只有初始化安装时需要读取配置文件,读取后,就写入索引,通过后台管理,然后重命名为 install_config.json_配置已失效_请通过后台设置管理
func loadInstallConfig() configStruct {
	defaultErr := errors.New("install_config.json加载失败,使用默认配置")
	if installed { //如果已经安装,需要从索引读取配置,这里暂时返回defaultConfig
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
	//Decode从输入流读取下一个json编码值并保存在v指向的值里
	err = json.Unmarshal([]byte(byteValue), &configJson)
	if err != nil {
		FuncLogError(defaultErr)
		return defaultConfig
	}

	if configJson.JwtSecret == "" { //如果没有配置jwtSecret,产生随机字符串
		configJson.JwtSecret = randStr(32)
	}

	return configJson
}

var defaultConfig = configStruct{
	BasePath: "",
	//默认的加密Secret
	//JwtSecret:   "gpress+jwtSecret-2023",
	JwtSecret:   randStr(32),
	Theme:       "default",
	JwttokenKey: "jwttoken", //jwt的key
	Timeout:     1800,       //半个小时超时
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
	//清空配置,重新创建
	deleteAll(ctx, configIndexName)

	configIndex := IndexMap[configIndexName]

	//basePath
	basePath := make(map[string]string)
	basePathId := FuncGenerateStringID()
	basePath["id"] = basePathId
	basePath["configKey"] = "basePath"
	basePath["configValue"] = config.BasePath
	err := configIndex.Index(basePathId, basePath)
	if err != nil {
		return err
	}

	//jwtSecret
	jwtSecret := make(map[string]string)
	jwtSecretId := FuncGenerateStringID()
	jwtSecret["id"] = jwtSecretId
	jwtSecret["configKey"] = "jwtSecret"
	jwtSecret["configValue"] = config.JwtSecret
	err = configIndex.Index(jwtSecretId, jwtSecret)
	if err != nil {
		return err
	}

	//jwttokenKey
	jwttokenKey := make(map[string]string)
	jwttokenKeyId := FuncGenerateStringID()
	jwttokenKey["id"] = jwttokenKeyId
	jwttokenKey["configKey"] = "jwttokenKey"
	jwttokenKey["configValue"] = config.JwttokenKey
	err = configIndex.Index(jwttokenKeyId, jwttokenKey)
	if err != nil {
		return err
	}

	//serverPort
	serverPort := make(map[string]string)
	serverPortId := FuncGenerateStringID()
	serverPort["id"] = serverPortId
	serverPort["configKey"] = "serverPort"
	serverPort["configValue"] = config.ServerPort
	err = configIndex.Index(serverPortId, serverPort)
	if err != nil {
		return err
	}

	//theme
	theme := make(map[string]string)
	themeId := FuncGenerateStringID()
	theme["id"] = themeId
	theme["configKey"] = "theme"
	theme["configValue"] = config.Theme
	err = configIndex.Index(themeId, theme)
	if err != nil {
		return err
	}

	//timeout
	timeout := make(map[string]interface{})
	timeoutId := FuncGenerateStringID()
	timeout["id"] = timeoutId
	timeout["configKey"] = "timeout"
	timeout["configValue"] = strconv.Itoa(config.Timeout)
	err = configIndex.Index(timeoutId, timeout)
	if err != nil {
		return err
	}

	return nil
}

func findConfig() (configStruct, error) {
	configIndex := IndexMap[configIndexName]
	query := bleve.NewQueryStringQuery("*")
	serarchRequest := bleve.NewSearchRequestOptions(query, 100, 0, false)
	serarchRequest.Fields = []string{"*"}
	config := defaultConfig
	result, err := configIndex.SearchInContext(context.Background(), serarchRequest)
	if err != nil {
		return config, err
	}
	for _, v := range result.Hits {
		configKey := v.Fields["configKey"].(string)
		switch configKey {
		case "basePath":
			config.BasePath = v.Fields["configValue"].(string)
		case "jwtSecret":
			config.JwtSecret = v.Fields["configValue"].(string)
		case "jwttokenKey":
			config.JwttokenKey = v.Fields["configValue"].(string)
		case "serverPort":
			config.ServerPort = v.Fields["configValue"].(string)
		case "theme":
			config.Theme = v.Fields["configValue"].(string)
		case "timeout":
			strv := v.Fields["configValue"].(string)
			config.Timeout, err = strconv.Atoi(strv)
			if err != nil {
				config.Timeout = defaultConfig.Timeout
				return config, err
			}
		}
	}
	return config, nil

}
