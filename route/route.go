package route

import (
	"context"
	"gitee.com/gpress/gpress/constant"
	"gitee.com/gpress/gpress/gbleve"
	"gitee.com/gpress/gpress/logger"
	"gitee.com/gpress/gpress/service"
	"gitee.com/gpress/gpress/util"
	"github.com/blevesearch/bleve/v2"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol"
	"html/template"
	"net/http"
)

var ThemePath string
var h *server.Hertz
var funcMap = template.FuncMap{"md5": util.FuncMD5, "basePath": funcBasePath, "T": util.FuncT, "safeHTML": funcSafeHTML, "relURL": funcRelURL, "sass": funcSass}

func InitRoute(port string) {
	h = server.Default(server.WithHostPorts(port), server.WithBasePath("/"))
}
func RunServer() {

	h.GET("/", funcIndex)
	h.GET("/hello", func(ctx context.Context, c *app.RequestContext) {
		c.String(http.StatusOK, "hello gpress!")
	})
	h.GET("/test", func(ctx context.Context, c *app.RequestContext) {
		r, err := service.FindIndexFieldResult(ctx, constant.INDEX_NAV_MENU_NAME, 1)
		if err != nil {
			logger.FuncLogError(err)
			panic(err)
		}
		c.JSON(http.StatusOK, r)
	})

	// 测试新增数据
	h.GET("/add", func(ctx context.Context, c *app.RequestContext) {
		test := make(map[string]interface{}) // 新建map
		test["MenuName"] = "一级菜单名称"
		test["HrefURL"] = "localhost:8080"
		test["HrefTarget"] = "跳转方式"
		test["PID"] = "0" // 顶级菜单目录
		test["ComCode"] = "使用逗号分割,字符串,测试"
		test["Active"] = 1    // 是否有效
		test["themePC"] = "1" // 是否pc主题
		test["ModuleIndexCode"] = "Module的索引名称"
		test["TemplateID"] = "010101"      // 模板Id
		test["ChildTemplateID"] = "010201" // 子页面模板Id
		test["SortNo"] = "1"               // 排序
		test["ID"] = "001"
		// m, _ := saveNewIndex(c.Request.Context(), test, indexNavMenuName)
		r := gbleve.IndexMap[constant.INDEX_NAV_MENU_NAME].Index("001", test)
		c.JSON(http.StatusOK, r)
	})
	h.GET("/update", func(ctx context.Context, c *app.RequestContext) {
		test := make(map[string]interface{})
		test["ID"] = "001"

		test["ChildTemplateID"] = "010202" // 子页面模板Id
		test["SortNo"] = "1"               // 排序
		// r := IndexMap[indexNavMenuName].Index("001", test)
		x := service.UpdateIndex(ctx, constant.INDEX_NAV_MENU_NAME, "001", test)
		// m, _ := saveNexIndex(test, indexNavMenuName)
		c.JSON(http.StatusOK, x)
	})
	h.GET("/add2", func(ctx context.Context, c *app.RequestContext) {
		test := make(map[string]interface{})
		test["MenuName"] = "一级菜单"
		test["HrefURL"] = "localhost:8080"
		test["HrefTarget"] = "跳转方式"
		test["PID"] = "0"
		test["ComCode"] = "阿斯弗,sfs"
		test["TemplateID"] = "模板Id"
		test["ModuleIndexCode"] = "ModuleIndexCode"
		test["ChildTemplateID"] = "子页面模板Id"
		test["Active"] = 1
		test["themePC"] = "PC主题"
		m, _ := service.SaveNewIndex(ctx, test, constant.INDEX_NAV_MENU_NAME)
		c.JSON(http.StatusOK, m)
	})

	h.GET("/add3", func(ctx context.Context, c *app.RequestContext) {
		test := make(map[string]interface{})
		test["MenuName"] = "我是个子集"
		test["HrefURL"] = "localhost:8080"
		test["HrefTarget"] = "跳转方式"
		test["PID"] = "7216c38e-78fb-4ad9-95bf-294582faa685"
		test["ComCode"] = "阿斯弗,sfs"
		test["TemplateID"] = "模板Id"
		test["ModuleIndexCode"] = "ModuleIndexCode"
		test["ChildTemplateID"] = "子页面模板Id"
		test["Active"] = 1
		test["themePC"] = "PC主题"
		m, _ := service.SaveNewIndex(ctx, test, constant.INDEX_NAV_MENU_NAME)
		c.JSON(http.StatusOK, m)
	})

	h.GET("/getThis", func(ctx context.Context, c *app.RequestContext) {
		index := gbleve.IndexMap[constant.INDEX_NAV_MENU_NAME]
		queryIndexCode := bleve.NewNumericRangeInclusiveQuery(&service.Active, &service.Active, &service.Inclusive, &service.Inclusive)
		// 查询指定字段
		queryIndexCode.SetField("Active")
		// query := bleve.NewQueryStringQuery("")
		serarch := bleve.NewSearchRequestOptions(queryIndexCode, 1000, 0, false)
		// 查询所有字段
		serarch.Fields = []string{"*"}

		result, _ := index.SearchInContext(context.Background(), serarch)
		c.JSON(http.StatusOK, result)
	})

	h.GET("/getnav", func(ctx context.Context, c *app.RequestContext) {
		result, _ := service.GetNavMenu("0")
		c.JSON(http.StatusOK, result)
	})

	// 安装
	h.GET("/admin/install", func(ctx context.Context, c *app.RequestContext) {
		if gbleve.Installed { // 如果已经安装过了,跳转到登录
			c.Redirect(http.StatusOK, []byte("/admin/login"))
			return
		}
		c.HTML(http.StatusOK, "/admin/install.html", nil)
	})
	h.POST("/admin/install", func(ctx context.Context, c *app.RequestContext) {
		if gbleve.Installed { // 如果已经安装过了,跳转到登录
			c.Redirect(http.StatusOK, []byte("/admin/login"))
			return
		}
		// 使用后端管理界面配置,jwtSecret也有后端随机产生
		account := c.PostForm("account")
		password := c.PostForm("password")
		err := service.InsertUser(ctx, account, password)
		if err != nil {
			c.Redirect(http.StatusOK, []byte("/admin/error"))
			return
		}
		// 安装成功,更新安装状态
		service.UpdateInstall(ctx)
		c.Redirect(http.StatusOK, []byte("/admin/login"))
	})

	// 后台管理员登录
	h.GET("/admin/login", func(ctx context.Context, c *app.RequestContext) {
		if !gbleve.Installed { // 如果没有安装,跳转到安装
			c.Redirect(http.StatusOK, []byte("/admin/install"))
			return
		}
		c.HTML(http.StatusOK, "/admin/login.html", nil)
	})
	h.POST("/admin/login", func(ctx context.Context, c *app.RequestContext) {
		if !gbleve.Installed { // 如果没有安装,跳转到安装
			c.Redirect(http.StatusOK, []byte("/admin/install"))
			return
		}
		account := c.PostForm("account")
		password := c.PostForm("password")
		userId, err := service.FindUserId(ctx, account, password)
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
		c.SetCookie(service.Config.JwtTokenKey, jwttoken, service.Config.Timeout, "/", "", protocol.CookieSameSiteStrictMode, true, true)

		c.Redirect(http.StatusOK, []byte("/admin/index"))
	})
	// admin路由组
	admin := h.Group("/admin")
	admin.Use(adminHandler())
	// 后台管理员首页
	admin.GET("/index", func(ctx context.Context, c *app.RequestContext) {
		// 获取从jwttoken中解码的userId
		userId, ok := c.Get(constant.TOKEN_USER_ID)
		if !ok || userId == "" {
			c.Redirect(http.StatusOK, []byte("/admin/login"))
			return
		}
		c.HTML(http.StatusOK, "/admin/index.html", nil)
	})
	h.Spin()
}
