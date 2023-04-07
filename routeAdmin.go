package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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

	h.GET("/getnav", func(ctx context.Context, c *app.RequestContext) {
		result, _ := getNavMenu("0")
		c.JSON(http.StatusOK, result)
	})

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
	// 后台管理员首页
	adminGroup.GET("/reload", func(ctx context.Context, c *app.RequestContext) {
		err := loadTemplate(true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, ERR: err})
			return
		}
		//此处为hertz bug,已经调用了 h.SetHTMLTemplate(tmpl),但是c.HTMLRender依然是老的内存地址
		//c.HTMLRender = render.HTMLProduction{Template: tmpl}
		//c.HTML(http.StatusOK, "admin/index.html", nil)
		c.JSON(http.StatusOK, ResponseData{StatusCode: 1})
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
			params.WriteByte('&')
		}
		params.WriteString(k)
		params.WriteByte('=')
		params.WriteString(c.Query(k))
		i++
	}
	responseData, err := funcIndexList(urlPathParam, "*", q, pageNo, params.String())
	if err != nil { //索引不存在
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
	funcIndexById(ctx, c, "look.html")
}

// funcUpdatePre 修改页面
func funcUpdatePre(ctx context.Context, c *app.RequestContext) {
	funcIndexById(ctx, c, "update.html")
}

// 修改内容
func funcUpdate(ctx context.Context, c *app.RequestContext) {
	jsonBody, err := c.Body()
	if err != nil { //没有id,认为是新增
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "获取body内容错误"})
		c.Abort() // 终止后续调用
		FuncLogError(err)
		return
	}
	newMap := make(map[string]interface{}, 0)
	err = json.Unmarshal(jsonBody, &newMap)
	if err != nil { //没有id,认为是新增
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
	//indexName := bleveDataDir + urlPathParam
	_, ok, _ := openBleveIndex(urlPathParam)
	if !ok { //索引不存在
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "数据不存在"})
		c.Abort() // 终止后续调用
		return
	}

	err = updateIndex(ctx, urlPathParam, id, newMap)
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
	//indexName := bleveDataDir + urlPathParam
	_, ok, _ := openBleveIndex(urlPathParam)
	if !ok { //索引不存在
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "数据不存在"})
		c.Abort() // 终止后续调用
		return
	}
	jsonBody, err := c.Body()
	if err != nil { //没有id,认为是新增
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "获取body内容错误"})
		c.Abort() // 终止后续调用
		FuncLogError(err)
		return
	}
	newMap := make(map[string]interface{}, 0)
	err = json.Unmarshal(jsonBody, &newMap)
	if err != nil { //没有id,认为是新增
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "转换json数据错误"})
		c.Abort() // 终止后续调用
		FuncLogError(err)
		return
	}
	responseData, err := saveNewIndex(ctx, urlPathParam, newMap)
	if err != nil { //没有id,认为是新增
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
	//indexName := bleveDataDir + urlPathParam
	err := deleteById(ctx, urlPathParam, id)
	if err != nil { //没有id,认为是新增
		c.JSON(http.StatusInternalServerError, ResponseData{StatusCode: 0, Message: "删除数据失败"})
		c.Abort() // 终止后续调用
		FuncLogError(err)
	}
	c.JSON(http.StatusOK, ResponseData{StatusCode: 1, Message: "删除数据成功"})
}

// permissionHandler 权限拦截器
func permissionHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		jwttoken := c.Cookie(config.JwttokenKey)
		userId, err := userIdByToken(string(jwttoken))
		if err != nil || userId == "" {
			c.Redirect(http.StatusOK, cRedirecURI("admin/login"))
			c.Abort() // 终止后续调用
			return
		}
		// 传递从jwttoken获取的userId
		c.Set(tokenUserId, userId)
	}
}

// funcIndexById 根据Id查询索引信息
func funcIndexById(ctx context.Context, c *app.RequestContext, htmlfile string) {
	id := c.Query("id")
	if id == "" {
		c.Redirect(http.StatusOK, cRedirecURI("admin/error"))
		c.Abort() // 终止后续调用
		return
	}
	urlPathParam := c.Param("urlPathParam")
	//indexName := bleveDataDir + urlPathParam
	responseData, err := funcIndexOne(urlPathParam, "*", id)
	if err != nil { //索引不存在
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
