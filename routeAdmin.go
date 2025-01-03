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
	"encoding/hex"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	//"github.com/bytedance/go-tagexpr/v2/binding"
	"gitee.com/chunanyong/zorm"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/route/param"
	"golang.org/x/crypto/sha3"
)

// alphaNumericReg 传入的列名只能是字母数字或下划线,长度不超过20
var alphaNumericReg = regexp.MustCompile("^[a-zA-Z0-9_]{1,20}$")

// init 初始化函数
func init() {
	// adminGroup 初始化管理员路由组
	var adminGroup = h.Group("/admin")
	//设置权限
	adminGroup.Use(permissionHandler())

	//设置json处理函数
	//binding.ResetJSONUnmarshaler(json.Unmarshal)
	/*
		binding.Default().ResetJSONUnmarshaler(func(data []byte, v interface{}) error {
			dec := json.NewDecoder(bytes.NewBuffer(data))
			dec.UseNumber()
			return dec.Decode(v)
		})
	*/

	// 异常页面
	h.GET("/admin/error", func(ctx context.Context, c *app.RequestContext) {
		cHtmlAdmin(c, http.StatusOK, "admin/error.html", nil)
	})

	// 安装
	h.GET("/admin/install", funcAdminInstallPre)
	h.POST("/admin/install", funcAdminInstall)

	// 生成30位随机数,钱包签名随机校验.如果32位metamask会解析成16进制字符串,可能是metamask的bug
	h.POST("/admin/random", func(ctx context.Context, c *app.RequestContext) {
		generateChainRandStr()
		c.JSON(http.StatusOK, ResponseData{StatusCode: 1, Data: chainRandStr})
	})

	// 后台管理员登录
	h.GET("/admin/login", funcAdminLoginPre)
	h.POST("/admin/login", funcAdminLogin)

	// 后台管理员使用区块链账号登录
	h.GET("/admin/chainlogin", funcAdminChainloginPre)
	h.POST("/admin/chainlogin", funcAdminChainlogin)

	// 后台管理员首页
	adminGroup.GET("/index", func(ctx context.Context, c *app.RequestContext) {
		cHtmlAdmin(c, http.StatusOK, "admin/index.html", nil)
	})

	// 刷新站点,重新加载资源包含模板和对应的静态文件
	adminGroup.GET("/reload", funcAdminReload)

	//上传文件
	adminGroup.POST("/upload", funcUploadFile)
	//上传主题文件
	adminGroup.POST("/themeTemplate/uploadTheme", funcUploadTheme)

	// 通用list列表
	adminGroup.GET("/:urlPathParam/list", funcList)
	// 查询主题模板
	adminGroup.GET("/themeTemplate/list", funcListThemeTemplate)
	// 查询Content列表,根据CategoryId like
	adminGroup.GET("/content/list", funcContentList)

	// 通用查看
	adminGroup.GET("/:urlPathParam/look", funcLook)
	// 内容预览
	adminGroup.GET("/content/look", funcContentPreview)
	// 导航菜单预览
	adminGroup.GET("/category/look", funcCategoryPreview)

	//跳转到修改页面
	adminGroup.GET("/:urlPathParam/update", funcUpdatePre)
	// 修改Config
	adminGroup.POST("/config/update", funcUpdateConfig)
	// 修改Site
	adminGroup.POST("/site/update", funcUpdateSite)
	// 修改User
	adminGroup.POST("/user/update", funcUpdateUser)
	// 修改Category
	adminGroup.POST("/category/update", funcUpdateCategory)
	// 修改Content
	adminGroup.POST("/content/update", funcUpdateContent)
	// 修改ThemeTemplate
	adminGroup.POST("/themeTemplate/update", funcUpdateThemeTemplate)

	//跳转到保存页面
	adminGroup.GET("/:urlPathParam/save", funcSavePre)
	//保存Category
	adminGroup.POST("/category/save", funcSaveCategory)
	//保存Content
	adminGroup.POST("/content/save", funcSaveContent)

	//ajax POST删除数据
	adminGroup.POST("/:urlPathParam/delete", funcDelete)

}

