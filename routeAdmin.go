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
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/cloudwego/hertz/pkg/route/param"
	"golang.org/x/crypto/sha3"
)

// adminGroup路由组,使用变量声明,优先级高于init函数
var adminGroup = initAdminGroup()

func initAdminGroup() *route.RouterGroup {
	// 设置日志级别
	hlog.SetLevel(hlog.LevelError)
	//设置json处理函数
	//binding.ResetJSONUnmarshaler(json.Unmarshal)
	/*
		binding.Default().ResetJSONUnmarshaler(func(data []byte, v interface{}) error {
			dec := json.NewDecoder(bytes.NewBuffer(data))
			dec.UseNumber()
			return dec.Decode(v)
		})
	*/
	// 初始化模板
	err := initTemplate()
	if err != nil { // 初始化模板异常
		panic("初始化模板异常")
	}
	admin := h.Group("/admin")
	admin.Use(permissionHandler())
	return admin
}

// 初始化函数
func init() {

	// 异常页面
	h.GET("/admin/error", func(ctx context.Context, c *app.RequestContext) {
		c.HTML(http.StatusOK, "admin/error.html", nil)
	})

	// 安装
	h.GET("/admin/install", func(ctx context.Context, c *app.RequestContext) {
		if installed { // 如果已经安装过了,跳转到登录
			c.Redirect(http.StatusOK, cRedirecURI("admin/login"))
			c.Abort() // 终止后续调用
			return
		}
		c.HTML(http.StatusOK, "admin/install.html", nil)
	})
	h.POST("/admin/install", func(ctx context.Context, c *app.RequestContext) {
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
	})

	// 生成30位随机数,钱包签名随机校验.如果32位metamask会解析成16进制字符串,可能是metamask的bug
	h.POST("/admin/random", func(ctx context.Context, c *app.RequestContext) {
		generateChainRandStr()
		//返回到前端
		c.JSON(http.StatusOK, ResponseData{StatusCode: 1, Data: chainRandStr})
	})

	// 后台管理员登录
	h.GET("/admin/login", func(ctx context.Context, c *app.RequestContext) {
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
		c.HTML(http.StatusOK, "admin/login.html", responseData)
	})

	h.POST("/admin/login", func(ctx context.Context, c *app.RequestContext) {
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
		/*
			password := c.PostForm("password")
			bytehex := sha3.Sum512([]byte("admin"))
			str := hex.EncodeToString(bytehex[:])
			if password == str {
				fmt.Println(password)
			}
		*/
		jwttoken, _ := newJWTToken(userId)

		// c.HTML(http.StatusOK, "admin/index.html", nil)
		c.SetCookie(config.JwttokenKey, jwttoken, config.Timeout, "/", "", protocol.CookieSameSiteStrictMode, false, true)
		errorLoginCount.Store(0)
		c.Redirect(http.StatusOK, cRedirecURI("admin/index"))
	})

	// 后台管理员使用区块链账号登录
	h.GET("/admin/chainlogin", func(ctx context.Context, c *app.RequestContext) {
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
		c.HTML(http.StatusOK, "admin/chainlogin.html", responseData)
	})

	h.POST("/admin/chainlogin", func(ctx context.Context, c *app.RequestContext) {
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
	})

	// 后台管理员首页
	adminGroup.GET("/index", func(ctx context.Context, c *app.RequestContext) {
		// 获取从jwttoken中解码的userId
		userId, ok := c.Get(tokenUserId)
		if !ok || userId == "" {
			c.Redirect(http.StatusOK, cRedirecURI("admin/login"))
			c.Abort() // 终止后续调用
			return
		}

		c.HTML(http.StatusOK, "admin/index.html", nil)
	})
	// 刷新站点,重新加载资源包含模板和对应的静态文件
	adminGroup.GET("/reload", func(ctx context.Context, c *app.RequestContext) {
		err := loadTemplate()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, ERR: err})
			c.Abort() // 终止后续调用
			return
		}
		//重新生成静态文件
		go genStaticFile()

		//此处为hertz bug,已经调用了 h.SetHTMLTemplate(tmpl),但是c.HTMLRender依然是老的内存地址
		//c.HTMLRender = render.HTMLProduction{Template: tmpl}
		//c.HTML(http.StatusOK, "admin/index.html", nil)
		c.JSON(http.StatusOK, ResponseData{StatusCode: 1})
	})

	//上传文件
	adminGroup.POST("/upload", func(ctx context.Context, c *app.RequestContext) {
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
		c.JSON(http.StatusOK, ResponseData{StatusCode: 1, Data: funcBasePath() + path})
	})

	//上传主题文件
	adminGroup.POST("/themeTemplate/uploadTheme", func(ctx context.Context, c *app.RequestContext) {
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
		//解压压缩包
		err = unzip(path, themeDir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, ERR: err})
			c.Abort() // 终止后续调用
			return
		}
		os.Remove(path)
		c.JSON(http.StatusOK, ResponseData{StatusCode: 1, Data: funcBasePath() + path})
	})

	// 内容预览
	adminGroup.GET("/content/look", funcContentPreview)
	// 栏目预览
	adminGroup.GET("/category/look", funcCategoryPreview)

	// 查询Content列表,根据CategoryId like
	adminGroup.GET("/content/list", funcContentList)

	// 查询主题模板
	adminGroup.GET("/themeTemplate/list", funcThemeTemplateList)

	// 修改主题模板文件
	adminGroup.POST("/themeTemplate/update", funcThemeTemplatePost)

	// 通用list列表,先都使用get方法
	adminGroup.GET("/:urlPathParam/list", funcList)
	//adminGroup.POST("/:urlPathParam/list", funcList)

	// 通用查看
	adminGroup.GET("/:urlPathParam/look", funcLook)

	//跳转到修改页面
	adminGroup.GET("/:urlPathParam/update", funcUpdatePre)
	//ajax POST提交JSON信息,返回方法JSON
	adminGroup.POST("/:urlPathParam/update", funcUpdate)
	//跳转到保存页面
	adminGroup.GET("/:urlPathParam/save", funcSavePre)
	//ajax POST提交JSON信息,返回方法JSON
	adminGroup.POST("/:urlPathParam/save", funcSave)

	//ajax POST提交JSON信息,返回方法JSON
	adminGroup.POST("/:urlPathParam/delete", funcDelete)

}

