package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/blevesearch/bleve/v2"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

var templateDir = datadir + "template/"

// 模板路径,正常应该从siteInfo里获取,这里用于演示
var themePath = "theme/default/"

// hertz对象,可以在其他地方使用
var h *server.Hertz

func init() {
	h = server.Default(server.WithHostPorts(":8080"))
	//h.Use(gzip.Gzip(gzip.DefaultCompression))
}

func main() {

	funcMap := template.FuncMap{"md5": funcMD5, "basePath": funcBasePath}
	h.SetFuncMap(funcMap)

	//h.LoadHTMLFiles(themePath + "index.html")
	//h.LoadHTMLGlob(datadir + "html/theme/default/*")
	// 手动声明template对象,自己控制文件路径,默认是使用文件名,多个文件夹会存在问题
	tmpl := template.New("").Delims("", "").Funcs(funcMap)
	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		// 分隔符统一为 / 斜杠
		path = filepath.ToSlash(path)
		//相对路径
		relativePath := path[len(templateDir):]
		// 如果是静态资源
		if strings.Contains(path, "/js/") || strings.Contains(path, "/css/") || strings.Contains(path, "/image/") {
			h.StaticFile(relativePath, path)
		} else if strings.HasSuffix(path, ".html") { // 模板文件
			//创建对应的模板
			t := tmpl.New(relativePath)
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			//对应模板内容
			_, err = t.Parse(string(b))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		FuncLogError(err)
		panic(err)
	}
	//设置模板
	h.SetHTMLTemplate(tmpl)

	//设置默认的静态文件,实际路径会拼接为 datadir/public
	h.Static("/public", datadir)

	// 默认首页
	h.GET("/", funcIndex)

	h.GET("/hello", func(ctx context.Context, c *app.RequestContext) {
		c.String(http.StatusOK, "hello gpress!")
	})

	h.GET("/test", func(ctx context.Context, c *app.RequestContext) {
		r, err := findIndexFieldResult(ctx, indexNavMenuName, 1)
		if err != nil {
			FuncLogError(err)
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
		r := IndexMap[indexNavMenuName].Index("001", test)
		c.JSON(http.StatusOK, r)
	})
	h.GET("/update", func(ctx context.Context, c *app.RequestContext) {
		test := make(map[string]interface{})
		test["ID"] = "001"

		test["ChildTemplateID"] = "010202" // 子页面模板Id
		test["SortNo"] = "1"               // 排序
		// r := IndexMap[indexNavMenuName].Index("001", test)
		x := updateIndex(ctx, indexNavMenuName, "001", test)
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
		m, _ := saveNewIndex(ctx, test, indexNavMenuName)
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
		m, _ := saveNewIndex(ctx, test, indexNavMenuName)
		c.JSON(http.StatusOK, m)
	})
	h.GET("/getThis", func(ctx context.Context, c *app.RequestContext) {
		index := IndexMap[indexNavMenuName]
		queryIndexCode := bleve.NewNumericRangeInclusiveQuery(&active, &active, &inclusive, &inclusive)
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
		fmt.Println("1")
		result, _ := getNavMenu("0")
		c.JSON(http.StatusOK, result)
	})
	// 后台管理员登录
	h.GET("/admin/login", func(ctx context.Context, c *app.RequestContext) {
		c.HTML(http.StatusOK, "admin/login.html", nil)
	})
	h.POST("/admin/login", func(ctx context.Context, c *app.RequestContext) {
		userName := c.PostForm("userName")
		password := c.PostForm("password")
		fmt.Printf("userName:%s,password:%s", userName, password)
		c.JSON(http.StatusOK, "jwttoken-test")
	})

	// 后台管理员首页
	h.GET("/admin/index", func(ctx context.Context, c *app.RequestContext) {

		c.HTML(http.StatusOK, "admin/index.html", nil)
	})

	// 启动服务
	h.Spin()
}

// 请求响应函数
func funcIndex(ctx context.Context, c *app.RequestContext) {
	c.HTML(http.StatusOK, themePath+"index.html", map[string]string{"name": "test"})
}

// 测试自定义函数
func funcMD5(in string) ([]string, error) {
	list := make([]string, 2)

	hash := md5.Sum([]byte(in))
	list[0] = in
	list[1] = hex.EncodeToString(hash[:])
	return list, nil
}

//funcBasePath 基础路径,前端所有的资源请求必须带上 {{basePath}}
func funcBasePath() string {
	return ""
}