// funcAdminInstallPre 跳转到安装界面
func funcAdminInstallPre(ctx context.Context, c *app.RequestContext) {
	if installed { // 如果已经安装过了,跳转到登录
		c.Redirect(http.StatusOK, cRedirecURI("admin/login"))
		c.Abort() // 终止后续调用
		return
	}
	cHtmlAdmin(c, http.StatusOK, "admin/install.html", nil)
}

// funcAdminInstall 后台安装
func funcAdminInstall(ctx context.Context, c *app.RequestContext) {
	if installed { // 如果已经安装过了,跳转到登录
		c.Redirect(http.StatusOK, cRedirecURI("admin/login"))
		c.Abort() // 终止后续调用
		return
	}
	// 使用后端管理界面配置,jwtSecret也有后端随机产生
	user := User{}
	user.Account = c.PostForm("account")
	user.UserName = c.PostForm("account")
	user.Password = c.PostForm("password")
	user.ChainType = c.PostForm("chainType")
	user.ChainAddress = c.PostForm("chainAddress")
	// 重新hash密码,避免拖库后撞库
	sha3Bytes := sha3.Sum512([]byte(user.Password))
	user.Password = hex.EncodeToString(sha3Bytes[:])

	loginHtml := "admin/login?message=恭喜您,成功安装gpress,现在请登录"
	if user.ChainAddress != "" && user.ChainType != "" { //如果使用了address作为登录方式
		user.Account = ""
		user.UserName = ""
		loginHtml = "admin/chainlogin"
	}
	err := insertUser(ctx, user)
	if err != nil {
		c.Redirect(http.StatusOK, cRedirecURI("admin/error"))
		c.Abort() // 终止后续调用
		return
	}
	// 安装成功,更新安装状态
	updateInstall(ctx)
	c.Redirect(http.StatusOK, cRedirecURI(loginHtml))
}

// funcAdminLoginPre 跳转到登录界面
func funcAdminLoginPre(ctx context.Context, c *app.RequestContext) {
	if !installed { // 如果没有安装,跳转到安装
		c.Redirect(http.StatusOK, cRedirecURI("admin/install"))
		c.Abort() // 终止后续调用
		return
	}
	responseData := make(map[string]string, 0)
	message, ok := c.GetQuery("message")
	if ok {
		responseData["message"] = message
	}
	if errorLoginCount.Load() >= errCount { //连续错误3次显示验证码
		responseData["showCaptcha"] = "1"
		generateCaptcha()
		responseData["captchaBase64"] = captchaBase64
	}
	c.SetCookie(config.JwttokenKey, "", config.Timeout, "/", "", protocol.CookieSameSiteStrictMode, false, true)
	cHtmlAdmin(c, http.StatusOK, "admin/login.html", responseData)
}

// funcAdminLogin 后台登录
func funcAdminLogin(ctx context.Context, c *app.RequestContext) {
	if !installed { // 如果没有安装,跳转到安装
		c.Redirect(http.StatusOK, cRedirecURI("admin/install"))
		c.Abort() // 终止后续调用
		return
	}

	if errorLoginCount.Load() >= errCount { //连续错误3次显示验证码
		answer := c.PostForm("answer")
		if answer != captchaAnswer { //答案不对
			c.Redirect(http.StatusOK, cRedirecURI("admin/login?message=验证码错误"))
			c.Abort() // 终止后续调用
			return
		}
	}

	account := strings.TrimSpace(c.PostForm("account"))
	password := strings.TrimSpace(c.PostForm("password"))
	if account == "" || password == "" { // 用户不存在或者异常
		c.Redirect(http.StatusOK, cRedirecURI("admin/login?message=账户或密码不能为空"))
		c.Abort() // 终止后续调用
		return
	}
	// 重新hash密码,避免拖库后撞库
	sha3Bytes := sha3.Sum512([]byte(password))
	password = hex.EncodeToString(sha3Bytes[:])

	userId, err := findUserId(ctx, account, password)
	if userId == "" || err != nil { // 用户不存在或者异常
		errorLoginCount.Add(1)
		c.Redirect(http.StatusOK, cRedirecURI("admin/login?message=账户或密码错误"))
		c.Abort() // 终止后续调用
		return
	}
	jwttoken, _ := newJWTToken(userId)
	// c.HTML(http.StatusOK, "admin/index.html", nil)
	c.SetCookie(config.JwttokenKey, jwttoken, config.Timeout, "/", "", protocol.CookieSameSiteStrictMode, false, true)
	errorLoginCount.Store(0)
	c.Redirect(http.StatusOK, cRedirecURI("admin/index"))
}

