package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	//设置自定义i函数.要在LoadHTMLGlob加载模板之前.
	router.SetFuncMap(template.FuncMap{"md5": MD5})

	//router.GET("/", IndexApi)
	//LoadHTMLFiles(templates...)
	//router.LoadHTMLGlob("templates/**/*")
	router.LoadHTMLGlob(datadir + "theme/default/default/*")

	router.GET("/test", func(c *gin.Context) {
		fmt.Println("1")
		r, err := findIndexFields(navMenuName, 1)
		fmt.Println(err)
		if err != nil {
			panic(err)
		}
		c.JSON(200, r)
	})

	router.GET("/add", func(c *gin.Context) {
		fmt.Println("1")
		test := make(map[string]interface{})
		test["MenuName"] = "测试菜单"
		test["HrefURL"] = "localhost:8080"
		test["HrefTarget"] = "跳转方式"
		test["PID"] = "0"
		test["ComCode"] = "阿斯弗,sfs"
		test["TemplateID"] = "模板Id"
		test["ModuleIndexCode"] = "ModuleIndexCode"
		test["ChildTemplateID"] = "子页面模板Id"
		test["ComCode"] = "SortNo"
		test["Active"] = "1"
		test["themePC"] = "PC主题"
		m, _ := saveNexIndex(test, navMenuName)
		c.JSON(200, m)
	})
	router.GET("/getThis", func(c *gin.Context) {
		fmt.Println("1")
		index := IndexMap[navMenuName]
		queryIndexCode := bleve.NewTermQuery("测试菜单")
		//查询指定字段
		queryIndexCode.SetField("MenuName")
		query := bleve.NewConjunctionQuery(queryIndexCode)
		serarch := bleve.NewSearchRequestOptions(query, 1000, 0, false)
		//查询所有字段
		serarch.Fields = []string{"*"}

		result, _ := index.SearchInContext(context.Background(), serarch)
		c.JSON(200, result)
	})

	router.Run(":8080")
}

//请求响应函数
func IndexApi(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"name": "test"})

}

//自定义函数
func MD5(in string) ([]string, error) {
	list := make([]string, 2)

	hash := md5.Sum([]byte(in))
	list[0] = in
	list[1] = hex.EncodeToString(hash[:])
	return list, nil
}
