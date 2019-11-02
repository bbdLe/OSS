package locate

import (
	"OSS/app/dataServer/config"
	"OSS/comm/rabbitmq"
	"os"
	"path"
	"strconv"
)

func IsExist(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func Locate() {
	mq := rabbitmq.New(config.ServerCfg.Server.RabbitMq)
	defer mq.Close()
	mq.Bind("dataServer")
	c := mq.Cosume()
	for msg := range c {
		file, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			panic(err)
		}
		if IsExist(path.Join(config.ServerCfg.Server.StoragePath, "objects", file)) {
			mq.Send(msg.ReplyTo, config.ServerCfg.Server.Address)
		}
	}
}