// alphaNumericReg 传入的列名只能是字母数字或下划线,长度不超过20
var alphaNumericReg = regexp.MustCompile("^[a-zA-Z0-9_]{1,20}$")

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
	if err != nil { //表不存在
		c.Redirect(http.StatusOK, cRedirecURI("admin/error"))
		c.Abort() // 终止后续调用
		return
	}

	//优先使用自定义模板文件
	listFile := "admin/" + urlPathParam + "/list.html"
	/*
		t := tmpl.Lookup(listFile)
		if t == nil { //不存在自定义模板,使用通用模板
			listFile = "admin/list.html"
		}
	*/
	//queryString := c.Request.QueryString()
	//responseData.QueryString = string(queryString)
	//responseData.UrlPathParam = urlPathParam
	c.HTML(http.StatusOK, listFile, responseData)
}

// funcLook 通用查看,根据id查看
func funcLook(ctx context.Context, c *app.RequestContext) {
	funcTableById(ctx, c, "look.html")
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
	//funcTableById(ctx, c, "look.html")
}

// funcCategoryPreview 栏目预览
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
	//funcTableById(ctx, c, "look.html")
}

// funcContentList 查询Content列表,根据CategoryId like
func funcContentList(ctx context.Context, c *app.RequestContext) {
	urlPathParam := "content"
	//获取页码
	pageNoStr := c.DefaultQuery("pageNo", "1")
	q := strings.TrimSpace(c.Query("q"))
	pageNo, _ := strconv.Atoi(pageNoStr)
	comCode := strings.TrimSpace(c.Query("comCode"))
	values := make([]interface{}, 0)
	sql := ""
	if comCode != "" {
		sql = " * from content where comCode like ?  order by sortNo desc "
		values = append(values, comCode+"%")
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
	if err != nil { //表不存在
		c.Redirect(http.StatusOK, cRedirecURI("admin/error"))
		c.Abort() // 终止后续调用
		return
	}
	//优先使用自定义模板文件
	listFile := "admin/" + urlPathParam + "/list.html"
	c.HTML(http.StatusOK, listFile, responseData)
}

// funcThemeTemplateList 所有的主题文件列表
func funcThemeTemplateList(ctx context.Context, c *app.RequestContext) {
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
	//优先使用自定义模板文件
	listFile := "admin/" + urlPathParam + "/list.html"

	filePath := c.Query("file")
	if filePath == "" || strings.Contains(filePath, "..") {
		c.HTML(http.StatusOK, listFile, responseData)
		return
	}
	filePath = filepath.ToSlash(filePath)
	fileContent, err := os.ReadFile(themeDir + filePath)
	if err != nil {
		responseData.ERR = err
		c.HTML(http.StatusOK, listFile, responseData)
		return
	}
	responseData.ExtMap["file"] = string(fileContent)
	c.HTML(http.StatusOK, listFile, responseData)
}

func funcThemeTemplatePost(ctx context.Context, c *app.RequestContext) {
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

// funcUpdatePre 修改页面
func funcUpdatePre(ctx context.Context, c *app.RequestContext) {
	funcTableById(ctx, c, "update.html")
}

// 修改内容
func funcUpdate(ctx context.Context, c *app.RequestContext) {
	urlPathParam := c.Param("urlPathParam")
	funcUpdateTable(ctx, c, urlPathParam)
}

// 修改内容
func funcUpdateTable(ctx context.Context, c *app.RequestContext, urlPathParam string) {
	var entity zorm.IEntityStruct
	var err error
	var now = time.Now().Format("2006-01-02 15:04:05")
	id := ""
	mastUpdateColumn := []string{"status"}
	switch urlPathParam {
	case tableConfigName:
		ptrObj := &Config{}
		err = c.Bind(ptrObj)
		id = ptrObj.Id
		ptrObj.UpdateTime = now
		entity = ptrObj
	case tableUserName:
		ptrObj := &User{}
		err = c.Bind(ptrObj)
		id = ptrObj.Id
		ptrObj.UpdateTime = now
		if ptrObj.Password != "" {
			// 重新hash密码,避免拖库后撞库
			sha3Bytes := sha3.Sum512([]byte(ptrObj.Password))
			ptrObj.Password = hex.EncodeToString(sha3Bytes[:])
		}
		entity = ptrObj
	case tableSiteName:
		ptrObj := &Site{}
		err = c.Bind(ptrObj)
		id = ptrObj.Id
		ptrObj.UpdateTime = now
		entity = ptrObj
	case tableCategoryName:
		ptrObj := &Category{}
		err = c.Bind(ptrObj)
		id = ptrObj.Id
		ptrObj.UpdateTime = now
		if ptrObj.Pid != "" {
			f := zorm.NewSelectFinder(urlPathParam, "comCode").Append(" where id =?", ptrObj.Pid)
			zorm.QueryRow(ctx, f, &(ptrObj.ComCode))
			ptrObj.ComCode = ptrObj.ComCode + ptrObj.Id + ","
		} else {
			ptrObj.ComCode = "," + ptrObj.Id + ","
		}
		mastUpdateColumn = append(mastUpdateColumn, "hrefTarget")
		entity = ptrObj
	case tableContentName:
		ptrObj := &Content{}
		err = c.Bind(ptrObj)
		id = ptrObj.Id
		ptrObj.UpdateTime = now
		if ptrObj.Markdown != "" {
			content, toc, err := setMarkdownHtml(ptrObj.Markdown)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "markdown转html错误"})
				c.Abort() // 终止后续调用
				FuncLogError(err)
				return
			}
			ptrObj.Content = content
			ptrObj.Toc = toc
		}
		if ptrObj.CategoryID != "" {
			f := zorm.NewSelectFinder(tableCategoryName, "comCode,name as categoryName").Append(" where id =?", ptrObj.CategoryID)
			zorm.QueryRow(ctx, f, ptrObj)
		}
		mastUpdateColumn = append(mastUpdateColumn, "markdown", "content")
		entity = ptrObj
	default:
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "表不存在!"})
		c.Abort() // 终止后续调用
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "转换json数据错误"})
		c.Abort() // 终止后续调用
		FuncLogError(err)
		return
	}
	if id == "" { //没有id,终止调用
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "id不能为空"})
		c.Abort() // 终止后续调用
		return
	}

	_, err = zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		ctx, err = zorm.BindContextMustUpdateCols(ctx, mastUpdateColumn)
		_, err = zorm.UpdateNotZeroValue(ctx, entity)
		return nil, err
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "更新数据失败"})
		c.Abort() // 终止后续调用
		FuncLogError(err)
		return
	}
	if urlPathParam == tableContentName {
		genSearchDataJson()
	}
	c.JSON(http.StatusOK, ResponseData{StatusCode: 1})
}

