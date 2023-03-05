package main

import (
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// funcT 多语言i18n适配,例如 {{ T "nextPage" }}
func funcT(key string) (string, error) {
	return key, nil
}

// funcSafeHTML 转义html字符串
func funcSafeHTML(html string) (string, error) {
	ss := template.HTMLEscapeString(html)
	return ss, nil
}

// funcRelURL 真实的url
func funcRelURL(url string) (string, error) {
	return funcSafeHTML(url)
}

// funcSass 编译sass,生成css
func funcSass(sassFile string) (string, error) {
	sassPath := templateDir + "theme/" + config.Theme + "/assets/" + sassFile
	pathHash := hashSha256(sassPath)

	//文件夹路径
	filePath := templateDir + "theme/" + config.Theme + "/css/" + pathHash + ".css"

	//url 访问路径
	fileUrl := themePath + "css/" + pathHash + ".css"

	_, err := os.Stat(filePath)
	if os.IsExist(err) { //如果文件已经存在了,直接返回
		return funcRelURL(fileUrl)
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

	//c := exec.Command("bash", "-c", cmd)  // mac or linux
	err = cmd.Run()
	if err != nil {
		FuncLogError(err)
	}

	return funcSafeHTML(fileUrl)
}
