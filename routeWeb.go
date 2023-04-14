package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// 初始化函数
func init() {

	// 默认首页
	h.GET("/", funcTable)
	// 导航菜单列表
	h.GET("/navMenu/:urlPathParam", funcListNavMenu)
	// 查看标签
	h.GET("/tag/:urlPathParam", funcListTags)
	// 查看内容
	h.GET("/post/:urlPathParam", funcOneContent)

}

// funcTable 模板首页
func funcTable(ctx context.Context, c *app.RequestContext) {
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

// hrefURLRoute href 需要跳转的地址,hrefURL原地址
func hrefURLRoute(href string, realURL string) error {
	if href == "" || realURL == "" {
		return errors.New("跳转路径为空")
	}

	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") { //外部跳转
		h.GET("/"+realURL, func(ctx context.Context, c *app.RequestContext) { //注册内部地址,解析跳转到外部
			c.Redirect(consts.StatusMovedPermanently, []byte(href))
			c.Abort() // 终止后续调用
		})
		return nil
	}
	//内部跳转, 跳转到内部路由,例如 navMenu/about 跳转到 content/about
	h.GET(href, func(ctx context.Context, c *app.RequestContext) {
		c.Redirect(consts.StatusFound, cRedirecURI(config.BasePath+realURL))
		c.Abort() // 终止后续调用
	})

	return nil
}
