package main

import (
	"context"
	"html/template"
	"os/exec"
	"path/filepath"
	"runtime"
)

var funcMap = template.FuncMap{

	"basePath":   funcBasePath,
	"T":          funcT,
	"safeHTML":   funcSafeHTML,
	"relURL":     funcRelURL,
	"indexFiled": funcIndexFiled,
	//"md5":      funcMD5,
	//"sass":       funcSass,
	//"themePath":  funcThemePath,
	//"themeFile":  funcThemeFile,

}

// funcBasePath 基础路径,前端所有的资源请求必须带上 {{basePath}}
func funcBasePath() string {
	return config.BasePath
}

/*
func funcThemePath() string {
	return themePath
}

func funcThemeFile(file string) string {
	return "/theme/" + config.Theme + "/" + file
}

// 测试自定义函数
func funcMD5(in string) ([]string, error) {
	list := make([]string, 2)

	hash := md5.Sum([]byte(in))
	list[0] = in
	list[1] = hex.EncodeToString(hash[:])
	return list, nil
}
*/
// funcT 多语言i18n适配,例如 {{ T "nextPage" }}
func funcT(key string) (string, error) {
	return key, nil
}

// funcSafeHTML 转义html字符串
func funcSafeHTML(html string) (string, error) {
	ss := template.HTMLEscapeString(html)
	return ss, nil
}

// funcRelURL 拼接url路径的
func funcRelURL(url string) (string, error) {
	return funcSafeHTML(themePath + url)
}

// funcSass 编译sass,生成css
func funcSass(sassFile string) (string, error) {
	sassPath := templateDir + "theme/" + config.Theme + "/assets/" + sassFile
	pathHash := hashSha256(sassPath)

	//生成的css路径
	filePath := templateDir + "theme/" + config.Theme + "/css/" + pathHash + ".css"
	//filePath = filepath.FromSlash(filePath)

	//url 访问路径
	fileUrl := themePath + "css/" + pathHash + ".css"

	if pathExists(filePath) { //如果文件已经存在了,直接返回
		return funcSafeHTML(fileUrl)
	}
	var cmd *exec.Cmd
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	if goos == "windows" {
		cmdStr := datadir + "dart-sass/" + goos + "-" + goarch + "/sass.bat --style=compressed --charset --no-source-map " + sassPath + ":" + filePath
		// 分隔符统一系统符号
		cmdStr = filepath.FromSlash(cmdStr)
		cmd = exec.Command("cmd", "/C", cmdStr) // windows
	} else if goos == "linux" || goos == "darwin" {
		cmdStr := datadir + "dart-sass/" + goos + "-" + goarch + "/sass --style=compressed --charset --no-source-map " + sassPath + ":" + filePath
		cmd = exec.Command("bash", "-c", cmdStr) // mac or linux
	}
	err := cmd.Run()
	if err != nil {
		FuncLogError(err)
		return fileUrl, err
	}
	fileUrl, err = funcSafeHTML(fileUrl)
	//增加静态资源映射,使用了static文件夹,不需要再映射了
	//h.StaticFile(fileUrl, filePath)
	return fileUrl, err
}

// funcIndexFiled 根据indexName查找字段
func funcIndexFiled(indexName string) ([]IndexFieldStruct, error) {
	ctx := context.Background()
	indexField, err := findIndexFieldStruct(ctx, indexName)
	return indexField, err
}
