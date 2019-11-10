package locate

import (
	"OSS/app/apiServer/config"
	"OSS/comm/rabbitmq"
	"OSS/comm/rs"
	"OSS/comm/types"
	"encoding/json"
	"time"
)

func Locate(file string) (locateInfo map[int]string){
	mq := rabbitmq.New(config.ServerCfg.Server.RabbitMq)
	mq.Publish("dataServer", file)
	c := mq.Cosume()
	go func() {
		time.Sleep(time.Second)
		mq.Close()
	}()
	locateInfo = make(map[int]string)
	for i := 0; i < rs.AllShares; i++ {
		msg := <- c
		if len(msg.Body) == 0 {
			return
		}
		var info types.LocateMessage
		err := json.Unmarshal(msg.Body, &info)
		if err != nil {
			return
		}
		locateInfo[info.Id] = info.Addr
	}
	return
}

func Exist(name string) bool {
	return len(Locate(name)) >= rs.DataShards
}