// funcAdminChainloginPre 跳转到区块链登录页面
func funcAdminChainloginPre(ctx context.Context, c *app.RequestContext) {
	if !installed { // 如果没有安装,跳转到安装
		c.Redirect(http.StatusOK, cRedirecURI("admin/install"))
		c.Abort() // 终止后续调用
		return
	}
	var responseData map[string]string = nil
	message, ok := c.GetQuery("message")
	if ok {
		responseData = make(map[string]string, 0)
		responseData["message"] = message
	}
	c.SetCookie(config.JwttokenKey, "", config.Timeout, "/", "", protocol.CookieSameSiteStrictMode, false, true)
	cHtmlAdmin(c, http.StatusOK, "admin/chainlogin.html", responseData)
}

// funcAdminChainlogin 区块链登录
func funcAdminChainlogin(ctx context.Context, c *app.RequestContext) {
	if !installed { // 如果没有安装,跳转到安装
		c.Redirect(http.StatusOK, cRedirecURI("admin/install"))
		c.Abort() // 终止后续调用
		return
	}
	//获取签名
	signature := c.PostForm("signature")
	userId, chainType, chainAddress, err := findUserAddress(ctx)
	if userId == "" || chainType == "" || chainAddress == "" || err != nil {
		c.Redirect(http.StatusOK, cRedirecURI("admin/chainlogin?message=地址异常"))
		c.Abort() // 终止后续调用
		return
	}
	verify := false
	switch chainType {
	case "ETH":
		verify, err = verifySecp256k1Signature(chainAddress, chainRandStr, signature)
	case "XUPER":
		verify, err = verifyXuperSignature(chainAddress, []byte(signature), []byte(chainRandStr))
	default:
		c.Redirect(http.StatusOK, cRedirecURI("admin/chainlogin?message=暂不支持此类型区块链账户"))
		c.Abort() // 终止后续调用
		return
	}

	if !verify || err != nil {
		c.Redirect(http.StatusOK, cRedirecURI("admin/chainlogin?message=签名校验失败"))
		c.Abort() // 终止后续调用
		return
	}
	jwttoken, _ := newJWTToken(userId)
	c.SetCookie(config.JwttokenKey, jwttoken, config.Timeout, "/", "", protocol.CookieSameSiteStrictMode, false, true)
	c.Redirect(http.StatusOK, cRedirecURI("admin/content/list"))
}

// funcAdminReload 刷新站点,会重新加载模板文件,生成静态文件
func funcAdminReload(ctx context.Context, c *app.RequestContext) {
	err := loadTemplate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, ERR: err})
		c.Abort() // 终止后续调用
		return
	}
	//重新生成静态文件
	go genStaticFile()
	c.JSON(http.StatusOK, ResponseData{StatusCode: 1})
}

// funcUploadFile 上传文件
func funcUploadFile(ctx context.Context, c *app.RequestContext) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, ERR: err})
		c.Abort() // 终止后续调用
		return
	}
	path := "public/upload/" + zorm.FuncGenerateStringID(ctx) + filepath.Ext(fileHeader.Filename)
	newFileName := datadir + path
	err = c.SaveUploadedFile(fileHeader, newFileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, ERR: err})
		c.Abort() // 终止后续调用
		return
	}
	c.JSON(http.StatusOK, ResponseData{StatusCode: 1, Data: path})
}

