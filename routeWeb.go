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
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

// routeCategoryMap 动态添加的导航菜单路由map[pathURL]categoryID
var routeCategoryMap = make(map[string]string, 0)

// init 初始化函数
func init() {

	//初始化静态文件
	initStaticFS()

	// 异常页面
	h.GET("/error", funcError)

	// 默认首页
	h.GET("/", funcIndex)
	h.GET("/page/:pageNo", funcIndex)
	h.GET("/page/:pageNo/", funcIndex)

	// 导航菜单列表
	h.GET("/category/:urlPathParam", funcListCategory)
	h.GET("/category/:urlPathParam/", funcListCategory)
	h.GET("/category/:urlPathParam/page/:pageNo", funcListCategory)
	h.GET("/category/:urlPathParam/page/:pageNo/", funcListCategory)

	// 查看标签
	h.GET("/tag/:urlPathParam", funcListTags)
	h.GET("/tag/:urlPathParam/", funcListTags)
	h.GET("/tag/:urlPathParam/page/:pageNo", funcListTags)
	h.GET("/tag/:urlPathParam/page/:pageNo/", funcListTags)
	// 查看内容
	h.GET("/post/:urlPathParam", funcOneContent)
	h.GET("/post/:urlPathParam/", funcOneContent)

	//@TODO 静态文件映射的 /favicon.ico 还是进入通配funcListCategoryFilepath
	h.GET("/favicon.ico", func(ctx context.Context, c *app.RequestContext) {
		c.File(datadir + site.Favicon)
	})

	//初始化导航菜单路由
	initCategoryRoute()

	// 通配其他动态路径
	h.GET("/*filepath", funcListCategoryFilepath)

}

// funcIndex 模板首页
func funcIndex(ctx context.Context, c *app.RequestContext) {
	data := warpRequestMap(c)
	cHtml(c, http.StatusOK, "index.html", data)
}

// funcError 错误页面
func funcError(ctx context.Context, c *app.RequestContext) {
	cHtml(c, http.StatusOK, "error.html", nil)
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
	cHtml(c, http.StatusOK, templateFile, data)
}

// funcListTags 标签列表
func funcListTags(ctx context.Context, c *app.RequestContext) {
	data := warpRequestMap(c)
	urlPathParam := c.Param("urlPathParam")

	data["UrlPathParam"] = urlPathParam
	cHtml(c, http.StatusOK, "tag.html", data)
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
	cHtml(c, http.StatusOK, templateFile, data)
}

// funcListCategoryFilepath 通配的filepath映射
func funcListCategoryFilepath(ctx context.Context, c *app.RequestContext) {
	//@TODO 静态文件映射的 /favicon.ico 还是进入到这个方法,造成了异常
	key := string(c.URI().Path())
	key = trimRightSlash(key) // 去掉最后的/, 例如: /web/ 实际是 /web
	//从url路径分析获得的内容id,例如: /web/nginx-use-hsts contentID是nginx-use-hsts
	contentID := ""
	pageNo := ""
	//获取路径的对应的 categoryID
	categoryID, has := routeCategoryMap[key]
	if has { //导航菜单的路径
		c.Set("urlPathParam", categoryID)
		funcListCategory(ctx, c)
		return
	}

	//拆分url
	urls := strings.Split(key, "/")

	//处理分页请求,例如 /web/page/1
	if len(urls) > 1 && urls[len(urls)-2] == "page" {
		pageNo = urls[len(urls)-1]
		urls = urls[:len(urls)-2]
		key = strings.Join(urls, "/")
		categoryID, has = routeCategoryMap[key]
		if has {
			c.Set("urlPathParam", categoryID)
			c.Set("pageNo", pageNo)
			funcListCategory(ctx, c)
			return
		}
	}
	if len(urls) < 3 { // 类似 /web 却没有注册,返回404
		cHtml(c, http.StatusNotFound, "error.html", nil)
		return
	}
	//处理导航/内容的请求,例如: /web/nginx-use-hsts
	contentID = urls[len(urls)-1]
	urls = urls[:len(urls)-1]
	key = strings.Join(urls, "/")
	//获取 categoryID
	categoryID, has = routeCategoryMap[key]

	if !has { //没有注册的categoryID,返回404的error
		cHtml(c, http.StatusNotFound, "error.html", nil)
		return
	}
	if contentID != "" { //内容页面
		c.Set("urlPathParam", contentID)
		funcOneContent(ctx, c)
	} else { //导航菜单页面
		c.Set("urlPathParam", categoryID)
		funcListCategory(ctx, c)
	}

}

// initCategoryRoute 初始化导航菜单的映射路径
func initCategoryRoute() {
	categorys, _ := findAllCategory(context.Background())
	for i := 0; i < len(categorys); i++ {
		category := categorys[i]
		//导航菜单的访问映射
		h.GET(trimRightSlash(category.PathURL), addListCategoryRoute(category.Id))
		h.GET(category.PathURL, addListCategoryRoute(category.Id))
		//导航菜单分页数据的访问映射
		h.GET(category.PathURL+"page/:pageNo", addListCategoryRoute(category.Id))
		h.GET(category.PathURL+"page/:pageNo/", addListCategoryRoute(category.Id))
		//导航菜单下文章的访问映射
		h.GET(category.PathURL+":urlPathParam", funcOneContent)
		h.GET(category.PathURL+":urlPathParam/", funcOneContent)
	}
}

// addListCategoryRoute 增加导航菜单的GET请求路由,用于自定义设置导航的路由
func addListCategoryRoute(categoryID string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Set("urlPathParam", categoryID)
		funcListCategory(ctx, c)
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
