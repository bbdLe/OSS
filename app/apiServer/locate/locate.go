package locate

import (
	"OSS/app/apiServer/config"
	"OSS/comm/rabbitmq"
	"strconv"
	"time"
)

func Locate(file string) string {
	mq := rabbitmq.New(config.ServerCfg.Server.RabbitMq)
	mq.Publish("dataServer", file)
	c := mq.Cosume()
	go func() {
		time.Sleep(time.Second)
		mq.Close()
	}()
	msg := <- c
	s, _ := strconv.Unquote(string(msg.Body))
	return s
}