// funcUploadTheme 上传主题
func funcUploadTheme(ctx context.Context, c *app.RequestContext) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, ERR: err})
		c.Abort() // 终止后续调用
		return
	}
	ext := filepath.Ext(fileHeader.Filename)
	if ext != ".zip" { //不是zip
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, ERR: err})
		c.Abort() // 终止后续调用
		return
	}
	path := themeDir + fileHeader.Filename
	err = c.SaveUploadedFile(fileHeader, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, ERR: err})
		c.Abort() // 终止后续调用
		return
	}
	defer func() {
		_ = os.Remove(path)
	}()
	//解压压缩包
	err = unzip(path, themeDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, ERR: err})
		c.Abort() // 终止后续调用
		return
	}
	c.JSON(http.StatusOK, ResponseData{StatusCode: 1, Data: path})
}

// funcList 通用list列表
func funcList(ctx context.Context, c *app.RequestContext) {
	urlPathParam := c.Param("urlPathParam")
	//获取页码
	pageNoStr := c.DefaultQuery("pageNo", "1")
	pageNo, _ := strconv.Atoi(pageNoStr)
	q := strings.TrimSpace(c.Query("q"))
	mapParams := make(map[string]interface{}, 0)
	//获取所有的参数
	c.Bind(&mapParams)
	//删除掉固定的两个
	delete(mapParams, "pageNo")
	delete(mapParams, "q")
	where := " WHERE 1=1 "
	var values []interface{} = make([]interface{}, 0)
	for k := range mapParams {
		if !alphaNumericReg.MatchString(k) {
			continue
		}
		value := c.Query(k)
		if strings.TrimSpace(value) == "" {
			continue
		}
		where = where + " and " + k + "=? "
		values = append(values, value)
	}
	sql := "* from " + urlPathParam + where + " order by sortNo desc "
	var responseData ResponseData
	var err error
	if len(values) == 0 {
		responseData, err = funcSelectList(urlPathParam, q, pageNo, defaultPageSize, sql)
	} else {
		responseData, err = funcSelectList(urlPathParam, q, pageNo, defaultPageSize, sql, values)
	}
	responseData.UrlPathParam = urlPathParam
	if err != nil {
		c.Redirect(http.StatusOK, cRedirecURI("admin/error"))
		c.Abort() // 终止后续调用
		return
	}
	listFile := "admin/" + urlPathParam + "/list.html"
	cHtmlAdmin(c, http.StatusOK, listFile, responseData)
}

// funcLook 通用查看,根据id查看
func funcLook(ctx context.Context, c *app.RequestContext) {
	funcLookById(ctx, c, "look.html")
}

// funcContentPreview 内容预览
func funcContentPreview(ctx context.Context, c *app.RequestContext) {
	id := c.Query("id")
	if id == "" {
		c.Redirect(http.StatusOK, cRedirecURI("admin/error"))
		c.Abort() // 终止后续调用
		return
	}
	params := make([]param.Param, 0, 1)
	params = append(params, param.Param{
		Key:   "urlPathParam",
		Value: id,
	})
	c.Params = params
	funcOneContent(ctx, c)
}

// funcCategoryPreview 导航菜单预览
func funcCategoryPreview(ctx context.Context, c *app.RequestContext) {
	id := c.Query("id")
	if id == "" {
		c.Redirect(http.StatusOK, cRedirecURI("admin/error"))
		c.Abort() // 终止后续调用
		return
	}
	params := make([]param.Param, 0, 1)
	params = append(params, param.Param{
		Key:   "urlPathParam",
		Value: id,
	})
	c.Params = params

	funcListCategory(ctx, c)
}

