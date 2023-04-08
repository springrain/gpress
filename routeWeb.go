package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

// 初始化函数
func init() {

	// 默认首页
	h.GET("/", funcIndex)
	// 导航菜单列表
	h.GET("/navMenu/:urlPathParam", funcListNavMenu)
	// 查看标签
	h.GET("/tag/:urlPathParam", funcListTags)
	// 查看内容
	h.GET("/post/:urlPathParam", funcOneContent)

}

// funcIndex 模板首页
func funcIndex(ctx context.Context, c *app.RequestContext) {
	data := warpRequestMap(c)
	c.HTML(http.StatusOK, "index.html", data)
}

func funcListNavMenu(ctx context.Context, c *app.RequestContext) {
	data := warpRequestMap(c)
	data["urlPathParam"] = c.Param("urlPathParam")
	c.HTML(http.StatusOK, "navMenu.html", data)
}
func funcListTags(ctx context.Context, c *app.RequestContext) {
	data := warpRequestMap(c)
	data["urlPathParam"] = c.Param("urlPathParam")
	c.HTML(http.StatusOK, "tag.html", data)
}
func funcOneContent(ctx context.Context, c *app.RequestContext) {
	data := warpRequestMap(c)
	data["urlPathParam"] = c.Param("urlPathParam")
	c.HTML(http.StatusOK, "content.html", data)
}

func warpRequestMap(c *app.RequestContext) map[string]interface{} {
	pageNoStr := c.DefaultQuery("pageNo", "1")
	pageNo, _ := strconv.Atoi(pageNoStr)
	q := strings.TrimSpace(c.Query("q"))

	data := make(map[string]interface{}, 0)
	data["pageNo"] = pageNo
	data["q"] = q
	return data
}
