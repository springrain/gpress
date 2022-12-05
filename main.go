package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"

	"github.com/blevesearch/bleve/v2"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func main() {
	// server.Default() creates a Hertz with recovery middleware.
	// If you need a pure hertz, you can use server.New()
	h := server.Default(server.WithHostPorts(":8080"))

	// h.StaticFS("/", &app.FS{Root: datadir + "theme/default/default/*", GenerateIndexPages: true})

	h.GET("/hello", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "Hello hertz!")
	})
	h.SetFuncMap(template.FuncMap{"md5": MD5})
	h.LoadHTMLGlob(datadir + "theme/default/default/*")
	h.GET("/", IndexApi)
	h.GET("/test", func(ctx context.Context, c *app.RequestContext) {
		fmt.Println("1")
		r, err := findIndexFieldResult(ctx, indexNavMenuName, 1)
		fmt.Println(err)
		if err != nil {
			panic(err)
		}
		c.JSON(200, r)
	})

	// 测试新增数据
	h.GET("/add", func(ctx context.Context, c *app.RequestContext) {
		fmt.Println("1")
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
		c.JSON(200, r)
	})
	h.GET("/update", func(ctx context.Context, c *app.RequestContext) {
		fmt.Println("1")
		test := make(map[string]interface{})
		test["ID"] = "001"

		test["ChildTemplateID"] = "010202" // 子页面模板Id
		test["SortNo"] = "1"               // 排序
		// r := IndexMap[indexNavMenuName].Index("001", test)
		x := updateIndex(ctx, indexNavMenuName, "001", test)
		// m, _ := saveNexIndex(test, indexNavMenuName)
		c.JSON(200, x)
	})
	h.GET("/add2", func(ctx context.Context, c *app.RequestContext) {
		fmt.Println("1")
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
		c.JSON(200, m)
	})
	h.GET("/add3", func(ctx context.Context, c *app.RequestContext) {
		fmt.Println("1")
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
		c.JSON(200, m)
	})
	h.GET("/getThis", func(ctx context.Context, c *app.RequestContext) {
		fmt.Println("1")
		index := IndexMap[indexNavMenuName]
		queryIndexCode := bleve.NewNumericRangeInclusiveQuery(&active, &active, &inclusive, &inclusive)
		// 查询指定字段
		queryIndexCode.SetField("Active")
		// query := bleve.NewQueryStringQuery("")
		serarch := bleve.NewSearchRequestOptions(queryIndexCode, 1000, 0, false)
		// 查询所有字段
		serarch.Fields = []string{"*"}

		result, _ := index.SearchInContext(context.Background(), serarch)
		c.JSON(200, result)
	})

	h.GET("/getnav", func(ctx context.Context, c *app.RequestContext) {
		fmt.Println("1")
		result, _ := getNavMenu("0")
		c.JSON(200, result)
	})
	// 启动服务
	h.Spin()
}

// 请求响应函数

func IndexApi(ctx context.Context, c *app.RequestContext) {
	c.HTML(http.StatusOK, "index.html", map[string]string{"name": "test"})
}

// 自定义函数
func MD5(in string) ([]string, error) {
	list := make([]string, 2)

	hash := md5.Sum([]byte(in))
	list[0] = in
	list[1] = hex.EncodeToString(hash[:])
	return list, nil
}
