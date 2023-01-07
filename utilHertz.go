package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

var funcMap = template.FuncMap{"md5": funcMD5, "basePath": funcBasePath}

// initTemplate 初始化模板
func initTemplate() error {
	h.SetFuncMap(funcMap)
	//h.LoadHTMLFiles(themePath + "index.html")
	//h.LoadHTMLGlob(datadir + "html/theme/default/*")
	// 手动声明template对象,自己控制文件路径,默认是使用文件名,多个文件夹会存在问题
	err := loadTemplate()
	//设置默认的静态文件,实际路径会拼接为 datadir/public
	h.Static("/public", datadir)
	return err
}

// loadTemplate 用于更新重复加载
func loadTemplate() error {
	tmpl := template.New("").Delims("", "").Funcs(funcMap)
	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		// 分隔符统一为 / 斜杠
		path = filepath.ToSlash(path)
		//相对路径
		relativePath := path[len(templateDir)-1:]
		// 如果是静态资源
		if strings.Contains(path, "/js/") || strings.Contains(path, "/css/") || strings.Contains(path, "/image/") {
			if !strings.HasSuffix(path, consts.FSCompressedFileSuffix) { //过滤掉压缩包
				h.StaticFile(relativePath, path)
			}
		} else if strings.HasSuffix(path, ".html") { // 模板文件
			//创建对应的模板
			t := tmpl.New(relativePath)
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			//对应模板内容
			_, err = t.Parse(string(b))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		FuncLogError(err)
		//panic(err)
		return err
	}
	//设置模板
	h.SetHTMLTemplate(tmpl)
	return nil
}

// funcBasePath 基础路径,前端所有的资源请求必须带上 {{basePath}}
func funcBasePath() string {
	return config.BasePath
}

// 加载配置文件
func loadConfigFile() configStruct {
	defaultErr := errors.New("config.json加载失败,使用默认配置")
	// 打开文件
	jsonFile, err := os.Open(datadir + "config.json")
	if err != nil {
		FuncLogError(defaultErr)
		return defaultConfig
	}
	// 关闭文件
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	configJson := configStruct{}
	//Decode从输入流读取下一个json编码值并保存在v指向的值里
	err = json.Unmarshal([]byte(byteValue), &configJson)
	if err != nil {
		FuncLogError(defaultErr)
		return defaultConfig
	}
	return configJson
}

var defaultConfig = configStruct{
	BasePath: "",
	//默认的加密Secret
	JwtSecret:   "gpress+jwtSecret-2023",
	JwttokenKey: "jwttoken", //jwt的key
	Timeout:     1800,       //半个小时超时
}

type configStruct struct {
	BasePath    string `json:"basePath"`
	JwtSecret   string `json:"jwtSecret"`
	JwttokenKey string `json:"jwttokenKey"`
	Timeout     int    `json:"timeout"`
}
