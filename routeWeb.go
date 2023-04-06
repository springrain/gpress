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

}

// funcIndex 模板首页
func funcIndex(ctx context.Context, c *app.RequestContext) {
	pageNoStr := c.DefaultQuery("pageNo", "1")
	pageNo, _ := strconv.Atoi(pageNoStr)
	q := strings.TrimSpace(c.Query("q"))

	data := make(map[string]interface{}, 0)
	data["pageNo"] = pageNo
	data["q"] = q
	c.HTML(http.StatusOK, "index.html", data)
}
