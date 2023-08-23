package main

import (
	"fmt"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

// nats 服务器
var ns *server.Server

// nats客户端链接
var nc *nats.Conn

func initNatsServer() error {
	var err error
	// 创建一个 NATS Server 配置
	natsConfig := &server.Options{
		Host:           "127.0.0.1",
		Port:           4222,
		NoLog:          true,
		NoSigs:         true,
		MaxControlLine: 1024,
	}

	if config.ExternalNats { //如果是使用外部独立的Nats服务
		natsConfig, err = server.ProcessConfigFile(datadir + "nats.conf")
		if err != nil {
			return fmt.Errorf("server.ProcessConfigFile error: %w", err)
		}
	}

	// 创建 NATS Server 实例
	ns, err = server.NewServer(natsConfig)
	if err != nil {
		return fmt.Errorf("server.NewServer error: %w", err)
	}

	// 启动 NATS Server 实例
	//go func() {
	ns.Start()
	// Connect to a server
	nsUrl := fmt.Sprintf("nats://%s:%d", natsConfig.Host, natsConfig.Port)
	nc, err = nats.Connect(nsUrl)

	// Simple Publisher
	//nc.Publish("gpress", []byte("hello gpress"))

	// Simple Async Subscriber
	nc.Subscribe("gpress", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})
	nc.Publish("gpress", []byte("hello gpress"))
	//}()
	// 等待一段时间
	//time.Sleep(time.Second * 3)
	return err
}
