package service

import (
	"context"
	"encoding/json"
	"errors"
	"gitee.com/gpress/gpress/constant"
	"gitee.com/gpress/gpress/gbleve"
	"gitee.com/gpress/gpress/logger"
	"gitee.com/gpress/gpress/model"
	"gitee.com/gpress/gpress/util"
	"github.com/blevesearch/bleve/v2"
	"io"
	"os"
	"strconv"
	"time"
)

var Config model.ConfigStruct

// LoadInstallConfig 加载配置文件,只有初始化安装时需要读取配置文件,读取后,就写入索引,通过后台管理,然后重命名为 install_config.json_配置已失效_请通过后台设置管理
func LoadInstallConfig(installed bool) model.ConfigStruct {
	defaultErr := errors.New("install_config.json加载失败,使用默认配置")
	if installed { // 如果已经安装,需要从索引读取配置,这里暂时返回defaultConfig
		config, err := findConfig()
		if err != nil {
			return defaultConfig
		}
		return config
	}
	// 打开文件
	jsonFile, err := os.Open(constant.DATA_DIR + "install_config.json")
	if err != nil {
		logger.FuncLogError(defaultErr)
		return defaultConfig
	}
	// 关闭文件
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		logger.FuncLogError(defaultErr)
		return defaultConfig
	}
	configJson := model.ConfigStruct{}
	// Decode从输入流读取下一个json编码值并保存在v指向的值里
	err = json.Unmarshal([]byte(byteValue), &configJson)
	if err != nil {
		logger.FuncLogError(defaultErr)
		return defaultConfig
	}

	if configJson.JwtSecret == "" { // 如果没有配置jwtSecret,产生随机字符串
		configJson.JwtSecret = util.RandStr(32)
	}

	return configJson
}

func findConfig() (model.ConfigStruct, error) {
	configIndex := gbleve.IndexMap[constant.CONFIG_INDEX_NAME]
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
			config.JwtTokenKey = v.Fields["configValue"].(string)
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

// insertConfig 插入config
func InsertConfig(ctx context.Context, config model.ConfigStruct) error {
	// 清空配置,重新创建
	DeleteAll(ctx, constant.CONFIG_INDEX_NAME)

	configIndex := gbleve.IndexMap[constant.CONFIG_INDEX_NAME]

	// basePath
	basePath := make(map[string]string)
	basePathId := util.FuncGenerateStringID()
	basePath["id"] = basePathId
	basePath["configKey"] = "basePath"
	basePath["configValue"] = config.BasePath
	err := configIndex.Index(basePathId, basePath)
	if err != nil {
		return err
	}

	// jwtSecret
	jwtSecret := make(map[string]string)
	jwtSecretId := util.FuncGenerateStringID()
	jwtSecret["id"] = jwtSecretId
	jwtSecret["configKey"] = "jwtSecret"
	jwtSecret["configValue"] = config.JwtSecret
	err = configIndex.Index(jwtSecretId, jwtSecret)
	if err != nil {
		return err
	}

	// jwttokenKey
	jwttokenKey := make(map[string]string)
	jwttokenKeyId := util.FuncGenerateStringID()
	jwttokenKey["id"] = jwttokenKeyId
	jwttokenKey["configKey"] = "jwttokenKey"
	jwttokenKey["configValue"] = config.JwtTokenKey
	err = configIndex.Index(jwttokenKeyId, jwttokenKey)
	if err != nil {
		return err
	}

	// serverPort
	serverPort := make(map[string]string)
	serverPortId := util.FuncGenerateStringID()
	serverPort["id"] = serverPortId
	serverPort["configKey"] = "serverPort"
	serverPort["configValue"] = config.ServerPort
	err = configIndex.Index(serverPortId, serverPort)
	if err != nil {
		return err
	}

	// theme
	theme := make(map[string]string)
	themeId := util.FuncGenerateStringID()
	theme["id"] = themeId
	theme["configKey"] = "theme"
	theme["configValue"] = config.Theme
	err = configIndex.Index(themeId, theme)
	if err != nil {
		return err
	}

	// timeout
	timeout := make(map[string]interface{})
	timeoutId := util.FuncGenerateStringID()
	timeout["id"] = timeoutId
	timeout["configKey"] = "timeout"
	timeout["configValue"] = strconv.Itoa(config.Timeout)
	err = configIndex.Index(timeoutId, timeout)
	if err != nil {
		return err
	}

	return nil
}

// UpdateInstall 更新安装状态
func UpdateInstall(ctx context.Context) error {
	// 将config配置写入到索引,写入前先把config表清空
	err := InsertConfig(ctx, Config)
	if err != nil {
		return err
	}

	now := strconv.FormatInt(time.Now().UnixNano(), 10)
	// 删除 install 文件
	err = os.Rename(constant.TEMPLATE_DIR+"admin/install.html", constant.TEMPLATE_DIR+"admin/install.html."+now)
	if err != nil {
		return err
	}

	// install_config.json 重命名为 install_config.json_配置已失效_请通过后台设置管理
	err = os.Rename(constant.DATA_DIR+"install_config.json", constant.DATA_DIR+"install_config.json."+now)
	if err != nil {
		return err
	}
	// 更改安装状态
	gbleve.Installed = true
	return nil
}