// funcContentList 查询Content列表,根据CategoryId like
func funcContentList(ctx context.Context, c *app.RequestContext) {
	urlPathParam := "content"
	//获取页码
	pageNoStr := c.DefaultQuery("pageNo", "1")
	q := strings.TrimSpace(c.Query("q"))
	pageNo, _ := strconv.Atoi(pageNoStr)
	id := strings.TrimSpace(c.Query("id"))
	values := make([]interface{}, 0)
	sql := ""
	if id != "" {
		sql = " * from content where id like ?  order by sortNo desc "
		values = append(values, id+"%")
	} else {
		sql = " * from content order by sortNo desc "
	}
	var responseData ResponseData
	var err error
	if len(values) == 0 {
		responseData, err = funcSelectList(urlPathParam, q, pageNo, defaultPageSize, sql)
	} else {
		responseData, err = funcSelectList(urlPathParam, q, pageNo, defaultPageSize, sql, values)
	}
	responseData.UrlPathParam = urlPathParam
	if err != nil {
		c.Redirect(http.StatusOK, cRedirecURI("admin/error"))
		c.Abort() // 终止后续调用
		return
	}
	listFile := "admin/" + urlPathParam + "/list.html"
	cHtmlAdmin(c, http.StatusOK, listFile, responseData)
}

// funcListThemeTemplate 所有的主题文件列表
func funcListThemeTemplate(ctx context.Context, c *app.RequestContext) {
	urlPathParam := "themeTemplate"
	var responseData ResponseData
	extMap := make(map[string]interface{})
	extMap["file"] = ""
	responseData.ExtMap = extMap
	list := make([]ThemeTemplate, 0)

	//遍历当前使用的模板文件夹
	err := filepath.Walk(themeDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 分隔符统一为 / 斜杠
		path = filepath.ToSlash(path)
		path = path[strings.Index(path, themeDir)+len(themeDir):]
		if path == "" {
			return err
		}
		//获取文件后缀
		ext := filepath.Ext(path)
		ext = strings.ToLower(ext)
		// 跳过压缩的 gz文件
		if ext == ".gz" {
			return nil
		}

		pid := filepath.ToSlash(filepath.Dir(path))
		if pid == "." {
			pid = ""
		}

		themeTemplate := ThemeTemplate{}
		themeTemplate.FilePath = path
		themeTemplate.Pid = pid
		themeTemplate.Id = path
		themeTemplate.FileSuffix = ext
		themeTemplate.Name = info.Name()
		if info.IsDir() {
			themeTemplate.FileType = "dir"
		} else {
			themeTemplate.FileType = "file"
		}
		list = append(list, themeTemplate)
		return nil
	})

	responseData.UrlPathParam = urlPathParam
	responseData.Data = list
	responseData.ERR = err
	listFile := "admin/" + urlPathParam + "/list.html"

	filePath := c.Query("file")
	if filePath == "" || strings.Contains(filePath, "..") {
		//c.HTML(http.StatusOK, listFile, responseData)
		cHtmlAdmin(c, http.StatusOK, listFile, responseData)
		return
	}
	filePath = filepath.ToSlash(filePath)
	fileContent, err := os.ReadFile(themeDir + filePath)
	if err != nil {
		responseData.ERR = err
		cHtmlAdmin(c, http.StatusOK, listFile, responseData)
		return
	}
	responseData.ExtMap["file"] = string(fileContent)
	cHtmlAdmin(c, http.StatusOK, listFile, responseData)
}

// funcUpdateThemeTemplate 更新主题模板
func funcUpdateThemeTemplate(ctx context.Context, c *app.RequestContext) {
	themeTemplate := ThemeTemplate{}
	c.Bind(&themeTemplate)
	filePath := filepath.ToSlash(themeTemplate.FilePath)
	if filePath == "" || strings.Contains(filePath, "..") {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0})
		c.Abort() // 终止后续调用
		return
	}
	if !pathExist(themeDir + filePath) {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0})
		c.Abort() // 终止后续调用
		return
	}

	//打开文件
	file, err := os.OpenFile(themeDir+filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0})
		c.Abort() // 终止后续调用
		return
	}
	defer file.Close() // 确保在函数结束时关闭文件

	// 写入内容
	_, err = file.WriteString(themeTemplate.FileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0})
		c.Abort() // 终止后续调用
		return
	}
	c.JSON(http.StatusOK, ResponseData{StatusCode: 1})
}

// funcUpdatePre 跳转到修改页面
func funcUpdatePre(ctx context.Context, c *app.RequestContext) {
	funcLookById(ctx, c, "update.html")
}

