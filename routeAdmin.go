package main

import (
	"context"
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	//"github.com/bytedance/go-tagexpr/v2/binding"
	"gitee.com/chunanyong/zorm"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/route"
)

// adminGroup路由组,使用变量声明,优先级高于init函数
var adminGroup = initAdminGroup()

func initAdminGroup() *route.RouterGroup {
	// 设置日志级别
	hlog.SetLevel(hlog.LevelInfo)
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
			return
		}
		c.HTML(http.StatusOK, "admin/install.html", nil)
	})
	h.POST("/admin/install", func(ctx context.Context, c *app.RequestContext) {
		if installed { // 如果已经安装过了,跳转到登录
			c.Redirect(http.StatusOK, cRedirecURI("admin/login"))
			return
		}
		// 使用后端管理界面配置,jwtSecret也有后端随机产生
		account := c.PostForm("account")
		password := c.PostForm("password")
		err := insertUser(ctx, account, password)
		if err != nil {
			c.Redirect(http.StatusOK, cRedirecURI("admin/error"))
			return
		}
		// 安装成功,更新安装状态
		updateInstall(ctx)
		c.Redirect(http.StatusOK, cRedirecURI("admin/login"))
	})

	// 后台管理员登录
	h.GET("/admin/login", func(ctx context.Context, c *app.RequestContext) {
		if !installed { // 如果没有安装,跳转到安装
			c.Redirect(http.StatusOK, cRedirecURI("admin/install"))
			return
		}
		c.SetCookie(config.JwttokenKey, "", config.Timeout, "/", "", protocol.CookieSameSiteStrictMode, true, true)
		c.HTML(http.StatusOK, "admin/login.html", nil)
	})
	h.POST("/admin/login", func(ctx context.Context, c *app.RequestContext) {
		if !installed { // 如果没有安装,跳转到安装
			c.Redirect(http.StatusOK, cRedirecURI("admin/install"))
			return
		}
		account := c.PostForm("account")
		password := c.PostForm("password")
		userId, err := findUserId(ctx, account, password)
		if userId == "" || err != nil { // 用户不存在或者异常
			c.Redirect(http.StatusOK, cRedirecURI("admin/login"))
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
		jwttoken, _ := newJWTToken(userId, nil)

		// c.HTML(http.StatusOK, "admin/index.html", nil)
		c.SetCookie(config.JwttokenKey, jwttoken, config.Timeout, "/", "", protocol.CookieSameSiteStrictMode, true, true)

		c.Redirect(http.StatusOK, cRedirecURI("admin/index"))
	})

	// 后台管理员首页
	adminGroup.GET("/index", func(ctx context.Context, c *app.RequestContext) {
		// 获取从jwttoken中解码的userId
		userId, ok := c.Get(tokenUserId)
		if !ok || userId == "" {
			c.Redirect(http.StatusOK, cRedirecURI("admin/login"))
			return
		}

		c.HTML(http.StatusOK, "admin/index.html", nil)
	})
	// 重新加载资源包含模板和对应的静态文件
	adminGroup.GET("/reload", func(ctx context.Context, c *app.RequestContext) {
		err := loadTemplate()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, ERR: err})
			return
		}
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
			return
		}
		path := "public/" + zorm.FuncGenerateStringID(ctx) + filepath.Ext(fileHeader.Filename)
		newFileName := datadir + path
		err = c.SaveUploadedFile(fileHeader, newFileName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, ERR: err})
			return
		}
		c.JSON(http.StatusOK, ResponseData{StatusCode: 1, Data: funcBasePath() + path})
	})

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

	//ajax POST提交新增表信息
	adminGroup.POST("/tableInfo/save", funcTableInfoSave)

	//ajax POST提交新增字段信息
	adminGroup.POST("/tableField/save", funcTableFieldSave)

	//ajax POST提交JSON信息,返回方法JSON
	adminGroup.POST("/:urlPathParam/delete", funcDelete)

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
	var params strings.Builder
	i := 0
	for k := range mapParams {
		if i > 0 {
			params.WriteString(" and ")
		}
		params.WriteString(k)
		params.WriteByte('=')
		params.WriteString(c.Query(k))
		i++
	}
	where := params.String()
	sql := "* from " + urlPathParam
	if where != "" {
		sql += " WHERE " + where
	}
	sql += " order by sortNo desc "
	responseData, err := funcSelectList(q, pageNo, sql)
	responseData["urlPathParam"] = urlPathParam
	if err != nil { //表不存在
		c.Redirect(http.StatusOK, cRedirecURI("admin/error"))
		c.Abort() // 终止后续调用
		return
	}

	//优先使用自定义模板文件
	listFile := "admin/" + urlPathParam + "/list.html"
	t := tmpl.Lookup(listFile)
	if t == nil { //不存在自定义模板,使用通用模板
		listFile = "admin/list.html"
	}
	//queryString := c.Request.QueryString()
	//responseData.QueryString = string(queryString)
	//responseData.UrlPathParam = urlPathParam
	c.HTML(http.StatusOK, listFile, responseData)
}