// funcSavePre 保存页面
func funcSavePre(ctx context.Context, c *app.RequestContext) {
	urlPathParam := c.Param("urlPathParam")

	//优先使用自定义模板文件
	templateFile := "admin/" + urlPathParam + "/save.html"
	/*
		t := tmpl.Lookup(templateFile)
		if t == nil { //不存在自定义模板,使用通用模板
			templateFile = "admin/save.html"
		}
	*/
	responseData := ResponseData{UrlPathParam: urlPathParam}

	responseData.QueryStringMap = wrapQueryStringMap(c)

	c.HTML(http.StatusOK, templateFile, responseData)
}

// 保存内容
func funcSave(ctx context.Context, c *app.RequestContext) {
	urlPathParam := c.Param("urlPathParam")
	funcSaveTable(ctx, c, urlPathParam)
}

// 保存内容
func funcSaveTable(ctx context.Context, c *app.RequestContext, urlPathParam string) {
	var entity zorm.IEntityStruct
	var err error
	var now = time.Now().Format("2006-01-02 15:04:05")
	switch urlPathParam {
	case tableConfigName:
		ptrObj := &Config{}
		err = c.Bind(ptrObj)
		if ptrObj.Id == "" {
			ptrObj.Id = FuncGenerateStringID()
		}
		if ptrObj.CreateTime == "" {
			ptrObj.CreateTime = now
		}
		if ptrObj.UpdateTime == "" {
			ptrObj.UpdateTime = now
		}
		entity = ptrObj
	case tableUserName:
		ptrObj := &User{}
		err = c.Bind(ptrObj)
		if ptrObj.Id == "" {
			ptrObj.Id = FuncGenerateStringID()
		}
		if ptrObj.CreateTime == "" {
			ptrObj.CreateTime = now
		}
		if ptrObj.UpdateTime == "" {
			ptrObj.UpdateTime = now
		}
		entity = ptrObj
	case tableSiteName:
		ptrObj := &Site{}
		err = c.Bind(ptrObj)
		if ptrObj.Id == "" {
			ptrObj.Id = FuncGenerateStringID()
		}
		if ptrObj.CreateTime == "" {
			ptrObj.CreateTime = now
		}
		if ptrObj.UpdateTime == "" {
			ptrObj.UpdateTime = now
		}
		entity = ptrObj
	case tableCategoryName:
		ptrObj := &Category{}
		err = c.Bind(ptrObj)
		if ptrObj.Id == "" {
			ptrObj.Id = FuncGenerateStringID()
		}
		if ptrObj.CreateTime == "" {
			ptrObj.CreateTime = now
		}
		if ptrObj.UpdateTime == "" {
			ptrObj.UpdateTime = now
		}
		if ptrObj.Pid != "" {
			f := zorm.NewSelectFinder(urlPathParam, "comCode").Append(" where id =?", ptrObj.Pid)
			zorm.QueryRow(ctx, f, &(ptrObj.ComCode))
			ptrObj.ComCode = ptrObj.ComCode + ptrObj.Id + ","
		} else {
			ptrObj.ComCode = "," + ptrObj.Id + ","
		}
		entity = ptrObj
	case tableContentName:
		ptrObj := &Content{}
		err = c.Bind(ptrObj)
		if ptrObj.Id == "" {
			ptrObj.Id = FuncGenerateStringID()
		}
		if ptrObj.CreateTime == "" {
			ptrObj.CreateTime = now
		}
		if ptrObj.UpdateTime == "" {
			ptrObj.UpdateTime = now
		}

		if ptrObj.Markdown != "" {
			content, toc, err := setMarkdownHtml(ptrObj.Markdown)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "markdown转html错误"})
				c.Abort() // 终止后续调用
				FuncLogError(err)
				return
			}
			ptrObj.Content = content
			ptrObj.Toc = toc
		}
		if ptrObj.CategoryID != "" {
			f := zorm.NewSelectFinder(tableCategoryName, "comCode,name as categoryName").Append(" where id =?", ptrObj.CategoryID)
			zorm.QueryRow(ctx, f, ptrObj)
		}
		entity = ptrObj
	default:
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "表不存在!"})
		c.Abort() // 终止后续调用
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "转换json数据错误"})
		c.Abort() // 终止后续调用
		FuncLogError(err)
		return
	}

	count, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
		return zorm.Insert(ctx, entity)
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "保存数据失败"})
		c.Abort() // 终止后续调用
		FuncLogError(err)
		return
	}
	if urlPathParam == tableContentName {
		genSearchDataJson()
	}
	c.JSON(http.StatusOK, ResponseData{StatusCode: count.(int), Message: "保存成功!"})
}