// funcUpdateConfig 更新配置
func funcUpdateConfig(ctx context.Context, c *app.RequestContext) {
	now := time.Now().Format("2006-01-02 15:04:05")
	entity := &Config{}
	ok := funcUpdateInit(ctx, c, entity)
	if !ok {
		return
	}
	if !strings.HasPrefix(entity.BasePath, "/") {
		entity.BasePath = "/" + entity.BasePath
	}
	if !strings.HasSuffix(entity.BasePath, "/") {
		entity.BasePath = entity.BasePath + "/"
	}
	entity.UpdateTime = now
	funcUpdate(ctx, c, entity, entity.Id)
}

// funcUpdateSite 更新站点
func funcUpdateSite(ctx context.Context, c *app.RequestContext) {
	now := time.Now().Format("2006-01-02 15:04:05")
	entity := &Site{}
	ok := funcUpdateInit(ctx, c, entity)
	if !ok {
		return
	}
	entity.UpdateTime = now
	funcUpdate(ctx, c, entity, entity.Id)
}

// funcUpdateUser 更新用户信息
func funcUpdateUser(ctx context.Context, c *app.RequestContext) {
	now := time.Now().Format("2006-01-02 15:04:05")
	entity := &User{}
	ok := funcUpdateInit(ctx, c, entity)
	if !ok {
		return
	}
	if entity.Password != "" {
		// 重新hash密码,避免拖库后撞库
		sha3Bytes := sha3.Sum512([]byte(entity.Password))
		entity.Password = hex.EncodeToString(sha3Bytes[:])
	} else {
		f1 := zorm.NewSelectFinder(tableUserName, "password").Append("WHERE id=?", entity.Id)
		password := ""
		zorm.QueryRow(ctx, f1, &password)
		entity.Password = password
	}
	entity.UpdateTime = now
	funcUpdate(ctx, c, entity, entity.Id)
}

// funcUpdateCategory 更新导航菜单
func funcUpdateCategory(ctx context.Context, c *app.RequestContext) {
	entity := &Category{}
	ok := funcUpdateInit(ctx, c, entity)
	if !ok {
		return
	}
	funcUpdate(ctx, c, entity, entity.Id)
}

// funcUpdateContent 更新内容
func funcUpdateContent(ctx context.Context, c *app.RequestContext) {
	now := time.Now().Format("2006-01-02 15:04:05")
	entity := &Content{}
	ok := funcUpdateInit(ctx, c, entity)
	if !ok {
		return
	}
	if entity.Markdown != "" {
		content, toc, err := renderMarkdownHtml(entity.Markdown)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "markdown转html错误"})
			c.Abort() // 终止后续调用
			FuncLogError(ctx, err)
			return
		}
		entity.Content = content
		entity.Toc = toc
	}
	newId := ""
	if entity.CategoryID != "" {
		f := zorm.NewSelectFinder(tableCategoryName, "name as categoryName").Append(" where id =?", entity.CategoryID)
		zorm.QueryRow(ctx, f, entity)
		urls := strings.Split(entity.Id, "/")
		newId = entity.CategoryID + urls[len(urls)-1]
	}
	entity.UpdateTime = now
	zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		funcUpdate(ctx, c, entity, entity.Id)
		if newId != "" {
			finder := zorm.NewUpdateFinder(tableContentName).Append("id=? WHERE id=?", newId, entity.Id)
			return zorm.UpdateFinder(ctx, finder)
		}
		return nil, nil
	})

}

// funcUpdateInit 初始化更新的对象参数,先从数据库查询,再更新数据
func funcUpdateInit(ctx context.Context, c *app.RequestContext, entity zorm.IEntityStruct) bool {
	jsontmp := make(map[string]interface{}, 0)
	c.Bind(&jsontmp)
	id := jsontmp["id"]
	finder := zorm.NewSelectFinder(entity.GetTableName()).Append("WHERE id=?", id)
	has, err := zorm.QueryRow(ctx, finder, entity)
	if !has || err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "Id不存在"})
		c.Abort() // 终止后续调用
		FuncLogError(ctx, err)
		return false
	}
	err = c.Bind(entity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "转换json数据错误"})
		c.Abort() // 终止后续调用
		FuncLogError(ctx, err)
		return false
	}
	return true
}

