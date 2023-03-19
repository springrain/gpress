package route

import (
	"gitee.com/gpress/gpress/constant"
	"gitee.com/gpress/gpress/logger"
	"gitee.com/gpress/gpress/service"
	"gitee.com/gpress/gpress/util"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func funcT() error {
	h.SetFuncMap(funcMap)
	// h.LoadHTMLFiles(themePath + "index.html")
	// h.LoadHTMLGlob(datadir + "html/theme/default/*")
	// 手动声明template对象,自己控制文件路径,默认是使用文件名,多个文件夹会存在问题
	err := loadTemplate()
	// 设置默认的静态文件,实际路径会拼接为 datadir/public
	h.Static("/public", constant.DATA_DIR)
	return err
}

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
	sassPath := constant.TEMPLATE_DIR + "theme/" + service.Config.Theme + "/assets/" + sassFile
	pathHash := util.HashSha256(sassPath)

	//生成的css路径
	filePath := constant.TEMPLATE_DIR + "theme/" + service.Config.Theme + "/css/" + pathHash + ".css"
	//filePath = filepath.FromSlash(filePath)

	//url 访问路径
	fileUrl := ThemePath + "css/" + pathHash + ".css"

	_, err := os.Lstat(filePath)
	if !os.IsNotExist(err) { //如果文件已经存在了,直接返回
		return funcSafeHTML(fileUrl)
	}
	var cmd *exec.Cmd
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	if goos == "windows" {
		cmdStr := constant.DATA_DIR + "dart-sass/" + goos + "-" + goarch + "/sass.bat --style=compressed --charset --no-source-map " + sassPath + ":" + filePath
		// 分隔符统一系统符号
		cmdStr = filepath.FromSlash(cmdStr)
		cmd = exec.Command("cmd", "/C", cmdStr) // windows
	} else if goos == "linux" || goos == "darwin" {
		cmdStr := constant.DATA_DIR + "dart-sass/" + goos + "-" + goarch + "/sass --style=compressed --charset --no-source-map " + sassPath + ":" + filePath
		cmd = exec.Command("bash", "-c", cmdStr) // mac or linux
	}
	err = cmd.Run()
	if err != nil {
		logger.FuncLogError(err)
	}
	fileUrl, err = funcSafeHTML(fileUrl)
	//增加静态资源映射
	h.StaticFile(fileUrl, filePath)
	return fileUrl, err
}
