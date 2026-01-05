// Copyright (c) 2023 gpress Authors.
//
// This file is part of gpress.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gitee.com/chunanyong/zorm"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server/render"
)

// 后台模板渲染
var templateAdmin = template.New(appName+"-admin").Delims("", "").Funcs(funcMap)
var htmlRenderAdmin = render.HTMLProduction{Template: templateAdmin}

// 前端模板渲染,分为pc/wap/wx三种客户端,使用Map记录
var webTemplateMap = make(map[string]*template.Template)
var webRenderMap = make(map[string]render.HTMLRender)

// cHtml 渲染前端页面
func cHtml(c *app.RequestContext, code int, name string, obj interface{}) {
	_, htmlRender := getTheme(c)
	instance := htmlRender.Instance(name, obj)
	c.Render(code, instance)
}

// cHtmlAdmin 渲染后台界面
func cHtmlAdmin(c *app.RequestContext, code int, name string, obj interface{}) {
	instance := htmlRenderAdmin.Instance(name, obj)
	c.Render(code, instance)
}

// loadTemplate 加载页面模板
func loadTemplate() error {
	var err error
	site, err = funcSite()
	if err != nil {
		return err
	}

	//遍历后台admin模板
	err = walkTemplateDir(templateAdmin, templateDir+"admin/", templateDir, true)
	if err != nil {
		FuncLogError(nil, err)
		return err
	}
	// 重新初始化Template和Render的Map
	webTemplateMap = make(map[string]*template.Template)
	webRenderMap = make(map[string]render.HTMLRender)
	//遍历前端默认模板文件
	if site.Theme != "" {
		themeTemplate, has := webTemplateMap[site.Theme]
		if !has {
			themeTemplate = template.New(appName+site.Theme).Delims("", "").Funcs(funcMap)
			webTemplateMap[site.Theme] = themeTemplate
			err = walkTemplateDir(themeTemplate, themeDir+site.Theme+"/", themeDir+site.Theme+"/", false)
			if err != nil {
				FuncLogError(nil, err)
				return err
			}
		}
		htmlRender := render.HTMLProduction{Template: themeTemplate}
		webRenderMap["default"] = &htmlRender

	}
	if site.ThemePC != "" {
		themeTemplate, has := webTemplateMap[site.ThemePC]
		if !has {
			themeTemplate = template.New(appName+site.ThemePC).Delims("", "").Funcs(funcMap)
			webTemplateMap[site.ThemePC] = themeTemplate
			err = walkTemplateDir(themeTemplate, themeDir+site.ThemePC+"/", themeDir+site.ThemePC+"/", false)
			if err != nil {
				FuncLogError(nil, err)
				return err
			}
		}
		htmlRender := render.HTMLProduction{Template: themeTemplate}
		webRenderMap["pc"] = &htmlRender

	}
	//遍历手机wap的模板文件
	if site.ThemeWAP != "" {
		themeTemplate, has := webTemplateMap[site.ThemeWAP]
		if !has {
			themeTemplate = template.New(appName+site.ThemeWAP).Delims("", "").Funcs(funcMap)
			webTemplateMap[site.ThemeWAP] = themeTemplate
			err = walkTemplateDir(themeTemplate, themeDir+site.ThemeWAP+"/", themeDir+site.ThemeWAP+"/", false)
			if err != nil {
				FuncLogError(nil, err)
				return err
			}
		}
		htmlRender := render.HTMLProduction{Template: themeTemplate}
		webRenderMap["wap"] = &htmlRender

	}
	//遍历微信WX的模板文件
	if site.ThemeWX != "" {
		themeTemplate, has := webTemplateMap[site.ThemeWX]
		if !has {
			themeTemplate = template.New(appName+site.ThemeWX).Delims("", "").Funcs(funcMap)
			webTemplateMap[site.ThemeWX] = themeTemplate
			err = walkTemplateDir(themeTemplate, themeDir+site.ThemeWX+"/", themeDir+site.ThemeWX+"/", false)
			if err != nil {
				FuncLogError(nil, err)
				return err
			}
		}
		htmlRender := render.HTMLProduction{Template: themeTemplate}
		webRenderMap["wx"] = &htmlRender

	}
	return nil
}

