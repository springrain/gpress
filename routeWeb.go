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

// 初始化函数
func init() {
	// 异常页面
	h.GET("/error", func(ctx context.Context, c *app.RequestContext) {
		c.HTML(http.StatusOK, "error.html", nil)
	})

	// 默认首页
	h.GET("/", funcIndex)
	h.GET("/page/:pageNo", funcIndex)

	// 导航菜单列表
	h.GET("/category/:urlPathParam", funcListCategory)
	h.GET("/category/:urlPathParam/page/:pageNo", funcListCategory)

	// 查看标签
	h.GET("/tag/:urlPathParam", funcListTags)
	h.GET("/tag/:urlPathParam/page/:pageNo", funcListTags)
	// 查看内容
	h.GET("/post/:urlPathParam", funcOneContent)

}

// funcIndex 模板首页
func funcIndex(ctx context.Context, c *app.RequestContext) {
	data := warpRequestMap(c)
	c.HTML(http.StatusOK, "index.html", data)
}

func funcListCategory(ctx context.Context, c *app.RequestContext) {
	urlPathParam := c.Param("urlPathParam")
	/*
		status, err := findStatus(ctx, tableCategoryName, urlPathParam)
		if err != nil || status < 0 || status > 1 { //异常状态
			c.Redirect(http.StatusOK, cRedirecURI("error"))
			c.Abort() // 终止后续调用
			return
		}
	*/
	data := warpRequestMap(c)

	data["UrlPathParam"] = urlPathParam
	templateFile, err := findThemeTemplate(ctx, tableCategoryName, urlPathParam)
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
	urlPathParam := c.Param("urlPathParam")
	/*
		status, err := findStatus(ctx, tableContentName, urlPathParam)
		if err != nil || status < 0 || status > 1 { //异常状态
			c.Redirect(http.StatusOK, cRedirecURI("error"))
			c.Abort() // 终止后续调用
			return
		}
	*/
	data := warpRequestMap(c)
	data["UrlPathParam"] = urlPathParam

	templateFile, err := findThemeTemplate(ctx, tableContentName, urlPathParam)
	if err != nil || templateFile == "" {
		templateFile = "content.html"
	}

	c.HTML(http.StatusOK, templateFile, data)
}

func warpRequestMap(c *app.RequestContext) map[string]interface{} {
	pageNoStr := c.Param("pageNo")
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
/*
func findStatus(ctx context.Context, tableName string, id string) (int, error) {
	if tableName == "" || id == "" {
		return -1, errors.New("数据异常")
	}
	finder := zorm.NewSelectFinder(tableName, "status").Append("WHERE id=?", id)
	stauts := -1
	ok, err := zorm.QueryRow(ctx, finder, &stauts)
	if !ok {
		return -1, errors.New("数据异常")
	}
	return stauts, err
}
*/
