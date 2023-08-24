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
		ServerName: appName, //MQTT协议必须指定
		//Host:               "127.0.0.1",
		Port:               4222,
		JetStream:          true,
		JetStreamMaxMemory: 1 * 1024 * 1024 * 1024,  // 1G
		JetStreamMaxStore:  10 * 1024 * 1024 * 1024, // 10G
		StoreDir:           datadir + "natsdata",    // 持久化数据目录
		NoLog:              true,
		//Debug:              false,
		NoSigs:         true,
		MaxControlLine: 1024,
		//mqtt.js只支持webscoket协议,默认的path是 /mqtt
		Websocket: server.WebsocketOpts{
			NoTLS: true, //默认不使用TLS
			Port:  8083,
		},
		MQTT: server.MQTTOpts{ //启用MQTT协议,用于支持IOT/聊天等场景
			//Host: "127.0.0.1",
			Port: 1883,
		},
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
	//nc, err = nats.Connect(nsUrl, nats.TokenHandler(func() string { return "" }))
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
	stream, _ := js.Stream(ctx, appName)
	if stream == nil {
		stream, err = js.CreateStream(ctx, jetstream.StreamConfig{Name: appName, Subjects: []string{appName + ".*"}})
	}
	if err != nil {
		return fmt.Errorf("js.CreateStream error: %w", err)
	}

	//emit的事件订阅使用Subscribe模式,不进行持久化.可以和订阅发布使用同一个Subject
	//Queue Groups
	nc.QueueSubscribe("gpress.hello", "queueGroups", func(m *nats.Msg) {
		fmt.Printf("Received queueSubscribe message: %s\n", string(m.Data))
	})
	// Simple Async Subscriber
	nc.Subscribe("gpress.hello", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})

	// retrieve consumer handle from a stream
	cons, _ := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   appName + "-cons", //类似redis stream 组名称,同一个组内的消费者进行负载均衡
		AckPolicy: jetstream.AckExplicitPolicy,
	})

	// consume messages from the consumer in callback
	cons.Consume(func(msg jetstream.Msg) {
		fmt.Println("Received jetstream message: ", string(msg.Data()))
		msg.Ack()
	})

	// Simple Publisher
	nc.Publish("gpress.hello", []byte("hello gpress-->"+time.Now().Format("2006-01-02 15:04:05")))
	// TODO 如何后台发送一个mqtt消息
	//nc.Publish("/mqtt/hello", []byte("hello mqtt-->"+time.Now().Format("2006-01-02 15:04:05")))

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
