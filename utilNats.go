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
	natsOptions := &server.Options{
		ServerName: appName, //MQTT协议必须指定
		//Host:               "127.0.0.1",
		Port: 4222,
		//HTTPPort:           8222, 默认不启用http监控页面
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

		//自定义验证客户端连接权限
		CustomClientAuthentication: &NatsClientAuthentication{},
	}

	if config.ExternalNats { //如果是使用外部独立的Nats服务
		natsOptions, err = server.ProcessConfigFile(datadir + "nats.conf")
		if err != nil {
			return fmt.Errorf("server.ProcessConfigFile error: %w", err)
		}
	}

	// 创建 NATS Server 实例
	ns, err = server.NewServer(natsOptions)
	if err != nil {
		return fmt.Errorf("server.NewServer error: %w", err)
	}

	// 启动 NATS Server 实例
	//go func() {
	ns.Start()

	//nc, err = ns.InProcessConn(ns)
	// Connect to a server
	nsUrl := fmt.Sprintf("nats://%s:%d", natsOptions.Host, natsOptions.Port)
	//nc, err = nats.Connect(nsUrl, nats.TokenHandler(func() string { return "" }))
	options := make([]nats.Option, 0)
	//用户账号密码可以启动时随机产生,用于特殊权限判断
	options = append(options, nats.UserInfo("user", "password"))
	if natsOptions.Host == "" || natsOptions.Host == "0.0.0.0" || natsOptions.Host == "127.0.0.1" || natsOptions.Host == "localhost" { //避免建立TCP连接
		options = append(options, nats.InProcessServer(ns))
	}
	nc, err = nats.Connect(nsUrl, options...)

	if err != nil {
		return fmt.Errorf("nats.Connect(nsUrl) error: %w", err)
	}

	// create jetstream context from nats connection
	js, err = jetstream.New(nc)
	if err != nil {
		return fmt.Errorf("jetstream.New(nc) error: %w", err)
	}
	//ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	ctx := context.Background()

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

	// 订阅mqtt消息,用于历史记录
	nc.Subscribe("/.mqtt.hello", func(m *nats.Msg) {
		fmt.Printf("/.mqtt.hello==>: %s\n", string(m.Data))
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

	// 发布需要在订阅之后,订阅是不持久化的,没有订阅者,这条消息就删除了.

	// Simple Publisher
	nc.Publish("gpress.hello", []byte("hello gpress-->"+time.Now().Format("2006-01-02 15:04:05")))

	//https://docs.nats.io/running-a-nats-service/configuration/mqtt
	// http://127.0.0.1:8222/subsz?subs=1 查看对照 /mqtt/hello ==> /.mqtt.hello
	nc.Publish("/.mqtt.hello", []byte("hello mqtt-->"+time.Now().Format("2006-01-02 15:04:05")))

	// 等待一段时间
	//time.Sleep(time.Second * 3)
	return err
}

func closeNatsServer() {
	if nc != nil {
		nc.Close()
	}

	if ns != nil {
		ns.Shutdown()
	}

}

// 自定义连接权限验证接口,可以用于验证私钥签名
type NatsClientAuthentication struct{}

func (client *NatsClientAuthentication) Check(c server.ClientAuthentication) bool {

	//fmt.Printf("server.ClientAuthentication:%v", c)
	fmt.Printf("userName:%s,password:%s", c.GetOpts().Username, c.GetOpts().Password)

	//可以把password作为加密签名,登录成功之后,也可以作为token

	//Subscribe 和 Publish 需要自定义函数,验证用户权限

	//用户权限控制,需要限制某个用户能够订阅或者发布某个主题
	// https://docs.nats.io/running-a-nats-service/configuration/securing_nats/authorization
	c.RegisterUser(&server.User{
		Username: c.GetOpts().Username,
		Permissions: &server.Permissions{
			//订阅权限
			Subscribe: &server.SubjectPermission{
				//Allow: []string{},           //允许
				//Deny: []string{"/.mqtt.>"}, //拒绝
			},
			//发布权限
			Publish: &server.SubjectPermission{
				//Allow: []string{},           //允许
				//Deny: []string{"/.mqtt.>"}, //拒绝
			},
		},
	})

	return true
}
