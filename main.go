package main

import (
	"crypto/md5"
	"encoding/hex"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	//设置自定义i函数.要在LoadHTMLGlob加载模板之前.
	router.SetFuncMap(template.FuncMap{"md5": MD5})

	router.GET("/", IndexApi)
	//LoadHTMLFiles(templates...)
	//router.LoadHTMLGlob("templates/**/*")
	router.LoadHTMLGlob("templates/*")

	router.Run(":8080")
}

//请求响应函数
func IndexApi(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"name": "test"})
	return
}

//自定义函数
func MD5(in string) ([]string, error) {
	list := make([]string, 2)

	hash := md5.Sum([]byte(in))
	list[0] = in
	list[1] = hex.EncodeToString(hash[:])
	return list, nil
}
