package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"html/template"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

var tmpl *template.Template = template.New(defaultName).Delims("", "").Funcs(funcMap)

// initTemplate 初始化模板
func initTemplate() error {
	// h.SetFuncMap(funcMap)
	// h.LoadHTMLFiles(themePath + "index.html")
	// h.LoadHTMLGlob(datadir + "html/theme/default/*")
	// 手动声明template对象,自己控制文件路径,默认是使用文件名,多个文件夹会存在问题
	err := loadTemplate(false)
	// 设置模板
	h.SetHTMLTemplate(tmpl)
	// 设置默认的静态文件,实际路径会拼接为 datadir/public
	hStatic("/public", datadir)
	// 设置静态网页目录
	hStatic("/statichtml", datadir)
	return err
}

// loadTemplate 用于更新重复加载
func loadTemplate(reload bool) error {
	//声明新的template
	loadTmpl := template.New(defaultName).Delims("", "").Funcs(funcMap)

	staticFileMap := make(map[string]string)
	//遍历后台admin模板
	err := walkTemplateDir(loadTmpl, reload, templateDir+"admin/", templateDir, &staticFileMap)
	if err != nil {
		FuncLogError(err)
		return err
	}
	//遍历用户配置的主题模板
	err = walkTemplateDir(loadTmpl, reload, templateDir+"theme/"+config.Theme+"/", templateDir+"theme/"+config.Theme+"/", &staticFileMap)
	if err != nil {
		FuncLogError(err)
		return err
	}
	//此处为hertz bug,已经调用了 h.SetHTMLTemplate(tmpl),但是c.HTMLRender依然是老的内存地址.所以这里暂时不改变指针地址
	//https://github.com/cloudwego/hertz/issues/683
	*tmpl = *loadTmpl

	// 设置模板
	//h.SetHTMLTemplate(tmpl)
	if reload { //如果是reload,不处理静态文件
		return nil
	}

	//增加静态文件夹
	for k, v := range staticFileMap {
		//staticFS2 := http.Dir(v)

		hStatic(k, v)

		//h.Handle("GET", k+"/*filepath", http.FileServer(staticFS2))

	}

	/*
		// 直接映射 /statichtml,暂时不用每个都单独注册了
		// 遍历处理静态化文件
		filepath.Walk(statichtmlDir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() { // 只处理文件
				return nil
			}
			// 分隔符统一为 / 斜杠
			path = filepath.ToSlash(path)
			// 相对路径
			relativePath := path[len(statichtmlDir)-1:]
			// 设置静态化文件
			h.StaticFile(relativePath, path)
			return nil
		})
	*/
	return nil
}

func walkTemplateDir(loadTmpl *template.Template, reload bool, walkDir string, baseDir string, staticFileMap *map[string]string) error {
	//遍历模板文件夹
	err := filepath.Walk(walkDir, func(path string, info os.FileInfo, err error) error {
		// 分隔符统一为 / 斜杠
		path = filepath.ToSlash(path)
		// 相对路径

		// 如果是静态资源
		if !reload && (strings.Contains(path, "/js/") || strings.Contains(path, "/css/") || strings.Contains(path, "/image/")) {
			relativePath := path[len(baseDir)-1:]
			/*
				// 直接映射静态文件夹
				if !strings.HasSuffix(path, consts.FSCompressedFileSuffix) { // 过滤掉压缩包
				    h.StaticFile(relativePath, path)
				}
			*/
			if strings.Contains(relativePath, "/js/") { //如果是js文件夹
				key := relativePath[:strings.Index(relativePath, "/js/")+4]
				value := path[:strings.Index(path, key)]
				(*staticFileMap)[key] = value
			} else if strings.Contains(relativePath, "/css/") { //如果是css文件夹
				key := relativePath[:strings.Index(relativePath, "/css/")+5]
				value := path[:strings.Index(path, key)]
				(*staticFileMap)[key] = value
			} else if strings.Contains(relativePath, "/image/") { //如果是image文件夹
				key := relativePath[:strings.Index(relativePath, "/image/")+7]
				value := path[:strings.Index(path, key)]
				(*staticFileMap)[key] = value
			}

		} else if strings.HasSuffix(path, ".html") { // 模板文件
			relativePath := path[len(baseDir):]
			// 创建对应的模板
			t := loadTmpl.New(relativePath)
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
	return err
}

// isInstalled 是否已经安装过了
func isInstalled() bool {
	// 依赖bleveStatus变量,确保bleve在isInstalled之前初始化
	if !bleveStatus {
		FuncLogError(errors.New("bleveStatus状态为false"))
	}
	return !pathExists(templateDir + "admin/install.html")
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
	err = os.Rename(templateDir+"admin/install.html", templateDir+"admin/install.html."+now)
	if err != nil {
		return err
	}

	// install_config.json 重命名为 install_config.json_配置已失效_请通过后台设置管理
	err = os.Rename(datadir+"install_config.json", datadir+"install_config.json."+now)
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
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// hashSha256 使用sha256计算hash值
func hashSha256(str string) string {
	hashByte := sha256.Sum256([]byte(str))
	hashStr := hex.EncodeToString(hashByte[:])
	return hashStr
}

// ResponseData 返回数据包装器
type ResponseData struct {
	// 业务状态代码 // 异常 0, 成功 1,默认失败0,业务代码见说明
	StatusCode int `json:"statusCode"`
	// HttpCode http的状态码
	// HttpCode int `json:"httpCode,omitempty"`
	// 返回数据
	Data interface{} `json:"data,omitempty"`
	// 返回的信息内容,配合StatusCode
	Message string `json:"message,omitempty"`
	// 扩展的map,用于处理返回多个值的情况
	ExtMap map[string]interface{} `json:"extMap,omitempty"`
	// 列表的分页对象
	Page *Page `json:"page,omitempty"`
	// 查询条件的struct回传
	// QueryStruct interface{} `json:"queryStruct,omitempty"`
	QueryString string `json:"queryString,omitempty"`
	// 索引名称
	UrlPathIndexName string `json:"urlPathIndexName,omitempty"`
	// 索引字段信息
	// IndexField []IndexFieldStruct `json:"indexField,omitempty"`
	// 响应错误
	ERR error `json:"err,omitempty"`
}

func responData2Map(responseData ResponseData) map[string]interface{} {
	result := make(map[string]interface{}, 0)
	result["statusCode"] = responseData.StatusCode
	result["data"] = responseData.Data
	result["message"] = responseData.Message
	result["extMap"] = responseData.ExtMap
	result["page"] = responseData.Page
	result["queryString"] = responseData.QueryString
	result["urlPathIndexName"] = responseData.UrlPathIndexName
	result["err"] = responseData.ERR
	return result
}

func hStatic(relativePath, root string) {
	basePath := funcBasePath()
	filePath := ""
	if basePath == "/" || basePath == "" { //默认值
		filePath = root + relativePath
	} else if strings.HasPrefix(relativePath, basePath) { //去掉前缀
		filePath = root + relativePath[len(basePath):]
	} else {
		filePath = root + relativePath
	}
	h.StaticFS(relativePath, &app.FS{
		Root: filePath,
		PathRewrite: func(c *app.RequestContext) []byte {
			path := "/" + c.Param("filepath")
			return []byte(path)
		},
	},
	)
}

func cRedirecURI(uri string) []byte {
	return []byte(config.BasePath + uri)
}