// funcUpdate 更新表数据
func funcUpdate(ctx context.Context, c *app.RequestContext, entity zorm.IEntityStruct, id string) {
	if id == "" { //没有id,终止调用
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "id不能为空"})
		c.Abort() // 终止后续调用
		return
	}
	_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		_, err := zorm.Update(ctx, entity)
		return nil, err
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "更新数据失败"})
		c.Abort() // 终止后续调用
		FuncLogError(ctx, err)
		return
	}
	c.JSON(http.StatusOK, ResponseData{StatusCode: 1})
}

// funcSavePre 跳转到保存页面
func funcSavePre(ctx context.Context, c *app.RequestContext) {
	urlPathParam := c.Param("urlPathParam")
	templateFile := "admin/" + urlPathParam + "/save.html"
	responseData := ResponseData{UrlPathParam: urlPathParam}
	responseData.QueryStringMap = wrapQueryStringMap(c)
	cHtmlAdmin(c, http.StatusOK, templateFile, responseData)
}

// funcSaveCategory 保存导航菜单
func funcSaveCategory(ctx context.Context, c *app.RequestContext) {
	entity := &Category{}
	err := c.Bind(entity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "转换json数据错误"})
		c.Abort() // 终止后续调用
		FuncLogError(ctx, err)
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	if entity.CreateTime == "" {
		entity.CreateTime = now
	}
	if entity.UpdateTime == "" {
		entity.UpdateTime = now
	}
	if entity.Pid != "" {
		entity.Id = entity.Pid + entity.Id + "/"
	} else {
		entity.Id = "/" + entity.Id + "/"
	}
	has := validateIDExists(ctx, entity.Id)
	if has {
		c.JSON(http.StatusConflict, ResponseData{StatusCode: 0, Message: "URL路径重复,请修改路径标识"})
		c.Abort() // 终止后续调用
		return
	}
	count, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		return zorm.Insert(ctx, entity)
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "保存数据失败"})
		c.Abort() // 终止后续调用
		FuncLogError(ctx, err)
		return
	}

	//增加路由映射
	addCategoryRoute(entity.Id)
	// 增加自定义路由映射
	//routeCategoryMap[funcTrimSuffixSlash(entity.Id)] = entity.Id

	c.JSON(http.StatusOK, ResponseData{StatusCode: count.(int), Message: "保存成功!"})
}

// funcSaveContent 保存内容
func funcSaveContent(ctx context.Context, c *app.RequestContext) {
	entity := &Content{}
	err := c.Bind(entity)
	if err != nil || entity.Id == "" || entity.CategoryID == "" {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "转换json数据错误"})
		c.Abort() // 终止后续调用
		FuncLogError(ctx, err)
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	// 构建ID
	entity.Id = entity.CategoryID + entity.Id
	has := validateIDExists(ctx, entity.Id)
	if has {
		c.JSON(http.StatusConflict, ResponseData{StatusCode: 0, Message: "URL路径重复,请修改路径标识"})
		c.Abort() // 终止后续调用
		return
	}
	if entity.CreateTime == "" {
		entity.CreateTime = now
	}
	if entity.UpdateTime == "" {
		entity.UpdateTime = now
	}
	if entity.Markdown != "" {
		content, toc, err := renderMarkdownHtml(entity.Markdown)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "markdown转html错误"})
			c.Abort() // 终止后续调用
			FuncLogError(ctx, err)
			return
		}
		entity.Content = content
		entity.Toc = toc
	}

	f := zorm.NewSelectFinder(tableCategoryName, "name as categoryName").Append(" where id =?", entity.CategoryID)
	zorm.QueryRow(ctx, f, entity)

	count, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		return zorm.Insert(ctx, entity)
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "保存数据失败"})
		c.Abort() // 终止后续调用
		FuncLogError(ctx, err)
		return
	}
	c.JSON(http.StatusOK, ResponseData{StatusCode: count.(int), Message: "保存成功!"})
}