// funcLook 通用查看,根据id查看
func funcLook(ctx context.Context, c *app.RequestContext) {
	funcTableById(ctx, c, "look.html")
}

// funcUpdatePre 修改页面
func funcUpdatePre(ctx context.Context, c *app.RequestContext) {
	funcTableById(ctx, c, "update.html")
}

// 修改内容
func funcUpdate(ctx context.Context, c *app.RequestContext) {

	newMap := make(map[string]interface{}, 0)
	err := c.Bind(&newMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "转换json数据错误"})
		c.Abort() // 终止后续调用
		FuncLogError(err)
		return
	}

	id := ""
	if newMap["id"] != nil {
		id = newMap["id"].(string)
	}
	if id == "" { //没有id,认为是新增
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "id不能为空"})
		c.Abort() // 终止后续调用
		return
	}
	urlPathParam := c.Param("urlPathParam")
	//tableName := bleveDataDir + urlPathParam

	if !tableExist(urlPathParam) {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "数据不存在"})
		c.Abort() // 终止后续调用
		return
	}
	err = setMarkdownHtml(&newMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "markdown转html错误"})
		c.Abort() // 终止后续调用
		FuncLogError(err)
		return
	}

	entityMap := zorm.NewEntityMap(urlPathParam)
	for k, v := range newMap {
		entityMap.Set(k, v)
	}
	entityMap.PkColumnName = "id"
	entityMap.Set("updateTime", time.Now().Format("2006-01-02 15:04:05"))
	err = updateTable(ctx, entityMap)
	if err != nil { //没有id,认为是新增
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "更新数据失败"})
		c.Abort() // 终止后续调用
		FuncLogError(err)
		return
	}
	c.JSON(http.StatusOK, ResponseData{StatusCode: 1})
}

// funcSavePre 保存页面
func funcSavePre(ctx context.Context, c *app.RequestContext) {
	urlPathParam := c.Param("urlPathParam")
	//优先使用自定义模板文件
	updateFile := "admin/" + urlPathParam + "/save.html"
	t := tmpl.Lookup(updateFile)
	if t == nil { //不存在自定义模板,使用通用模板
		updateFile = "admin/save.html"
	}
	c.HTML(http.StatusOK, updateFile, responData2Map(ResponseData{UrlPathParam: urlPathParam}))
}

// 保存内容
func funcSave(ctx context.Context, c *app.RequestContext) {
	urlPathParam := c.Param("urlPathParam")
	//tableName := bleveDataDir + urlPathParam
	if !tableExist(urlPathParam) {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "数据不存在"})
		c.Abort() // 终止后续调用
		return
	}

	newMap := make(map[string]interface{}, 0)

	//err = json.Unmarshal(jsonBody, &newMap)
	err := c.Bind(&newMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "转换json数据错误"})
		c.Abort() // 终止后续调用
		FuncLogError(err)
		return
	}
	err = setMarkdownHtml(&newMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "markdown转html错误"})
		c.Abort() // 终止后续调用
		FuncLogError(err)
		return
	}
	entityMap := zorm.NewEntityMap(urlPathParam)
	for k, v := range newMap {
		entityMap.Set(k, v)
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	entityMap.Set("createTime", now)
	entityMap.Set("updateTime", now)
	responseData, err := saveEntityMap(ctx, entityMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "保存数据失败"})
		c.Abort() // 终止后续调用
		FuncLogError(err)
		return
	}
	c.JSON(http.StatusOK, responData2Map(responseData))
}

// 修改内容
func funcDelete(ctx context.Context, c *app.RequestContext) {
	id := c.PostForm("id")
	//id := c.Query("id")
	if id == "" { //没有id,认为是新增
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "id不能为空"})
		c.Abort() // 终止后续调用
		return
	}
	urlPathParam := c.Param("urlPathParam")
	//tableName := bleveDataDir + urlPathParam
	err := deleteById(ctx, urlPathParam, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "删除数据失败"})
		c.Abort() // 终止后续调用
		FuncLogError(err)
	}
	c.JSON(http.StatusOK, ResponseData{StatusCode: 1, Message: "删除数据成功"})
}

