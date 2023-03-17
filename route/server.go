package route

/**
路由包
*/
import (
	"gitee.com/gpress/gpress/constant"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func RunServer() {
	h := server.Default(server.WithHostPorts(constant.ServerPort), server.WithBasePath("/"))
	h.GET("/", funcIndex)

	h.Spin()
}
