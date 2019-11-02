package heartbeat

import (
	"OSS/app/dataServer/config"
	"OSS/comm/rabbitmq"
	"time"
)

func Heartbeat() {
	mq := rabbitmq.New(config.ServerCfg.Server.RabbitMq)
	defer mq.Close()
	for {
		mq.Publish("apiServer", config.ServerCfg.Server.Address)
		time.Sleep(time.Second * 5)
	}
}
