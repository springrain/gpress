package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func init() {
	// 异常页面
	h.GET("/nats", func(ctx context.Context, c *app.RequestContext) {
		subject := c.Query("subject")
		message := c.Query("message")
		nc.Publish(subject, []byte(message))
	})
}
