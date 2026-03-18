// Copyright (c) 2023 gpress Authors.
//
// This file is part of gpress.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

// init 初始化函数
func init() {

	//初始化静态文件
	initStaticFS()

	// 注册中间件, 指定处理.md后缀
	h.Use(removeExtMiddleware(".md"))

	// 异常页面
	h.GET("/error", funcError)
	h.GET("/error.md", funcError)

	// 默认首页
	h.GET("/", funcIndex)
	h.GET("/index.md", funcIndex)
	h.GET("/page/:pageNo", funcIndex)
	h.GET("/page/:pageNo/", funcIndex)

	// 查看标签
	h.GET("/tag/:urlPathParam", funcListTags)
	h.GET("/tag/:urlPathParam/", funcListTags)
	h.GET("/tag/:urlPathParam/page/:pageNo", funcListTags)
	h.GET("/tag/:urlPathParam/page/:pageNo/", funcListTags)

	//初始化导航菜单路由
	initCategoryRoute()

}

// funcIndex 模板首页
func funcIndex(ctx context.Context, c *app.RequestContext) {
	data := warpRequestMap(c)
	templateFile := "index.html"
	if _, has := c.Get(markdownKey); has {
		templateFile = "index.md"
	}
	cHtml(c, http.StatusOK, templateFile, data)
}

// funcError 错误页面
func funcError(ctx context.Context, c *app.RequestContext) {
	templateFile := "error.html"
	if _, has := c.Get(markdownKey); has {
		templateFile = "error.md"
	}
	cHtml(c, http.StatusOK, templateFile, nil)
}

// funcListCategory 导航菜单数据列表
func funcListCategory(ctx context.Context, c *app.RequestContext) {
	data := warpRequestMap(c)
	urlPathParam := c.Param("urlPathParam")
	if urlPathParam == "" { //导航菜单路径访问,例如:/web
		urlPathParam = c.GetString("urlPathParam")
	}
	data["UrlPathParam"] = urlPathParam
	templateFile, err := findThemeTemplate(ctx, tableCategoryName, urlPathParam)
	if err != nil || templateFile == "" {
		templateFile = "category.html"
	}
	if _, has := c.Get(markdownKey); has {
		templateFile = strings.TrimSuffix(templateFile, ".html") + ".md"
	}
	cHtml(c, http.StatusOK, templateFile, data)
}

// funcListTags 标签列表
func funcListTags(ctx context.Context, c *app.RequestContext) {
	data := warpRequestMap(c)
	urlPathParam := c.Param("urlPathParam")

	data["UrlPathParam"] = urlPathParam

	templateFile := "tag.html"
	if _, has := c.Get(markdownKey); has {
		templateFile = "tag.md"
	}
	cHtml(c, http.StatusOK, templateFile, data)
}

// funcOneContent 查询一篇文章
func funcOneContent(ctx context.Context, c *app.RequestContext) {
	data := warpRequestMap(c)
	urlPathParam := c.Param("urlPathParam")
	if urlPathParam == "" { //导航菜单路径访问,例如:/web/nginx-use-hsts
		urlPathParam = c.GetString("urlPathParam")
	}
	data["UrlPathParam"] = urlPathParam

	templateFile, err := findThemeTemplate(ctx, tableContentName, urlPathParam)
	if err != nil || templateFile == "" {
		templateFile = "content.html"
	}
	if _, has := c.Get(markdownKey); has {
		templateFile = strings.TrimSuffix(templateFile, ".html") + ".md"
	}
	cHtml(c, http.StatusOK, templateFile, data)
}

// initCategoryRoute 初始化导航菜单的映射路径
func initCategoryRoute() {
	categories, _ := findAllCategory(context.Background())
	for i := 0; i < len(categories); i++ {
		category := categories[i]
		categoryID := category.Id
		addCategoryRoute(categoryID)
	}
}

// addCategoryRoute 增加导航菜单的路由
func addCategoryRoute(categoryID string) {

	// 处理重复注册路由的panic,不对外抛出
	defer func() {
		if r := recover(); r != nil {
			FuncLogPanic(nil, fmt.Errorf("addCategoryRoute panic recovered: %v", r))
		}
	}()

	//导航菜单的访问映射
	h.GET(funcTrimSuffixSlash(categoryID), addListCategoryRoute(categoryID))
	h.GET(funcTrimSuffixSlash(categoryID)+".md", addListCategoryRoute(categoryID))
	h.GET(categoryID, addListCategoryRoute(categoryID))
	//导航菜单分页数据的访问映射
	h.GET(categoryID+"page/:pageNo", addListCategoryRoute(categoryID))
	h.GET(categoryID+"page/:pageNo/", addListCategoryRoute(categoryID))
	//导航菜单下文章的访问映射
	h.GET(categoryID+":contentURI", addOneContentRoute(categoryID))
	h.GET(categoryID+":contentURI/", addOneContentRoute(categoryID))
}

// addListCategoryRoute 增加导航菜单的GET请求路由,用于自定义设置导航的路由
func addListCategoryRoute(categoryID string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Set("urlPathParam", categoryID)
		funcListCategory(ctx, c)
	}
}

// addOneContentRoute 增加内容的GET请求路由
func addOneContentRoute(categoryID string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		contentURI := c.Param("contentURI")
		key := categoryID + contentURI
		c.Set("urlPathParam", key)
		funcOneContent(ctx, c)
	}
}

// warpRequestMap 包装请求参数为map
func warpRequestMap(c *app.RequestContext) map[string]interface{} {
	pageNoStr := c.Param("pageNo")
	if pageNoStr == "" {
		pageNoStr = c.GetString("pageNo")
	}
	if pageNoStr == "" {
		//pageNoStr = c.DefaultQuery("pageNo", "1")
		pageNoStr = "1"
	}

	pageNo, _ := strconv.Atoi(pageNoStr)
	q := strings.TrimSpace(c.Query("q"))

	data := make(map[string]interface{}, 0)
	data["pageNo"] = pageNo
	data["q"] = q
	//设置用户角色,0是访客,1是管理员
	userType, ok := c.Get(userTypeKey)
	if ok {
		data[userTypeKey] = userType
	} else {
		data[userTypeKey] = 0
	}
	return data
}

// removeExtMiddleware 仅处理.md结尾的请求,非.md直接放行
func removeExtMiddleware(ext string) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		originalPath := string(ctx.Request.URI().Path())

		// 核心判断: 仅当路径非空/非根路径 且 以指定后缀结尾时才处理
		if originalPath == "" || originalPath == "/" || !strings.HasSuffix(originalPath, ext) {
			ctx.Next(c) // 不是.md结尾,直接放行
			return
		}

		// 1. 修改URI路径(去掉后缀)
		newPath := originalPath[:len(originalPath)-len(ext)]
		ctx.Request.URI().SetPath(newPath)

		// 2. 设置标记: 标记该请求是.md后缀转换而来(存在ctx中,处理器可读取)
		ctx.Set(markdownKey, true)

		// 3. 修改已解析的Params参数值(移除后缀)
		for i := range ctx.Params {
			ctx.Params[i].Value = strings.TrimSuffix(ctx.Params[i].Value, ext)
		}

		// 继续执行后续逻辑
		ctx.Next(c)
	}
}
