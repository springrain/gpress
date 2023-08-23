package main

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// nats 服务器
var ns *server.Server

// nats客户端链接
var nc *nats.Conn

var js jetstream.JetStream

func initNatsServer() error {
	var err error
	// 创建一个 NATS Server 配置
	natsConfig := &server.Options{
		Host:           "127.0.0.1",
		Port:           4222,
		JetStream:      true,
		StoreDir:       datadir + "natsdata",
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
	if err != nil {
		return fmt.Errorf("nats.Connect(nsUrl) error: %w", err)
	}
	// create jetstream context from nats connection
	js, err = jetstream.New(nc)
	if err != nil {
		return fmt.Errorf("jetstream.New(nc) error: %w", err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	// get existing stream handle
	stream, _ := js.Stream(ctx, "gpress")
	if stream == nil {
		stream, err = js.CreateStream(ctx, jetstream.StreamConfig{Name: "gpress", Subjects: []string{"gpress.*"}})
	}
	if err != nil {
		return fmt.Errorf("js.CreateStream error: %w", err)
	}
	// retrieve consumer handle from a stream
	cons, _ := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "CONS",
		AckPolicy: jetstream.AckExplicitPolicy,
	})

	// consume messages from the consumer in callback
	cons.Consume(func(msg jetstream.Msg) {
		fmt.Println("Received jetstream message: ", string(msg.Data()))
		msg.Ack()
	})

	// Simple Publisher
	nc.Publish("gpress.hello", []byte("hello gpress"))

	//}()
	// 等待一段时间
	//time.Sleep(time.Second * 3)
	return err
}

func closeNatsServer() {
	if ns != nil {
		ns.Shutdown()
	}
	if nc != nil {
		nc.Close()
	}
}
