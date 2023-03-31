package main

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/route"
)

// adminGroup路由组,使用变量声明,优先级高于init函数
var adminGroup = initAdminGroup()

func initAdminGroup() *route.RouterGroup {
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

	// 默认首页
	h.GET("/", funcIndex)

	h.GET("/getnav", func(ctx context.Context, c *app.RequestContext) {
		result, _ := getNavMenu("0")
		c.JSON(http.StatusOK, result)
	})

	// 异常页面
	h.GET("/admin/error", func(ctx context.Context, c *app.RequestContext) {
		c.HTML(http.StatusOK, "/admin/error.html", nil)
	})

	// 安装
	h.GET("/admin/install", func(ctx context.Context, c *app.RequestContext) {
		if installed { // 如果已经安装过了,跳转到登录
			c.Redirect(http.StatusOK, []byte("/admin/login"))
			return
		}
		c.HTML(http.StatusOK, "/admin/install.html", nil)
	})
	h.POST("/admin/install", func(ctx context.Context, c *app.RequestContext) {
		if installed { // 如果已经安装过了,跳转到登录
			c.Redirect(http.StatusOK, []byte("/admin/login"))
			return
		}
		// 使用后端管理界面配置,jwtSecret也有后端随机产生
		account := c.PostForm("account")
		password := c.PostForm("password")
		err := insertUser(ctx, account, password)
		if err != nil {
			c.Redirect(http.StatusOK, []byte("/admin/error"))
			return
		}
		// 安装成功,更新安装状态
		updateInstall(ctx)
		c.Redirect(http.StatusOK, []byte("/admin/login"))
	})

	// 后台管理员登录
	h.GET("/admin/login", func(ctx context.Context, c *app.RequestContext) {
		if !installed { // 如果没有安装,跳转到安装
			c.Redirect(http.StatusOK, []byte("/admin/install"))
			return
		}
		c.HTML(http.StatusOK, "/admin/login.html", nil)
	})
	h.POST("/admin/login", func(ctx context.Context, c *app.RequestContext) {
		if !installed { // 如果没有安装,跳转到安装
			c.Redirect(http.StatusOK, []byte("/admin/install"))
			return
		}
		account := c.PostForm("account")
		password := c.PostForm("password")
		userId, err := findUserId(ctx, account, password)
		if userId == "" || err != nil { // 用户不存在或者异常
			c.Redirect(http.StatusOK, []byte("/admin/login"))
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

		c.Redirect(http.StatusOK, []byte("/admin/index"))
	})

	// 后台管理员首页
	adminGroup.GET("/index", func(ctx context.Context, c *app.RequestContext) {
		// 获取从jwttoken中解码的userId
		userId, ok := c.Get(tokenUserId)
		if !ok || userId == "" {
			c.Redirect(http.StatusOK, []byte("/admin/login"))
			return
		}

		c.HTML(http.StatusOK, "/admin/index.html", nil)
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
		//c.HTML(http.StatusOK, "/admin/index.html", nil)
		c.JSON(http.StatusOK, ResponseData{StatusCode: 1})
	})

	// 通用列表
	adminGroup.GET("/:indexName/list", funcList)
	adminGroup.POST("/:indexName/list", funcList)

}

// funcIndex 模板首页
func funcIndex(ctx context.Context, c *app.RequestContext) {
	c.HTML(http.StatusOK, themePath+"index.html", map[string]string{"name": "test"})
}

// funcList 通用list列表
func funcList(ctx context.Context, c *app.RequestContext) {
	nameParam := c.Param("indexName")
	indexName := bleveDataDir + nameParam
	reponseData, err := findIndex(ctx, c, indexName)
	if err != nil { //索引不存在
		c.Redirect(http.StatusOK, []byte("/admin/error"))
		c.Abort() // 终止后续调用
		return
	}

	//优先使用自定义模板文件
	listFile := "/admin/" + nameParam + "List.html"
	t := tmpl.Lookup(listFile)
	if t == nil { //不存在自定义模板,使用通用模板
		listFile = "/admin/list.html"
	}

	c.HTML(http.StatusOK, listFile, reponseData)
}

// permissionHandler 权限拦截器
func permissionHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		jwttoken := c.Cookie(config.JwttokenKey)
		userId, err := userIdByToken(string(jwttoken))
		if err != nil || userId == "" {
			c.Redirect(http.StatusOK, []byte("/admin/login"))
			c.Abort() // 终止后续调用
			return
		}
		// 传递从jwttoken获取的userId
		c.Set(tokenUserId, userId)
	}
}