// walkTemplateDir 遍历模板文件夹
func walkTemplateDir(tmpl *template.Template, walkDir string, baseDir string, isAdmin bool) error {
	loadTmpl := template.New(tmpl.Name()).Delims("", "").Funcs(funcMap)
	//遍历模板文件夹
	err := filepath.WalkDir(walkDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 分隔符统一为 / 斜杠
		path = filepath.ToSlash(path)
		if !isAdmin && strings.Contains(path, "/admin/") { //如果用户主题,但是包含admin目录,不解析
			return fs.SkipDir
		}
		if strings.HasSuffix(path, ".html") { // 模板文件
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

	//更新模板对象
	*tmpl = *loadTmpl

	return err
}

// isInstalled 是否已经安装过了
func isInstalled() bool {
	// 检查表状态
	var dbStatus = checkDBStatus()

	// 依赖dbStatus变量,确保sqlite在isInstalled之前初始化
	if !dbStatus {
		err := errors.New(funcT("Table check failed,dbStatus is false"))
		FuncLogError(nil, err)
		panic(err)
	}

	return !pathExist(templateDir + "admin/install.html")
}

// updateInstall 更新安装状态
func updateInstall(ctx context.Context) error {
	// 将config配置写入到表,写入前先把config表清空
	err := insertConfig(ctx)
	if err != nil {
		return err
	}
	//如果文件存在就删除
	if pathExist(templateDir + "admin/install.html.bak") {
		os.Remove(templateDir + "admin/install.html.bak")
	}
	// 删除 install 文件
	err = os.Rename(templateDir+"admin/install.html", templateDir+"admin/install.html.bak")
	if err != nil {
		return err
	}

	// 更改安装状态
	installed = true
	return nil
}

// initStaticFS 初始化静态文件
func initStaticFS() {
	appFS := &app.FS{
		Root:     "./",
		Compress: false, //不使用hertz的压缩.gz,程序控制压缩.gz
		//CompressedFileSuffix: compressedFileSuffix,
		PathRewrite: func(c *app.RequestContext) []byte {
			//relativePath := c.Param("filepath")
			relativePath := string(c.URI().Path())
			//去掉basePath
			relativePath = strings.TrimPrefix(relativePath, funcBasePath())
			key := relativePath
			parts := strings.Split(key, "/")
			if len(parts) > 1 {
				key = parts[0]
				if key == "admin" { //后台管理
					key = key + "/" + parts[1]
				}
			}
			switch key {
			case "js", "css", "image": //处理静态文件,根据浏览器获取对应的主题
				theme, _ := getTheme(c)
				return []byte("/" + themeDir + theme + "/" + relativePath)
			case "public": //public目录下的静态文件
				return []byte("/" + datadir + relativePath)
			case "favicon.ico": //默认的favicon图标
				return []byte("/" + datadir + site.Favicon)
			case "admin/js", "admin/css", "admin/image": //后台管理的静态文件
				return []byte("/" + templateDir + relativePath)
			default: //其他从public目录下获取
				return []byte("/" + datadir + "public/" + relativePath)
			}

		},
	}
	//映射静态文件,兼容项目前缀路径. 不再使用,有动态通配拦截 /
	h.StaticFS("/", appFS)
}

// cRedirecURI 重定向到uri,拼接上basePath
func cRedirecURI(uri string) []byte {
	return []byte(config.BasePath + uri)
}

// FuncGenerateStringID 默认生成字符串ID的函数.方便自定义扩展
// FuncGenerateStringID Function to generate string ID by default. Convenient for custom extension
var FuncGenerateStringID func() string = generateStringID

// generateStringID 生成主键字符串
// generateStringID Generate primary key string
func generateStringID() string {
	// 使用 crypto/rand 真随机9位数
	randNum, randErr := rand.Int(rand.Reader, big.NewInt(1000000000))
	if randErr != nil {
		return ""
	}
	// 获取9位数,前置补0,确保9位数
	rand9 := fmt.Sprintf("%09d", randNum)

	// 获取纳秒 按照 年月日时分秒毫秒微秒纳秒 拼接为长度23位的字符串
	pk := time.Now().Format("2006.01.02.15.04.05.000000000")
	pk = strings.ReplaceAll(pk, ".", "")

	// 23位字符串+9位随机数=32位字符串,这样的好处就是可以使用ID进行排序
	pk = pk + rand9
	return pk
}

// pathExist 文件或者目录是否存在
func pathExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// findThemeTemplate 从数据库查询模板文件
func findThemeTemplate(ctx context.Context, tableName string, urlPathParam string) (string, error) {
	finder := zorm.NewSelectFinder(tableName, "template_file").Append(" WHERE id=?", urlPathParam)
	templatePath := ""
	has, err := zorm.QueryRow(ctx, finder, &templatePath)
	if err != nil {
		FuncLogError(ctx, err)
		return "", err

	}
	if !has {
		return "", err
	}
	return templatePath, err
}

// getTheme 根据Cookie和User-Agent获取配置的 theme
func getTheme(c *app.RequestContext) (string, render.HTMLRender) {
	//优先cookie,其次请求头
	userAgentByte := c.Cookie("User-Agent")
	if len(userAgentByte) == 0 {
		userAgentByte = c.GetHeader("User-Agent")
	}
	if len(userAgentByte) == 0 {
		return site.Theme, webRenderMap["default"]
	}
	userAgent := strings.ToLower(string(userAgentByte))

	if site.ThemeWX != "" && (strings.Contains(userAgent, "weixin") || strings.Contains(userAgent, "wechat") || strings.Contains(userAgent, "micromessenger")) { // 微信
		return site.ThemeWX, webRenderMap["wx"]
	} else if site.ThemeWAP != "" && (strings.Contains(userAgent, "android") || strings.Contains(userAgent, "phone") || strings.Contains(userAgent, "harmonyos") || strings.Contains(userAgent, "mobile") || strings.Contains(userAgent, "blackberry") || strings.Contains(userAgent, "ipod")) {
		return site.ThemeWAP, webRenderMap["wap"]
	} else if site.ThemePC != "" && (strings.Contains(userAgent, "windows") || strings.Contains(userAgent, "linux") || strings.Contains(userAgent, "macintosh") || strings.Contains(userAgent, "ipad") || strings.Contains(userAgent, "tablet")) {
		return site.ThemePC, webRenderMap["pc"]
	}
	return site.Theme, webRenderMap["default"]

}

// ResponseData 返回数据包装器
type ResponseData struct {
	// StatusCode 业务状态代码 // 异常 0, 成功 1,默认失败0,业务代码见说明
	StatusCode int `json:"statusCode"`

	// Data 返回数据
	Data interface{} `json:"data,omitempty"`

	// Message 返回的信息内容,配合StatusCode
	Message string `json:"message,omitempty"`

	// ExtMap 扩展的map,用于处理返回多个值的情况
	ExtMap map[string]interface{} `json:"extMap,omitempty"`

	// 列表的分页对象
	Page *zorm.Page `json:"page,omitempty"`

	// QueryStringMap 查询条件的struct回传
	QueryStringMap map[string]string `json:"queryStringMap,omitempty"`

	// UrlPathParam 表名称
	UrlPathParam string `json:"urlPathParam,omitempty"`

	// ERR 响应错误, json无法序列化error类型,使用Message返回错误信息
	//ERR error `json:"err,omitempty"`
}