// 删除内容
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
			//tableName := bleveDataDir + urlPathParam
			err := deleteById(ctx, urlPathParam, id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "删除数据失败"})
				c.Abort() // 终止后续调用
			}
			c.JSON(http.StatusOK, ResponseData{StatusCode: 1, Message: "删除数据成功"})
		}
	} else {
		//tableName := bleveDataDir + urlPathParam
		err := deleteById(ctx, urlPathParam, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "删除数据失败"})
			c.Abort() // 终止后续调用
		}
		if urlPathParam == tableContentName {
			genSearchDataJson()
		}
		c.JSON(http.StatusOK, ResponseData{StatusCode: 1, Message: "删除数据成功"})
	}
}

/*
// permissionHandler 权限拦截器
func permissionHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		rs := &requestState{c: c}
		ctx = context.WithValue(ctx, requestStateKey{}, rs)
		//c.Request.AppendBodyString("go test body")
		fmt.Println("------------开始调用wasm---- ")

		ctxNext, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {
			t1 := time.Now()
			outCtx, next, _, err := wasmFediOS.HandleHTTP(ctx, []byte("test extData 1"))
			c.Request.Header.Set(api.WasmCallMethodHeaderName, api.HeaderHandleFollowNodeFn)
			_, next, value, err := wasmFediOS.HandleHTTP(outCtx, []byte("test extData 2"))
			tm2 := time.Now().Sub(t1).Milliseconds()
			fmt.Println("------------结束调用wasm----,耗时:" + strconv.FormatInt(tm2, 10) + "毫秒,返回值:" + string(value))
			return next, err
		})
		next, ok := ctxNext.(bool)
		if err != nil || !ok || !next { //中止调用
			c.Abort()
			return
		}

	}
}
*/

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

