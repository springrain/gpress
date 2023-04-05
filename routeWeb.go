package main

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
)

// 初始化函数
func init() {

	// 默认首页
	h.GET("/", funcIndex)

}

// funcIndex 模板首页
func funcIndex(ctx context.Context, c *app.RequestContext) {
	c.HTML(http.StatusOK, "index.html", map[string]string{"name": "test"})
}
