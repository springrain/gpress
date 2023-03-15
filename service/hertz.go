package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"gitee.com/gpress/gpress/bleves"
	"gitee.com/gpress/gpress/configs"
	"gitee.com/gpress/gpress/hugo"
	"gitee.com/gpress/gpress/logger"
	"gitee.com/gpress/gpress/util"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"html/template"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// hertz对象,可以在其他地方使用
var h = server.Default(server.WithHostPorts(config.ServerPort), server.WithBasePath("/"))

var funcMap = template.FuncMap{"md5": util.FuncMD5, "basePath": funcBasePath, "T": hugo.FuncT, "safeHTML": hugo.FuncSafeHTML, "relURL": hugo.FuncRelURL, "sass": funcSass}

var installed = isInstalled()

var themePath = "/theme/" + config.Theme + "/"

// 检查索引状态
var bleveStatus = bleves.CheckBleveStatus()

// initTemplate 初始化模板
func initTemplate() error {
	h.SetFuncMap(funcMap)
	// h.LoadHTMLFiles(themePath + "index.html")
	// h.LoadHTMLGlob(datadir + "html/theme/default/*")
	// 手动声明template对象,自己控制文件路径,默认是使用文件名,多个文件夹会存在问题
	err := loadTemplate()
	// 设置默认的静态文件,实际路径会拼接为 datadir/public
	h.Static("/public", configs.DATA_DIR)
	return err
}

// loadTemplate 用于更新重复加载
func loadTemplate() error {
	tmpl := template.New("").Delims("", "").Funcs(funcMap)
	err := filepath.Walk(configs.TEMPLATE_DIR, func(path string, info os.FileInfo, err error) error {
		// 分隔符统一为 / 斜杠
		path = filepath.ToSlash(path)
		// 相对路径
		relativePath := path[len(configs.TEMPLATE_DIR)-1:]
		// 如果是静态资源
		if strings.Contains(path, "/js/") || strings.Contains(path, "/css/") || strings.Contains(path, "/image/") {
			if !strings.HasSuffix(path, consts.FSCompressedFileSuffix) { // 过滤掉压缩包
				h.StaticFile(relativePath, path)
			}
		} else if strings.HasSuffix(path, ".html") { // 模板文件
			// 创建对应的模板
			t := tmpl.New(relativePath)
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			// 对应模板内容
			_, err = t.Parse(string(b))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		logger.FuncLogError(err)
		// panic(err)
		return err
	}

	// 处理静态化文件
	err = filepath.Walk(configs.STATIC_HTML_DIR, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() { // 只处理文件
			return nil
		}
		// 分隔符统一为 / 斜杠
		path = filepath.ToSlash(path)
		// 相对路径
		relativePath := path[len(configs.STATIC_HTML_DIR)-1:]
		// 设置静态化文件
		h.StaticFile(relativePath, path)
		return nil
	})
	if err != nil {
		logger.FuncLogError(err)
		return err
	}

	// 设置模板
	h.SetHTMLTemplate(tmpl)
	return nil
}

// funcBasePath 基础路径,前端所有的资源请求必须带上 {{basePath}}
func funcBasePath() string {
	return config.basePath
}

// isInstalled 是否已经安装过了
func isInstalled() bool {
	// 依赖bleveStatus变量,确保bleve在isInstalled之前初始化
	if !bleveStatus {
		logger.FuncLogError(errors.New("bleveStatus状态为false"))
	}
	_, err := os.Lstat(configs.TEMPLATE_DIR + "admin/install.html")
	return os.IsNotExist(err)
}

// updateInstall 更新安装状态
func updateInstall(ctx context.Context) error {
	// 将config配置写入到索引,写入前先把config表清空
	err := insertConfig(ctx, config)
	if err != nil {
		return err
	}

	now := strconv.FormatInt(time.Now().UnixNano(), 10)
	// 删除 install 文件
	err = os.Rename(configs.TEMPLATE_DIR+"admin/install.html", configs.TEMPLATE_DIR+"admin/install.html."+now)
	if err != nil {
		return err
	}

	// install_config.json 重命名为 install_config.json_配置已失效_请通过后台设置管理
	err = os.Rename(configs.DATA_DIR+"install_config.json", configs.DATA_DIR+"install_config.json."+now)
	if err != nil {
		return err
	}
	// 更改安装状态
	installed = true
	return nil
}

// randStr 生成随机字符串
func randStr(n int) string {
	//rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = configs.LETTERS[rand.Intn(len(configs.LETTERS))]
	}
	return string(b)
}

// hashSha256 使用sha256计算hash值
func hashSha256(str string) string {
	hashByte := sha256.Sum256([]byte(str))
	hashStr := hex.EncodeToString(hashByte[:])
	return hashStr
}

// funcSass 编译sass,生成css
func funcSass(sassFile string) (string, error) {
	sassPath := configs.TEMPLATE_DIR + "theme/" + config.Theme + "/assets/" + sassFile
	pathHash := hashSha256(sassPath)

	//生成的css路径
	filePath := configs.TEMPLATE_DIR + "theme/" + config.Theme + "/css/" + pathHash + ".css"
	//filePath = filepath.FromSlash(filePath)

	//url 访问路径
	fileUrl := themePath + "css/" + pathHash + ".css"

	_, err := os.Lstat(filePath)
	if !os.IsNotExist(err) { //如果文件已经存在了,直接返回
		return hugo.FuncSafeHTML(fileUrl)
	}
	var cmd *exec.Cmd
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	if goos == "windows" {
		cmdStr := configs.DATA_DIR + "dart-sass/" + goos + "-" + goarch + "/sass.bat --style=compressed --charset --no-source-map " + sassPath + ":" + filePath
		// 分隔符统一系统符号
		cmdStr = filepath.FromSlash(cmdStr)
		cmd = exec.Command("cmd", "/C", cmdStr) // windows
	} else if goos == "linux" || goos == "darwin" {
		cmdStr := configs.DATA_DIR + "dart-sass/" + goos + "-" + goarch + "/sass --style=compressed --charset --no-source-map " + sassPath + ":" + filePath
		cmd = exec.Command("bash", "-c", cmdStr) // mac or linux
	}
	err = cmd.Run()
	if err != nil {
		logger.FuncLogError(err)
	}
	fileUrl, err = hugo.FuncSafeHTML(fileUrl)
	//增加静态资源映射
	h.StaticFile(fileUrl, filePath)
	return fileUrl, err
}