// funcTableById 根据Id查询表信息
func funcTableById(ctx context.Context, c *app.RequestContext, htmlfile string) {
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
	//tableName := bleveDataDir + urlPathParam
	data, err := funcSelectOne(urlPathParam, "* FROM "+urlPathParam+" WHERE id=? ", id)
	responseData.Data = data
	//responseData["UrlPathParam"] = urlPathParam
	responseData.UrlPathParam = urlPathParam

	if err != nil { //表不存在
		c.Redirect(http.StatusOK, cRedirecURI("admin/error"))
		c.Abort() // 终止后续调用
		return
	}
	//优先使用自定义模板文件
	lookFile := "admin/" + urlPathParam + "/" + htmlfile
	/*
		t := tmpl.Lookup(lookFile)
		if t == nil { //不存在自定义模板,使用通用模板
			lookFile = "admin/" + htmlfile
		}
	*/
	responseData.StatusCode = 1
	responseData.QueryStringMap = wrapQueryStringMap(c)
	c.HTML(http.StatusOK, lookFile, responseData)
}

func setMarkdownHtml(mkstring string) (string, string, error) {
	if mkstring != "" {
		_, tocHtml, html, err := conver2Html([]byte(mkstring))
		if err != nil {
			return "", "", err
		}
		return *html, *tocHtml, nil
	}
	return "", "", nil
}

func wrapQueryStringMap(c *app.RequestContext) map[string]string {
	queryStringMap := make(map[string]string, 0)
	c.BindQuery(&queryStringMap)
	for k := range queryStringMap {
		queryStringMap[k] = c.Query(k)
	}
	return queryStringMap
}
