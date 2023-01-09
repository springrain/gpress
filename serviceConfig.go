package main

import (
	"context"
	"strconv"

	"github.com/blevesearch/bleve/v2"
)

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