// 保存表内容
func funcTableInfoSave(ctx context.Context, c *app.RequestContext) {
	newMap := make(map[string]interface{}, 0)
	err := c.Bind(&newMap)
	tableCode := newMap["code"]
	if err != nil || tableCode == nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "新增表失败"})
		c.Abort() // 终止后续调用
		FuncLogError(err)
	}
	tableCodeString := tableCode.(string)
	createTableSQL := `CREATE TABLE IF NOT EXISTS ` + tableCodeString + ` (
		id TEXT PRIMARY KEY     NOT NULL
	 ) strict ;`
	_, err = execNativeSQL(ctx, createTableSQL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "新增表失败"})
		c.Abort() // 终止后续调用
		FuncLogError(err)
	}
	entityMap := zorm.NewEntityMap(tableInfoName)
	for k, v := range newMap {
		entityMap.Set(k, v)
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	entityMap.Set("createTime", now)
	entityMap.Set("updateTime", now)
	responseData, err := saveEntityMap(ctx, entityMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "保存数据失败"})
		c.Abort() // 终止后续调用
		FuncLogError(err)
		return
	}
	c.JSON(http.StatusOK, responData2Map(responseData))
}

// 保存字段内容
func funcTableFieldSave(ctx context.Context, c *app.RequestContext) {
	fieldStruct := TableFieldStruct{}
	err := c.Bind(&fieldStruct)
	if err != nil || fieldStruct.TableCode == "" || fieldStruct.FieldCode == "" {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "新增字段失败"})
		c.Abort() // 终止后续调用
		FuncLogError(err)
	}

	sqlType := "text"
	if fieldStruct.FieldType == 1 { //数字
		sqlType = "int"
	}
	// code
	createTableSQL := "alter table " + fieldStruct.TableCode + "  add column " + fieldStruct.FieldCode + " " + sqlType
	_, err = execNativeSQL(ctx, createTableSQL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "新增表失败"})
		c.Abort() // 终止后续调用
		FuncLogError(err)
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	fieldStruct.CreateTime = now
	fieldStruct.UpdateTime = now
	entityMap := zorm.NewEntityMap(tableFieldName)

	newMap := make(map[string]interface{}, 0)
	jsonByte, _ := json.Marshal(fieldStruct)
	json.Unmarshal(jsonByte, &newMap)
	for k, v := range newMap {
		entityMap.Set(k, v)
	}
	responseData, err := saveEntityMap(ctx, entityMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "保存数据失败"})
		c.Abort() // 终止后续调用
		FuncLogError(err)
		return
	}
	c.JSON(http.StatusOK, responData2Map(responseData))
}

// permissionHandler 权限拦截器
func permissionHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		jwttoken := string(c.Cookie(config.JwttokenKey))
		userId, err := userIdByToken(jwttoken)
		if err != nil || userId == "" {
			c.Redirect(http.StatusOK, cRedirecURI("admin/login"))
			c.Abort() // 终止后续调用
			return
		}
		c.SetCookie(config.JwttokenKey, jwttoken, config.Timeout, "/", "", protocol.CookieSameSiteStrictMode, true, true)
		// 传递从jwttoken获取的userId
		c.Set(tokenUserId, userId)
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
	//tableName := bleveDataDir + urlPathParam
	responseData, err := funcSelectOne("* FROM "+urlPathParam+" WHERE id=? ", id)
	responseData["urlPathParam"] = urlPathParam
	if err != nil { //表不存在
		c.Redirect(http.StatusOK, cRedirecURI("admin/error"))
		c.Abort() // 终止后续调用
		return
	}
	//优先使用自定义模板文件
	lookFile := "admin/" + urlPathParam + "/" + htmlfile
	t := tmpl.Lookup(lookFile)
	if t == nil { //不存在自定义模板,使用通用模板
		lookFile = "admin/" + htmlfile
	}
	c.HTML(http.StatusOK, lookFile, responseData)
}

func setMarkdownHtml(newMap *map[string]interface{}) error {
	markdown, ok := (*newMap)["markdown"]
	// markdown转html
	if ok {
		mkstring := markdown.(string)
		if mkstring != "" {
			_, tocHtml, html, err := conver2Html([]byte(mkstring))
			if err != nil {
				return err
			}

			if html != nil && *html != "" {
				(*newMap)["content"] = *html
				(*newMap)["toc"] = *tocHtml
			}
		}

	}
	return nil
}