// funcDelete 删除数据
func funcDelete(ctx context.Context, c *app.RequestContext) {
	id := c.PostForm("id")
	//id := c.Query("id")
	if id == "" { //没有id,终止调用
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "id不能为空"})
		c.Abort() // 终止后续调用
		return
	}
	urlPathParam := c.Param("urlPathParam")
	if urlPathParam == "category" {
		finder := zorm.NewSelectFinder(tableCategoryName, "*").Append(" where pid =?", id)
		page := zorm.NewPage()
		pageNo, _ := strconv.Atoi("1")
		page.PageNo = pageNo
		data := make([]Category, 0)
		zorm.Query(context.Background(), finder, &data, page)
		if len(data) != 0 {
			c.JSON(http.StatusOK, ResponseData{StatusCode: 0, Message: "无法删除有子级的导航!"})
		} else {
			err := deleteById(ctx, urlPathParam, id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "删除数据失败"})
				c.Abort() // 终止后续调用
			}
			c.JSON(http.StatusOK, ResponseData{StatusCode: 1, Message: "删除数据成功"})
		}
	} else {
		err := deleteById(ctx, urlPathParam, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "删除数据失败"})
			c.Abort() // 终止后续调用
		}
		c.JSON(http.StatusOK, ResponseData{StatusCode: 1, Message: "删除数据成功"})
	}
}

// permissionHandler 权限拦截器
func permissionHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		jwttoken := string(c.Cookie(config.JwttokenKey))
		//fmt.Println(config.JwtSecret)
		userId, err := userIdByToken(jwttoken)
		if err != nil || userId == "" {
			c.Redirect(http.StatusOK, cRedirecURI("admin/login"))
			c.Abort() // 终止后续调用
			return
		}
		c.SetCookie(config.JwttokenKey, jwttoken, config.Timeout, "/", "", protocol.CookieSameSiteStrictMode, false, true)
		// 传递从jwttoken获取的userId
		c.Set(tokenUserId, userId)
		// 设置用户类型是 管理员
		c.Set(userTypeKey, 1)
	}
}

// funcLookById 根据Id,跳转到查看页面
func funcLookById(ctx context.Context, c *app.RequestContext, templateFile string) {
	id := c.Query("id")
	if id == "" {
		c.Redirect(http.StatusOK, cRedirecURI("admin/error"))
		c.Abort() // 终止后续调用
		return
	}
	urlPathParam := c.Param("urlPathParam")
	if !alphaNumericReg.MatchString(urlPathParam) {
		c.Redirect(http.StatusOK, cRedirecURI("admin/error"))
		c.Abort() // 终止后续调用
		return
	}
	responseData := ResponseData{StatusCode: 0}
	data, err := funcSelectOne(urlPathParam, "* FROM "+urlPathParam+" WHERE id=? ", id)
	responseData.Data = data
	responseData.UrlPathParam = urlPathParam

	if err != nil {
		c.Redirect(http.StatusOK, cRedirecURI("admin/error"))
		c.Abort() // 终止后续调用
		return
	}
	lookFile := "admin/" + urlPathParam + "/" + templateFile
	responseData.StatusCode = 1
	responseData.QueryStringMap = wrapQueryStringMap(c)
	cHtmlAdmin(c, http.StatusOK, lookFile, responseData)
}

// renderMarkdownHtml 渲染markdown内容,返回html,toc和error
func renderMarkdownHtml(mkstring string) (string, string, error) {
	if mkstring != "" {
		_, tocHtml, html, err := conver2Html([]byte(mkstring))
		if err != nil {
			return "", "", err
		}
		return *html, *tocHtml, nil
	}
	return "", "", nil
}

// wrapQueryStringMap 包装查询参数Map
func wrapQueryStringMap(c *app.RequestContext) map[string]string {
	queryStringMap := make(map[string]string, 0)
	c.BindQuery(&queryStringMap)
	for k := range queryStringMap {
		queryStringMap[k] = c.Query(k)
	}
	return queryStringMap
}
