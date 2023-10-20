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
	h.GET("/category/:urlPathParam", funcListCategory)
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

func funcListCategory(ctx context.Context, c *app.RequestContext) {
	data := warpRequestMap(c)
	urlPathParam := c.Param("urlPathParam")
	data["UrlPathParam"] = urlPathParam
	templateFile, err := findPageTemplate(ctx, "category", urlPathParam)
	if err != nil || templateFile == "" {
		templateFile = "category.html"
	}

	c.HTML(http.StatusOK, templateFile, data)
}
func funcListTags(ctx context.Context, c *app.RequestContext) {
	data := warpRequestMap(c)
	urlPathParam := c.Param("urlPathParam")
	data["UrlPathParam"] = urlPathParam

	c.HTML(http.StatusOK, "tag.html", data)
}
func funcOneContent(ctx context.Context, c *app.RequestContext) {
	data := warpRequestMap(c)
	urlPathParam := c.Param("urlPathParam")
	data["UrlPathParam"] = urlPathParam

	templateFile, err := findPageTemplate(ctx, "content", urlPathParam)
	if err != nil || templateFile == "" {
		templateFile = "content.html"
	}

	whereCondition, ok := c.Get(whereConditionKey)
	if ok {
		data[whereConditionKey] = whereCondition
	} else {
		data[whereConditionKey] = " and status<2 "
	}

	c.HTML(http.StatusOK, templateFile, data)
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

/*
// hrefURLRoute href 需要跳转的地址,hrefURL原地址
func hrefURLRoute(realURL string, hrefURL string) error {
	if hrefURL == "" || realURL == "" {
		return errors.New("跳转路径为空")
	}

	if strings.HasPrefix(hrefURL, "http://") || strings.HasPrefix(hrefURL, "https://") { //外部跳转
		h.GET("/"+realURL, func(ctx context.Context, c *app.RequestContext) { //注册内部地址,解析跳转到外部
			c.Redirect(consts.StatusMovedPermanently, []byte(hrefURL))
			c.Abort() // 终止后续调用
		})
		return nil
	}
	//内部跳转, 跳转到内部路由,例如 category/about 跳转到 content/about
	h.GET("/"+realURL, func(ctx context.Context, c *app.RequestContext) {
		//https://github.com/cloudwego/hertz/issues/724
		c.Redirect(http.StatusOK, cRedirecURI(hrefURL))
		//c.Redirect(consts.StatusFound, []byte("/"+hrefURL))
		c.Abort() // 终止后续调用
	})

	return nil
}
*/
